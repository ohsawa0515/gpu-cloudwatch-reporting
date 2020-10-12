package main

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/cloudwatch/cloudwatchiface"
	"github.com/kelseyhightower/envconfig"
)

const (
	// DimensionAsg indicates dimension CloudWatch as auto scaling group.
	DimensionAsg = "AutoScalingGroupName"
	// DimensionEC2 indicates dimension CloudWatch as EC2.
	DimensionEC2 = "InstanceId"
)

type metric struct {
	metricName string
	unit       string
}

// Client -
type Client struct {
	instanceID           string
	autoScalingGroupName string
	cloudwatchSvc        cloudwatchiface.CloudWatchAPI
	ctx                  *context.Context
}

// CloudWatchConfig is the defined config by envconfig.
type CloudWatchConfig struct {
	NameSpace string `default:"GPUMonitor"`
	Region    string `default:"us-east-1"`
}

var (
	metrics  map[string]metric
	cwConfig CloudWatchConfig
)

func init() {
	metrics = map[string]metric{}
	metrics["gpu_usage"] = metric{metricName: "GPUUsage", unit: "Percent"}
	metrics["memory_usage"] = metric{metricName: "MemoryUsage", unit: "Percent"}
	metrics["temperature"] = metric{metricName: "Temperature", unit: "None"}

	if err := envconfig.Process("", &cwConfig); err != nil {
		log.Fatal(err)
	}
}

// ReportGpuMetrics send gpu metrics to CloudWatch.
func (client *Client) ReportGpuMetrics(metric string, value float64) error {
	if _, err := client.cloudwatchSvc.PutMetricData(&cloudwatch.PutMetricDataInput{
		Namespace:  aws.String(cwConfig.NameSpace),
		MetricData: client.buildCloudWatchMetricDatum(metric, value),
	}); err != nil {
		return err
	}
	return nil
}

func (client *Client) buildCloudWatchMetricDatum(metric string, value float64) []*cloudwatch.MetricDatum {
	var metricDatum []*cloudwatch.MetricDatum
	timestamp := time.Now().UTC()
	metricDatum = append(metricDatum, &cloudwatch.MetricDatum{
		Dimensions: []*cloudwatch.Dimension{
			{
				Name:  aws.String(DimensionEC2),
				Value: aws.String(client.instanceID),
			},
		},
		MetricName: aws.String(metrics[metric].metricName),
		Unit:       aws.String(metrics[metric].unit),
		Timestamp:  aws.Time(timestamp),
		Value:      aws.Float64(value),
	})

	// if asg is found, add dimension of asg name.
	if len(client.autoScalingGroupName) > 0 {
		metricDatum[0].Dimensions = append(metricDatum[0].Dimensions, &cloudwatch.Dimension{
			Name:  aws.String(DimensionAsg),
			Value: aws.String(client.autoScalingGroupName),
		})
		metricDatum = append(metricDatum, &cloudwatch.MetricDatum{
			Dimensions: []*cloudwatch.Dimension{
				{
					Name:  aws.String(DimensionAsg),
					Value: aws.String(client.autoScalingGroupName),
				},
			},
			MetricName: aws.String(metrics[metric].metricName),
			Unit:       aws.String(metrics[metric].unit),
			Timestamp:  aws.Time(timestamp),
			Value:      aws.Float64(value),
		})
	}
	return metricDatum
}
