package main

import (
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/service/cloudwatch/cloudwatchiface"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

// NewClient creates a new client of CloudWatch.
func NewClient(ec2mdSvc *ec2metadata.EC2Metadata, ec2Svc ec2iface.EC2API, cwSvc cloudwatchiface.CloudWatchAPI) (*Client, error) {
	instanceID, err := getInstanceID(ec2mdSvc)
	if err != nil {
		return nil, err
	}
	autoScalingGroupName, err := getAutoScalingGroupName(ec2Svc, instanceID)
	if err != nil {
		return nil, err
	}
	return &Client{
		instanceID:           instanceID,
		autoScalingGroupName: autoScalingGroupName,
		cloudwatchSvc:        cwSvc,
	}, nil
}
