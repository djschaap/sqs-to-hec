package receivesqs

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

var svc *sqs.SQS

func OpenSvc() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc = sqs.New(sess)
	return
}

func DeleteMessages(queue_url string, sqs_result *sqs.ReceiveMessageOutput) {
	for i, m := range sqs_result.Messages {
		resultDelete, err := svc.DeleteMessage(&sqs.DeleteMessageInput{
			QueueUrl:      &queue_url,
			ReceiptHandle: m.ReceiptHandle,
		})
		if err != nil {
			fmt.Println("Delete Error [", i, "]", err)
			// ignore error and attempt to CONTINUE deletions
		} else if false {
			fmt.Println("Message [", i, "] deleted", resultDelete)
		}
	}
}

func ReceiveMessages(queue_url string, max_messages int64, max_wait int64) (*sqs.ReceiveMessageOutput, error) {
	result, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
		AttributeNames: []*string{
			aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
		},
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl:            &queue_url,
		MaxNumberOfMessages: aws.Int64(max_messages),
		VisibilityTimeout:   aws.Int64(10), // seconds
		WaitTimeSeconds:     aws.Int64(max_wait),
	})
	if err != nil {
		//fmt.Println("Error", err)
		return nil, err
	}
	//message_count := len(result.Messages)

	return result, nil
}
