package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func main() {
	sess := session.Must(session.NewSession(
		&aws.Config{Region: aws.String(cwConfig.Region)},
	))
	client, err := NewClient(ec2metadata.New(sess), ec2.New(sess), cloudwatch.New(sess))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("REGION: %s", cwConfig.Region)
	log.Printf("NAMESPACE: %s", cwConfig.NameSpace)
	log.Printf("SEND_INTERVAL_SECOND: %s", nvidiaConfig.SendIntervalSecond)
	log.Printf("COLLECT_INTERVAL_SECOND: %s", nvidiaConfig.CollectIntervalSecond)

	ctx := context.Background()
	if err := client.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
