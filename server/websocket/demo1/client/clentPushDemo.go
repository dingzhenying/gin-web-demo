package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	wserver "gin-web-demo/server/websocket/demo1/server"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func main() {
	pushURL := "http://127.0.0.1:12345/push"
	contentType := "application/json"
	// 保存token +userId

	for {
		timestamp := strconv.FormatInt(time.Now().UTC().Unix()*1000, 10)
		rand.Seed(int64(time.Now().Nanosecond()))
		values := 100 * rand.Float32()

		message := fmt.Sprintf("["+
			"{\"namespace\":\"000000\","+
			"\"internalSeriesId\":\"hiacloud0003000098L[]\","+
			"\"regions\":10,"+
			"\"t\": %s,"+
			"\"s\":0,"+
			"\"v\":\"L#%2.2f \","+
			"\"gatewayId\":\"hiacloud\","+
			"\"pointId\":\"0003000098L\""+
			"}"+
			"]", timestamp, values)

		//message:= fmt.Sprintf("Hello in %s", time.Now().Format("2006-01-02 15:04:05.000"))

		pm := wserver.PushMessage{
			UserID:  "jack",
			Event:   "topic1",
			Message: message,
		}
		fmt.Println(message)
		b, _ := json.Marshal(pm)

		http.DefaultClient.Post(pushURL, contentType, bytes.NewReader(b))

		time.Sleep(time.Second)
	}
}
