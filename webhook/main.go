package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/yevhenshymotiuk/ap-curriculum-bot/curriculum"
	"github.com/yevhenshymotiuk/telegram-lambda-helpers/apigateway"
)

func getObjectFromS3Bucket(
	bucketName string,
	objectName string,
) *s3.GetObjectOutput {
	sess, _ := session.NewSession(&aws.Config{Region: aws.String("eu-north-1")})

	client := s3.New(sess)

	resp, err := client.GetObject(
		&s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(objectName),
		},
	)

	if err != nil {
		log.Fatalf("Unable to get file %q, %v", objectName, err)
	}

	return resp
}

func handler(
	request events.APIGatewayProxyRequest,
) (apigateway.Response, error) {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		return apigateway.Response404, err
	}

	update := tgbotapi.Update{}

	err = json.Unmarshal([]byte(request.Body), &update)
	if err != nil {
		return apigateway.Response404, err
	}

	assetsBucket := os.Getenv("ASSETS_BUCKET")
	curriculumFile1 := os.Getenv("CURRICULUM_FILE1")
	curriculumFile2 := os.Getenv("CURRICULUM_FILE2")

	message := update.Message
	var responseMessageText string

	switch message.Command() {
	case "today":
		resp1 := getObjectFromS3Bucket(assetsBucket, curriculumFile1)
		resp2 := getObjectFromS3Bucket(assetsBucket, curriculumFile2)

		w1, err := curriculum.NewWeek(io.Reader(resp1.Body))
		if err != nil {
			return apigateway.Response404, err
		}
		w2, err := curriculum.NewWeek(io.Reader(resp2.Body))
		if err != nil {
			return apigateway.Response404, err
		}

		t1, t2 := curriculum.Today(*w1), curriculum.Today(*w2)

		// Formatted today's curriculums for 2 subgroups
		var ftc1, ftc2 string
		ftc1 = t1.Format()

		if reflect.DeepEqual(t1, t2) {
			responseMessageText = fmt.Sprintf("Розклад однаковий для обох підгруп:\n%s", ftc1)
		} else {
			ftc2 = t2.Format()
			responseMessageText = fmt.Sprintf("Підгрупа 1:\n%s\n\nПідгрупа 2:\n%s", ftc1, ftc2)
		}
	default:
		responseMessageText = `¯\_(ツ)_/¯`
	}

	responseMessage := tgbotapi.NewMessage(message.Chat.ID, responseMessageText)
	bot.Send(responseMessage)

	return apigateway.Response200, nil
}

func main() {
	lambda.Start(handler)
}
