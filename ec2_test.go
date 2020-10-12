package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/service/ec2"
)

const (
	expectedTagValue   = "bar-asg"
	expectedInstanceID = "i-1234567890ab"
)

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
		t.Errorf("expected value of tag: %s, but got tag value: %v.", expectedTagValue, tag)
	}
}

func TestGetInstanceID(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/latest/meta-data/instance-id", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, expectedInstanceID)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	mockEC2MetaDataSvc := ec2metadata.New(mockSession, &aws.Config{
		Endpoint: aws.String(server.URL),
	})

	instanceID, err := getInstanceID(mockEC2MetaDataSvc)
	if err != nil {
		t.Errorf("expected no error, but got %v.", err)
	}
	if instanceID != expectedInstanceID {
		t.Errorf("expected value of tag: %s, but got tag value: %v.", expectedInstanceID, instanceID)
	}
}
