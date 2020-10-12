package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch/cloudwatchiface"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

type mockEC2Client struct {
	ec2iface.EC2API
}

type mockCloudWatchClient struct {
	cloudwatchiface.CloudWatchAPI
}

var mockSession = session.Must(session.NewSession(&aws.Config{
	Region: aws.String("mock-region"),
}))

func createMockClient() (*Client, error) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("/latest/meta-data/instance-id", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "i-1234567890ab")
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	mockEC2MetaDataSvc := ec2metadata.New(mockSession, &aws.Config{
		Endpoint: aws.String(server.URL),
	})
	mockEC2Svc := &mockEC2Client{}
	mockCloudWatchSvc := &mockCloudWatchClient{}
	return NewClient(mockEC2MetaDataSvc, mockEC2Svc, mockCloudWatchSvc)
}
