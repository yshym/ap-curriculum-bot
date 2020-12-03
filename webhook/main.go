package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/yevhenshymotiuk/ap-curriculum-bot/curriculum"
	"github.com/yevhenshymotiuk/ap-curriculum-bot/helpers"
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

func specificDayText(ab, cf string, t time.Time) (string, error) {
	resp := getObjectFromS3Bucket(ab, cf)

	w, err := curriculum.NewWeek(io.Reader(resp.Body))
	if err != nil {
		return "", err
	}

	sd := curriculum.NewSpecificDay(
		*w,
		t,
	)

	return sd.Format(), nil
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
	curriculumFile := os.Getenv("CURRICULUM_FILE")

	message := update.Message
	log.Printf("Object: %+v\nText: %s", message, message.Text)
	var responseMessageText string

	switch message.Command() {
	case "today":
		responseMessageText, err = specificDayText(
			assetsBucket,
			curriculumFile,
			helpers.Now(),
		)
		if err != nil {
			return apigateway.Response404, err
		}
	case "tomorrow":
		responseMessageText, err = specificDayText(
			assetsBucket,
			curriculumFile,
			helpers.Now().AddDate(0, 0, 1),
		)
		if err != nil {
			return apigateway.Response404, err
		}
	case "day":
		if !strings.Contains(message.Text, " ") {
			responseMessageText = "You should provide an argument for '/day' command"
			break
		}

		splitted_text := strings.Split(message.Text, " ")

		t, err := helpers.FromFormatted(splitted_text[1])
		if err != nil {
			return apigateway.Response404, err
		}

		responseMessageText, err = specificDayText(
			assetsBucket,
			curriculumFile,
			*t,
		)
		if err != nil {
			return apigateway.Response404, err
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
