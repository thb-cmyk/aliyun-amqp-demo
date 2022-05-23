package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/thb-cmyk/aliyun-amqp-demo/amqpbasic"
	"github.com/thb-cmyk/aliyun-amqp-demo/databasic"
	worker "github.com/thb-cmyk/aliyun-amqp-demo/dataworker"
)

//参数说明，请参见AMQP客户端接入说明文档。
const accessKey = "LTAI5tDTMKSxBGCgnbWyaPvt"
const accessSecret = "3Vmfy9VcAjIdX8V8HQMzDdnxucfrm3"
const consumerGroupId = "DEFAULT_GROUP"
const clientId = "00:0c:29:c4:01:22"

//iotInstanceId：实例ID。
const iotInstanceId = "iot-06z00bp0nwmb9tp"

//接入域名，请参见AMQP客户端接入说明文档。
const host = "iot-06z00bp0nwmb9tp.amqp.iothub.aliyuncs.com"

func main() {
	var ok int

	address := "amqps://" + host + ":5671"
	timestamp := time.Now().Nanosecond() / 1000000
	//userName组装方法，请参见AMQP客户端接入说明文档。
	username := fmt.Sprintf("%s|authMode=aksign,signMethod=Hmacsha1,consumerGroupId=%s,authId=%s,iotInstanceId=%s,timestamp=%d|",
		clientId, consumerGroupId, accessKey, iotInstanceId, timestamp)
	stringToSign := fmt.Sprintf("authId=%s&timestamp=%d", accessKey, timestamp)
	hmacKey := hmac.New(sha1.New, []byte(accessSecret))
	hmacKey.Write([]byte(stringToSign))
	//计算签名，password组装方法，请参见AMQP客户端接入说明文档。
	password := base64.StdEncoding.EncodeToString(hmacKey.Sum(nil))

	session_id := amqpbasic.SessionIdentifyInit(address, username, password, "session001")

	session_test := new(amqpbasic.AmqpSessionHandler)

	ctx_top := context.Background()

	ok = session_test.SessionInit(session_id, 1, ctx_top)
	if ok == -1 {
		fmt.Printf("The works of creating a new session is failed!\n\r")
		return
	}
	fmt.Printf("The works of creating a new session is successful!\n\r")

	ok = session_test.LinkCreate("recv001")
	if ok == -1 {
		fmt.Printf("The works of creating a new link is failed!\n\r")
		return
	}
	fmt.Printf("The works of creating a new link is successful!\n\r")

	databasic.All_Init()
	databasic.Broker()
	databasic.ProceNode_register(Aliyun_handler, "aliyun")

	go amqpbasic.ReceiveThread(ctx_top)

	go func() {
		for {
			//ctx_temp_recv, _ := context.WithTimeout(ctx_top, 100*time.Microsecond)
			buf, num := session_test.ReceiverData("recv001", 1)
			if num <= 0 {
				continue
			}
			rawnode := databasic.RawNode_create("aliyun", buf[0])
			databasic.Send_raw(rawnode)
		}
	}()

	for {

	}

}

func Aliyun_handler(tasknode *databasic.TaskNode, rawnode *databasic.RawNode) bool {

	dataentry, ok := worker.DataEntry_Register(rawnode.Raw.([]byte), "test001")
	if !ok {
		fmt.Printf("dataentry register return a error!\n\r")
	}

	ok = worker.DataInsert(dataentry)
	if !ok {
		fmt.Printf("data insert return a error!\n\r")
	}
	return true
}
