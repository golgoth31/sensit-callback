// Copyright © 2018 David Sabatie <david.sabatie@notrenet.com>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.
package sensitsqs

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/golgoth31/sensit-callback/config"
	"github.com/golgoth31/sensit-callback/sensitTypes"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// var readqURL = "https://sqs.eu-west-1.amazonaws.com/289326297244/sensit.fifo"
// var writeqURL = "https://sqs.eu-west-1.amazonaws.com/289326297244/sensitReturn.fifo"

var getMess = true
var Data sensittypes.CallbackData
var cfg = config.Config

var readqURL string

// Define AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY and AWS_REGION as env vars
func init() {
	readqURL = cfg.GetString("input.mode.readqURL")
}

// GetMessage extract messages from SQS
func GetMessage(payChan chan sensittypes.CallbackData) {
	log.Print("[DEBUG] Starting SQS module")
	// sess := session.Must(session.NewSession())
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := sqs.New(sess)

	for getMess {
		result, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
			AttributeNames: []*string{
				aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
			},
			MessageAttributeNames: []*string{
				aws.String(sqs.QueueAttributeNameAll),
			},
			QueueUrl:            &readqURL,
			MaxNumberOfMessages: aws.Int64(1),
			VisibilityTimeout:   aws.Int64(20), // 20 seconds
			WaitTimeSeconds:     aws.Int64(10),
		})

		if err != nil {
			fmt.Printf("[ERROR] SQS Error: %v", err)
			return
		}

		if len(result.Messages) == 0 {
			log.Print("[DEBUG] Waiting for message")
		} else {
			for _, val := range result.Messages {
				err := json.Unmarshal([]byte(*val.Body), &Data)
				if err != nil {
					fmt.Println(err)
				}
				datatime, _ := strconv.ParseInt(Data.Time, 10, 64)
				log.Printf("[DEBUG] Message time: %v", time.Unix(datatime, 0))
				log.Printf("[DEBUG] Message data: %v", val.Body)
				Data.Timestamp = time.Unix(datatime, 0)
				payChan <- Data
				_, err = svc.DeleteMessage(&sqs.DeleteMessageInput{
					QueueUrl:      &readqURL,
					ReceiptHandle: val.ReceiptHandle,
				})

				if err != nil {
					log.Printf("[DEBUG] Delete Error: %v", err)
					return
				}
			}
		}
	}
}
