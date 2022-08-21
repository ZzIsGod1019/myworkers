package sqs

import (
	"errors"
	"myworkers/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func GetQueueURL(sess *session.Session, queue *string) (*sqs.GetQueueUrlOutput, error) {
	svc := sqs.New(sess)
	urlResult, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: queue,
	})
	if err != nil {
		return nil, err
	}

	return urlResult, nil
}

func getMessages(sess *session.Session, queueURL *string, timeout *int64) (*sqs.ReceiveMessageOutput, error) {
	svc := sqs.New(sess)

	msgResult, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
		AttributeNames: []*string{
			aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
		},
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl:            queueURL,
		MaxNumberOfMessages: aws.Int64(1),
		VisibilityTimeout:   timeout,
	})
	if err != nil {
		return nil, err
	}

	return msgResult, nil
}

func ReceiveMessage() (*sqs.Message, error) {
	queue := config.AwsSqsConf["queue"].(string)
	region := config.AwsSqsConf["region"].(string)
	timeout := int64(config.AwsSqsConf["timeout"].(int))
	accessKeyId := config.AwsSqsConf["access_key_id"].(string)
	secretAccessKey := config.AwsSqsConf["secret_access_key"].(string)

	if queue == "" {
		return nil, errors.New("队列名不能为空")
	}

	if timeout < 0 {
		timeout = 0
	}

	if timeout > 12*60*60 {
		timeout = 12 * 60 * 60
	}

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config: aws.Config{
			Region:      aws.String(region),
			Credentials: credentials.NewStaticCredentials(accessKeyId, secretAccessKey, ""),
		},
	}))

	urlResult, err := GetQueueURL(sess, &queue)
	if err != nil {
		return nil, err
	}

	queueURL := urlResult.QueueUrl

	msgResult, err := getMessages(sess, queueURL, &timeout)
	if err != nil {
		return nil, err
	}
	return msgResult.Messages[0], nil
}
