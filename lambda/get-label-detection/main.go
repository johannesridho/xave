package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rekognition"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type Configuration struct {
	Region                 string
	FbMessengerAccessToken string
	PsId                   string
}

type SnsMessage struct {
	JobId string `json:"JobId"`
}

type FbSendMessageReq struct {
	MessagingType string    `json:"messaging_type"`
	Recipient     Recipient `json:"recipient"`
	Message       Message   `json:"message"`
}

type Recipient struct {
	Id string `json:"id"`
}

type Message struct {
	Text string `json:"text"`
}

func main() {
	lambda.Start(GetLabelDetection)
}

func GetLabelDetection(ctx context.Context, event events.SNSEvent) (string, error) {
	log.Printf("start GetLabelDetection, event: %v", event)

	jsonMessage := event.Records[0].SNS.Message
	log.Printf("SNS message: %v", jsonMessage)

	snsMessage := SnsMessage{}
	json.Unmarshal([]byte(jsonMessage), &snsMessage)

	jobId := snsMessage.JobId
	log.Printf("Rekognition jobId: %s", jobId)

	config := Configuration{
		Region:                 os.Getenv("REGION"),
		FbMessengerAccessToken: os.Getenv("FB_MESSENGER_ACCESS_TOKEN"),
		PsId:                   os.Getenv("FB_MESSENGER_PSID"),
	}

	log.Printf("start GetLabelDetection")
	result, err := getLabelDetectionResult(config, jobId)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("GetLabelDetection result acquired, status: %s", *result.JobStatus)

	suspiciousLabels := [...]string{"flare", "fire", "flame", "quake", "duel", "kicking", "punching", "fighting", "martial art", "wrestling", "boxing"}
	suspiciousLabelDetected := make(map[string]bool)

	for _, label := range result.Labels {
		for _, suspiciousLabel := range suspiciousLabels {
			if strings.ToLower(*label.Label.Name) == suspiciousLabel {
				suspiciousLabelDetected[suspiciousLabel] = true
			} else if strings.ToLower(*label.Label.Name) == "grand theft auto" {
				suspiciousLabelDetected["road conflict/fighting"] = true
			}
		}
	}

	var strBuilder strings.Builder
	strBuilder.WriteString(fmt.Sprintf("Analysis result for video with job id: %s\n\n", jobId))

	if len(suspiciousLabelDetected) == 0 {
		strBuilder.WriteString("There is no suspicious activity detected in this video")
	} else {
		strBuilder.WriteString("Detected suspicious activities :\n")
		for key := range suspiciousLabelDetected {
			strBuilder.WriteString(fmt.Sprintf("%s\n", key))
		}
	}

	message := strBuilder.String()
	log.Println(message)
	sendToFb(config, message)

	return "success", nil
}

func getLabelDetectionResult(config Configuration, jobId string) (*rekognition.GetLabelDetectionOutput, error) {
	session, err := session.NewSession(&aws.Config{Region: aws.String(config.Region)})
	if err != nil {
		return nil, err
	}

	rek := rekognition.New(session)

	input := rekognition.GetLabelDetectionInput{JobId: &jobId}

	result, err := rek.GetLabelDetection(&input)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func sendToFb(config Configuration, message string) {
	log.Println(message)

	req := FbSendMessageReq{
		Message:       Message{Text: message},
		MessagingType: "RESPONSE",
		Recipient:     Recipient{Id: config.PsId},
	}

	url := fmt.Sprintf("https://graph.facebook.com/v3.2/me/messages?access_token=%s", config.FbMessengerAccessToken)

	payload, err := json.Marshal(req)
	res, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	responseBytes, err := ioutil.ReadAll(res.Body)
	log.Printf("received response: %s, status code: %d", string(responseBytes), res.StatusCode)
}
