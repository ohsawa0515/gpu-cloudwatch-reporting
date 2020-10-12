package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

// AutoScalingTagName indicates tag name created when it runs under auto scaling group.
const AutoScalingTagName = "aws:autoscaling:groupName"

func getInstanceID(svc *ec2metadata.EC2Metadata) (string, error) {
	return svc.GetMetadata("instance-id")
}

func getAutoScalingGroupName(svc ec2iface.EC2API, instanceID string) (string, error) {
	result, err := svc.DescribeTags(&ec2.DescribeTagsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("resource-id"),
				Values: []*string{aws.String(instanceID)},
			},
		},
	})
	if err != nil {
		return "", err
	}
	for _, t := range result.Tags {
		if *t.Key == AutoScalingTagName {
			return *t.Value, nil
		}
	}
	return "", nil
}
