package main

import (
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
)

func getInstanceId(sess *session.Session) (string, error) {
	svc := ec2metadata.New(sess)
	return svc.GetMetadata("instance-id")
}