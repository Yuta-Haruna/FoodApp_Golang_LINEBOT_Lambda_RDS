package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/line/line-bot-sdk-go/linebot"
)

// API Gatewayから受け取ったevents.APIGatewayProxyRequestのBody（JSON）をパースする
// https://app.quicktype.io/　に、実際に受け取ってたJSONメッセージを張り付けて、コード自動生成。

// ▼▼▼ https://app.quicktype.io/で自動生成したコード：ここから ▼▼▼
func UnmarshalLineRequest(data []byte) (LineRequest, error) {
	var r LineRequest
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *LineRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type LineRequest struct {
	Events      []Event `json:"events"`
	Destination string  `json:"destination"`
}

type Event struct {
	Type       string  `json:"type"`
	ReplyToken string  `json:"replyToken"`
	Source     Source  `json:"source"`
	Timestamp  int64   `json:"timestamp"`
	Message    Message `json:"message"`
}

type Message struct {
	Type string `json:"type"`
	ID   string `json:"id"`
	Text string `json:"text"`
}

type Source struct {
	UserID string `json:"userId"`
	Type   string `json:"type"`
}

// ▲▲▲ https://app.quicktype.io/で自動生成したコード：ここまで ▲▲▲

// Handler
// fmt.Printlnやlog.Fatalは、CloudWatchのログで確認可能
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// 受け取ったJSONメッセージをログに書き込む（デバッグ用）
	fmt.Println("*** body")
	fmt.Println(request.Body)

	// JSONデコード
	fmt.Println("*** JSON decode")
	myLineRequest, err := UnmarshalLineRequest([]byte(request.Body))
	if err != nil {
		log.Fatal(err)
	}

	// ボットの定義
	fmt.Println("*** linebot new")
	bot, err := linebot.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_ACCESS_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// リプライ実施
	fmt.Println("*** reply")
	var tmpReplyMessage string
	tmpReplyMessage = "🦜オウムさん：" + myLineRequest.Events[0].Message.Text
	if _, err = bot.ReplyMessage(myLineRequest.Events[0].ReplyToken, linebot.NewTextMessage(tmpReplyMessage)).Do(); err != nil {
		log.Fatal(err)
	}

	// 終了
	fmt.Println("*** end")
	return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}
