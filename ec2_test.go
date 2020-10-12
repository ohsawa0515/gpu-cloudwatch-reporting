package main

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

const expectedTagValue = "bar-asg"

type mockEC2Client struct {
	ec2iface.EC2API
}

func (m *mockEC2Client) DescribeTags(*ec2.DescribeTagsInput) (*ec2.DescribeTagsOutput, error) {
	return &ec2.DescribeTagsOutput{
		Tags: []*ec2.TagDescription{
			{
				Key:   aws.String("Name"),
				Value: aws.String("foo-instance"),
			},
			{
				Key:   aws.String("aws:autoscaling:groupName"),
				Value: aws.String(expectedTagValue),
			},
			{
				Key:   aws.String("baz"),
				Value: aws.String("hoge"),
			},
		},
	}, nil
}

func TestGetAutoScalingGroupName(t *testing.T) {
	mockSvc := &mockEC2Client{}
	tag, err := getAutoScalingGroupName(mockSvc, "i-1234567890abcdefg")
	if err != nil {
		t.Errorf("expected no error, but got %v.", err)
	}
	if tag != expectedTagValue {
		t.Errorf("Expected value of tag: %s, but got tag value: %v.", expectedTagValue, tag)
	}
}
