package sendhec

import (
	"crypto/tls"
	//"encoding/json"
	//"fmt"
	"github.com/aws/aws-sdk-go/service/sqs"
	hec "github.com/fuyufjh/splunk-hec-go"
	"net/http"
	"strconv"
	"time"
)

func FormatForHEC(sqs_messages []*sqs.Message) []hec.Event {
	var hec_messages []hec.Event

	for _, m := range sqs_messages {
		//message_json_bytes, _ := json.Marshal(m)
		//fmt.Println("Message:\n", string(message_json_bytes))
		//fmt.Println("  Body:", string(*m.Body))

		hec_message := hec.NewEvent(string(*m.Body))

		//if m.MessageAttributes["customer_code"] != nil {
		//fmt.Println("  customer_code:", string(*m.MessageAttributes["customer_code"].StringValue))
		//}
		if m.MessageAttributes["host"] != nil {
			v := string(*m.MessageAttributes["host"].StringValue)
			//fmt.Println("  host:", v)
			hec_message.SetHost(v)
		}
		if m.MessageAttributes["index"] != nil {
			v := string(*m.MessageAttributes["index"].StringValue)
			//fmt.Println("  index:", v)
			hec_message.SetIndex(v)
		}
		if m.MessageAttributes["source"] != nil {
			v := string(*m.MessageAttributes["source"].StringValue)
			//fmt.Println("  source:", v)
			hec_message.SetSource(v)
		}
		if m.MessageAttributes["sourcetype"] != nil {
			v := string(*m.MessageAttributes["sourcetype"].StringValue)
			//fmt.Println("  sourcetype:", v)
			hec_message.SetSourceType(v)
		}
		if m.Attributes["SentTimestamp"] != nil {
			i, err := strconv.Atoi(*m.Attributes["SentTimestamp"])
			if err == nil {
				time_s := int64(i / 1000)
				time_us := ( int64(i) - time_s * 1000 ) * 1000000
				v := time.Unix(time_s, time_us)
				//fmt.Println("  SentTimeStamp:", v)
				hec_message.SetTime(v)
			}
		}

		//if len(body["time"]) > 0 {
		//	hec_message.SetTime(body["time"])
		// cannot use body["time"] (type string) as type time.Time in argument to hec_message.SetTime
		//}

		// SetFields ?
		// https://docs.splunk.com/Documentation/Splunk/8.0.1/Data/HECExamples

		hec_messages = append(hec_messages, *hec_message)

		//for k, m := range hec_messages {
		//fmt.Println("  host[", k, "]: ", m.Host)
		//debug_bytes, _ := json.Marshal(m)
		//fmt.Println("  Message ", i, ":\n", string(debug_bytes))
		//}
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
