package sqstohec

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/djschaap/sqs-to-hec/receivesqs"
	"github.com/djschaap/sqs-to-hec/sendhec"
	hec "github.com/fuyufjh/splunk-hec-go"
	"net/http"
	"os"
	"regexp"
)

type HecConfig struct {
	Token string
	Url   string
}

type SqsConfig struct {
	Url string
}

type sess struct {
	HecConfig HecConfig
	SqsConfig SqsConfig
}

func (self *sess) RunForever() {
	var maxReceiveMessages int64 = 2 // TODO create MAX_RECEIVE_MESSAGES env var
	var maxReceiveWait int64 = 2     // in seconds
	//maxSendEvents := 2 // TODO create MAX_SEND_EVENTS env var

	receivesqs.OpenSvc()

	hasHecUrl, _ := regexp.MatchString(`\S`, self.HecConfig.Url)
	hasHecToken, _ := regexp.MatchString(`\S`, self.HecConfig.Token)
	if hasHecUrl && hasHecToken {
		hecClient := hec.NewCluster(
			[]string{self.HecConfig.Url},
			self.HecConfig.Token,
		)
		if true { // TODO set InsecureSkipVerify from env var
			hecClient.SetHTTPClient(&http.Client{Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}})
		}
	}

	for {
		mqResult, err := receivesqs.ReceiveMessages(self.SqsConfig.Url, maxReceiveMessages, maxReceiveWait)

		if err != nil {
			fmt.Println("ERROR from ReceiveMessages:", err)
			os.Exit(2)
			// TODO retry/attempt recovery?
		}

		var hecMmessages []hec.Event
		if len(mqResult.Messages) > 0 {
			hecMmessages = sendhec.FormatForHEC(mqResult.Messages)
			//debugBytes, _ := json.Marshal(mqResult.Messages)
			//fmt.Println("  DEBUG:", string(debugBytes))
		}

		if len(hecMmessages) > 0 {
			var sendErr error
			if hasHecUrl && hasHecToken {
				sendErr = sendhec.SendToHEC(self.HecConfig.Url, self.HecConfig.Token, hecMmessages)
				if sendErr != nil {
					fmt.Println("ERROR from SendToHEC:", sendErr)
					os.Exit(3)
					// TODO retry/attempt recovery?
				}
			} else {
				fmt.Println("WARNING: HEC_URL and/or HEC_TOKEN are missing")
				if hecMmessages != nil {
					for _, m := range hecMmessages {
						debugBytes, _ := json.Marshal(m)
						fmt.Println("NOT-SENT-TO-HEC:", string(debugBytes))
					}
				}
			}
			if sendErr == nil {
				receivesqs.DeleteMessages(self.SqsConfig.Url, mqResult)
			}
		}
	}
}

func New(
	hecConfig HecConfig,
	sqsConfig SqsConfig,
) sess {
	s := sess{
		HecConfig: hecConfig,
		SqsConfig: sqsConfig,
	}
	return s
}
