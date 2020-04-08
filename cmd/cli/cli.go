package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/djschaap/sqs-to-hec/pkg/receivesqs"
	"github.com/djschaap/sqs-to-hec/pkg/sendhec"
	hec "github.com/fuyufjh/splunk-hec-go"
	"net/http"
	"os"
	"regexp"
)

var (
	build_dt string
	commit   string
	version  string
)

func main() {
	fmt.Println("sqs-to-hec  Version:", version, " Commit:", commit,
		" Built at:", build_dt)
	var max_receive_messages int64 = 2 // TODO create MAX_RECEIVE_MESSAGES env var
	var max_receive_wait int64 = 2     // in seconds
	//max_send_events := 2 // TODO create MAX_SEND_EVENTS env var

	src_queue_url := os.Getenv("SRC_QUEUE")
	has_src_queue, _ := regexp.MatchString(`^https`, src_queue_url)
	if !has_src_queue {
		fmt.Println("ERROR: SRC_QUEUE must be set")
		os.Exit(1)
	}

	hec_url := os.Getenv("HEC_URL")
	hec_token := os.Getenv("HEC_TOKEN")
	has_hec_url, _ := regexp.MatchString(`\S`, hec_url)
	has_hec_token, _ := regexp.MatchString(`\S`, hec_token)

	receivesqs.OpenSvc()

	if has_hec_url && has_hec_token {
		hec_client := hec.NewCluster(
			[]string{hec_url},
			hec_token,
		)
		if true { // TODO set InsecureSkipVerify from env var
			hec_client.SetHTTPClient(&http.Client{Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}})
		}
	}

	for {
		sqs_result, err := receivesqs.ReceiveMessages(src_queue_url, max_receive_messages, max_receive_wait)

		if err != nil {
			fmt.Println("ERROR from ReceiveMessages:", err)
			os.Exit(2)
			// TODO retry/attempt recovery?
		}

		var hec_messages []hec.Event
		if len(sqs_result.Messages) > 0 {
			hec_messages = sendhec.FormatForHEC(sqs_result.Messages)
			//debug_bytes, _ := json.Marshal(sqs_result.Messages)
			//fmt.Println("  DEBUG:", string(debug_bytes))
		}

		if len(hec_messages) > 0 {
			var send_err error
			if has_hec_url && has_hec_token {
				send_err = sendhec.SendToHEC(hec_url, hec_token, hec_messages)
				if send_err != nil {
					fmt.Println("ERROR from SendToHEC:", send_err)
					os.Exit(3)
					// TODO retry/attempt recovery?
				}
			} else {
				fmt.Println("WARNING: HEC_URL and/or HEC_TOKEN are missing")
				if hec_messages != nil {
					for _, m := range hec_messages {
						debug_bytes, _ := json.Marshal(m)
						fmt.Println("NOT-SENT-TO-HEC:", string(debug_bytes))
					}
				}
			}
			if send_err == nil {
				receivesqs.DeleteMessages(src_queue_url, sqs_result)
			}
		}
	}
}
