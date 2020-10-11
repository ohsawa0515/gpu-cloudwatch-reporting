package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

const AutoScalingTagName = "aws:autoscaling:groupName"

func getInstanceId(sess *session.Session) (string, error) {
	svc := ec2metadata.New(sess)
	return svc.GetMetadata("instance-id")
}

func getAutoScalingGroupName(sess *session.Session, instanceID string) (string, error) {
	svc := ec2.New(sess)
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
