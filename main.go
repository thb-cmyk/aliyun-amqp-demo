package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/thb-cmyk/aliyun-amqp-demo/amqpbasic"

	"github.com/thb-cmyk/aliyun-amqp-demo/touser"
)

func main() {

}

func StartServer(configpath string) {
	//	setup a socket and listen the port
	configmap := touser.GetYamlConfig(configpath)
	host := touser.GetElement("host", configmap)
	timeinterval, err := strconv.Atoi(touser.GetElement("beatinginterval", configmap))
	touser.CheckError(err)
	netListen, err := net.Listen("tcp", host)
	touser.CheckError(err)
	defer netListen.Close()
	touser.Log("Waiting for clients")

	for {
		conn, err := netListen.Accept()
		if err != nil {
			continue
		}

		touser.Log(conn.RemoteAddr().String(), " tcp connect success")
		go handleConnection(conn, timeinterval)
	}

}

func handleConnection(conn net.Conn, timeinterval int) {
	defer conn.Close()
	timeout := make(chan byte)
	buffer := make([]byte, 512)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Printf("error information: %s\n\r", err)
		}
		touser.HeartBeating(conn, timeout, timeinterval)
		touser.GravelChannel(buffer, timeout)

		/* handle the received data */
		fmt.Printf("the data length is %d\n\r", n)
		fmt.Printf("the received data: %s\n\r", string(buffer))
	}
}

func Aliyun_Connect() {

	const accessKey = ""
	const accessSecret = ""
	const consumerGroupId = ""
	const clientId = ""
	const iotInstanceId = ""
	const host = ""

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

	ok := session_test.SessionInit(session_id, 1, ctx_top)
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

	go amqpbasic.ReceiveThread(ctx_top)

}
