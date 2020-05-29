package sendhec

import (
	"crypto/tls"
	"encoding/json"
	//"fmt"
	"github.com/aws/aws-sdk-go/service/sqs"
	hec "github.com/fuyufjh/splunk-hec-go"
	//"github.com/kr/pretty"
	//"log"
	"net/http"
	"strconv"
	"time"
)

var trace_mq bool

func FormatForHEC(sqs_messages []*sqs.Message) []hec.Event {
	var hec_messages []hec.Event

	for _, m := range sqs_messages {
		var mq_message_body map[string]interface{}
		mq_message_body_bytes := []byte(*m.Body)
		json.Unmarshal(mq_message_body_bytes, &mq_message_body)
		//pretty.Log("TRACE mq_message_body = ", mq_message_body) // DEBUG

		var hec_message *hec.Event
		if _, ok := mq_message_body["event"]; ok {
			//log.Println("DEBUG 30: NEW style") // expect HEC headers within MQ body; allows metrics
			hec_message = hec.NewEvent(mq_message_body["event"])
			if mq_message_body["host"] != nil {
				v := mq_message_body["host"]
				hec_message.SetHost(v.(string))
			}
			if mq_message_body["index"] != nil {
				v := mq_message_body["index"]
				hec_message.SetIndex(v.(string))
			}
			if mq_message_body["source"] != nil {
				v := mq_message_body["source"]
				hec_message.SetSource(v.(string))
			}
			if mq_message_body["sourcetype"] != nil {
				v := mq_message_body["sourcetype"]
				hec_message.SetSourceType(v.(string))
			}
			if mq_message_body["fields"] != nil {
				v := mq_message_body["fields"]
				hec_message.SetFields(v.(map[string]interface{}))
			}
		} else {
			//log.Println("DEBUG 40: OLD style") // use MQ attributes for HEC headers
			hec_message = hec.NewEvent(string(*m.Body))
			if m.MessageAttributes["host"] != nil {
				v := string(*m.MessageAttributes["host"].StringValue)
				hec_message.SetHost(v)
			}
			if m.MessageAttributes["index"] != nil {
				v := string(*m.MessageAttributes["index"].StringValue)
				hec_message.SetIndex(v)
			}
			if m.MessageAttributes["source"] != nil {
				v := string(*m.MessageAttributes["source"].StringValue)
				hec_message.SetSource(v)
			}
			if m.MessageAttributes["sourcetype"] != nil {
				v := string(*m.MessageAttributes["sourcetype"].StringValue)
				hec_message.SetSourceType(v)
			}
		}
		//time.Sleep(3 * 100000) // DEBUG

		if m.Attributes["SentTimestamp"] != nil {
			i, err := strconv.Atoi(*m.Attributes["SentTimestamp"])
			if err == nil {
				time_s := int64(i / 1000)
				time_us := (int64(i) - time_s*1000) * 1000000
				v := time.Unix(time_s, time_us)
				//fmt.Println("  SentTimeStamp:", v)
				hec_message.SetTime(v)
			}
		}

		// SetFields ?
		// https://docs.splunk.com/Documentation/Splunk/8.0.1/Data/HECExamples

		hec_messages = append(hec_messages, *hec_message)
	}

	return hec_messages
}

func SendToHEC(hec_url string, hec_token string, hec_messages []hec.Event) error {
	hec_client := hec.NewCluster(
		[]string{hec_url},
		hec_token,
	)
	if true {
		hec_client.SetHTTPClient(&http.Client{Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}})
	}

	for _, e := range hec_messages {
		err := hec_client.WriteBatch([]*hec.Event{&e})
		if err != nil {
			//fmt.Println("HEC_ERROR:", err)
			return err
		}
	}

	return nil
}
