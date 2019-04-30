package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rekognition"
	"log"
	"os"
)

type Configuration struct {
	Region                 string
	BucketName             string
	LabelDetectionSnsTopic string
	RekRole                string
}

func main() {
	lambda.Start(StartLabelDetectionHandler)
}

func StartLabelDetectionHandler(ctx context.Context, event events.S3Event) (string, error) {
	log.Printf("start handling event: %v", event)

	videoFileName := event.Records[0].S3.Object.Key
	log.Printf("S3 video file name: %v", videoFileName)

	config := Configuration{
		Region:                 os.Getenv("REGION"),
		BucketName:             os.Getenv("S3_BUCKET_NAME"),
		LabelDetectionSnsTopic: os.Getenv("SNS_TOPIC_ARN"),
		RekRole:                os.Getenv("REKOGNITION_ROLE_ARN"),
	}

	video := rekognition.Video{
		S3Object: &rekognition.S3Object{
			Bucket: &config.BucketName,
			Name:   &videoFileName,
		},
	}

	labelDetectionJobId, err := startLabelDetection(config, video, "AidentStartLabelDetection", &config.LabelDetectionSnsTopic)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("label detection job id = %s", labelDetectionJobId)

	return "success", nil
}

func startLabelDetection(config Configuration, video rekognition.Video, jobTag string, labelDetectionSnsTopic *string) (string, error) {
	session, err := session.NewSession(&aws.Config{Region: aws.String(config.Region)})
	if err != nil {
		return "", err
	}

	rek := rekognition.New(session)

	minConfidence := float64(50)
	notificationChannel := rekognition.NotificationChannel{SNSTopicArn: labelDetectionSnsTopic, RoleArn: &config.RekRole}

	input := rekognition.StartLabelDetectionInput{
		Video:               &video,
		MinConfidence:       &minConfidence,
		JobTag:              &jobTag,
		NotificationChannel: &notificationChannel,
	}

	output, err := rek.StartLabelDetection(&input)
	if err != nil {
		return "", err
	}

	jobId := *output.JobId

	return jobId, nil
}
