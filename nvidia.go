package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NVIDIA/gpu-monitoring-tools/bindings/go/nvml"
	"github.com/kelseyhightower/envconfig"
	"golang.org/x/sync/errgroup"
)

// NvidiaConfig is the defined config by envconfig.
type NvidiaConfig struct {
	// SendIntervalSecond indicates the number of seconds between sending metrics to CloudWatch.
	SendIntervalSecond time.Duration `default:"60s"  split_words:"true"`
	// CollectIntervalSecond indicates the number of seconds between collecting GPU metrics.
	// Send the average value which is collected every SendIntervalSecond to CloudWatch.
	CollectIntervalSecond time.Duration `default:"5s"  split_words:"true"`
}

var nvidiaConfig NvidiaConfig

func init() {
	if err := envconfig.Process("", &nvidiaConfig); err != nil {
		log.Fatal(err)
	}
}

func (client *Client) gpuUtilizationTicker(ctx context.Context, devices []*nvml.Device) error {
	sendTicker := time.NewTicker(nvidiaConfig.SendIntervalSecond)
	collectTicker := time.NewTicker(nvidiaConfig.CollectIntervalSecond)
	defer sendTicker.Stop()
	defer collectTicker.Stop()

	metric := float64(0)
	count := 0
	for {
		select {
		case <-sendTicker.C:
			sentMetric := metric / float64(count)
			if err := client.ReportGpuMetrics("gpu_usage", float64(sentMetric)); err != nil {
				return err
			}
			metric = 0
			count = 0
		case <-collectTicker.C:
			for i, device := range devices {
				st, err := device.Status()
				if err != nil {
					return fmt.Errorf("error getting device %d status: %v", i, err)
				}
				metric += float64(*st.Utilization.GPU)
				count++
			}
		case <-ctx.Done():
			log.Println("Stop GPU utilization")
			return ctx.Err()
		}
	}
}

func (client *Client) gpuMemoryUtilizationTicker(ctx context.Context, devices []*nvml.Device) error {
	sendTicker := time.NewTicker(nvidiaConfig.SendIntervalSecond)
	collectTicker := time.NewTicker(nvidiaConfig.CollectIntervalSecond)
	defer sendTicker.Stop()
	defer collectTicker.Stop()

	metric := float64(0)
	count := 0
	for {
		select {
		case <-sendTicker.C:
			sentMetric := metric / float64(count)
			if err := client.ReportGpuMetrics("memory_usage", float64(sentMetric)); err != nil {
				return err
			}
			metric = 0
			count = 0
		case <-collectTicker.C:
			for i, device := range devices {
				st, err := device.Status()
				if err != nil {
					return fmt.Errorf("error getting device %d status: %v", i, err)
				}
				metric += float64(*st.Utilization.Memory)
				count++
			}
		case <-ctx.Done():
			log.Println("Stop GPU memory utilization ticker")
			return ctx.Err()
		}
	}
}

func (client *Client) gpuTemperatureTicker(ctx context.Context, devices []*nvml.Device) error {
	sendTicker := time.NewTicker(nvidiaConfig.SendIntervalSecond)
	collectTicker := time.NewTicker(nvidiaConfig.CollectIntervalSecond)
	defer sendTicker.Stop()
	defer collectTicker.Stop()

	metric := float64(0)
	count := 0
	for {
		select {
		case <-sendTicker.C:
			sentMetric := metric / float64(count)
			if err := client.ReportGpuMetrics("temperature", float64(sentMetric)); err != nil {
				return err
			}
		case <-collectTicker.C:
			for i, device := range devices {
				st, err := device.Status()
				if err != nil {
					return fmt.Errorf("error getting device %d status: %v", i, err)
				}
				metric += float64(*st.Temperature)
				count++
			}
		case <-ctx.Done():
			log.Println("Stop GPU temperature ticker")
			return ctx.Err()
		}
	}
}

func signalContext(ctx context.Context) context.Context {
	parent, cancelParent := context.WithCancel(ctx)

	go func() {
		defer cancelParent()

		sig := make(chan os.Signal, 1)
		signal.Notify(sig,
			syscall.SIGHUP,
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGQUIT)
		defer signal.Stop(sig)

		select {
		case <-parent.Done():
			log.Println("Cancel from parent")
			return
		case <-sig:
			log.Println("Cancel from signal")
			return
		}
	}()

	return parent
}

// Run runs collect and send gpu metrics.
func (client *Client) Run(ctx context.Context) error {
	if err := nvml.Init(); err != nil {
		return err
	}
	defer nvml.Shutdown()

	count, err := nvml.GetDeviceCount()
	if err != nil {
		return fmt.Errorf("error getting device count: %v", err)
	}

	var devices []*nvml.Device
	for i := uint(0); i < count; i++ {
		device, err := nvml.NewDevice(i)
		if err != nil {
			return fmt.Errorf("error getting device %d: %v", i, err)
		}
		devices = append(devices, device)
	}

	parent := signalContext(ctx)
	eg, child := errgroup.WithContext(parent)

	eg.Go(func() error {
		return client.gpuUtilizationTicker(child, devices)
	})

	eg.Go(func() error {
		return client.gpuMemoryUtilizationTicker(child, devices)
	})

	eg.Go(func() error {
		return client.gpuTemperatureTicker(child, devices)
	})

	return eg.Wait()
}
