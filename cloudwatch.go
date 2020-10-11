package main

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/cloudwatch/cloudwatchiface"
	"github.com/kelseyhightower/envconfig"
)

var (
	Metrics  map[string]Metric
	cwConfig CloudWatchConfig
)

const (
	DimensionAsg = "AutoScalingGroupName"
	DimensionEC2 = "InstanceID"
)

type Metric struct {
	MetricName string
	Unit       string
}

type Client struct {
	instanceID           string
	autoScalingGroupName string
	cloudwatchSvc        cloudwatchiface.CloudWatchAPI
	ctx                  *context.Context
}

type CloudWatchConfig struct {
	NameSpace string `default:"GPUMonitor"`
	Region    string `default:"us-east-1"`
}

func init() {
	Metrics = map[string]Metric{}
	Metrics["gpu_usage"] = Metric{MetricName: "GPUUsage", Unit: "Percent"}
	Metrics["memory_usage"] = Metric{MetricName: "MemoryUsage", Unit: "Percent"}
	Metrics["temperature"] = Metric{MetricName: "Temperature", Unit: "None"}

	if err := envconfig.Process("", &cwConfig); err != nil {
		log.Fatal(err)
	}
}

func NewClient(ctx context.Context) (*Client, error) {
	sess := session.Must(session.NewSession(
		&aws.Config{Region: aws.String(cwConfig.Region)},
	))
	instanceID, err := getInstanceId(sess)
	if err != nil {
		return nil, err
	}
	autoScalingGroupName, err := getAutoScalingGroupName(sess, instanceID)
	if err != nil {
		return nil, err
	}
	return &Client{
		instanceID:           instanceID,
		autoScalingGroupName: autoScalingGroupName,
		cloudwatchSvc:        cloudwatch.New(sess),
	}, nil
}

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
		MetricName: aws.String(Metrics[metric].MetricName),
		Unit:       aws.String(Metrics[metric].Unit),
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
			MetricName: aws.String(Metrics[metric].MetricName),
			Unit:       aws.String(Metrics[metric].Unit),
			Timestamp:  aws.Time(timestamp),
			Value:      aws.Float64(value),
		})
	}
	return metricDatum
}
