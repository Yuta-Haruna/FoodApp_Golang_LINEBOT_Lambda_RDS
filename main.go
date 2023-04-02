package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/go-sql-driver/mysql"
	"github.com/line/line-bot-sdk-go/linebot"
)

// API Gatewayから受け取ったevents.APIGatewayProxyRequestのBody（JSON）をパースする
// https://app.quicktype.io/　に、実際に受け取ってたJSONメッセージを張り付けて、コード自動生成。

// 接続先DB情報
var db *sql.DB

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
	fmt.Println("*** linebotConnect")

	if err != nil {
		fmt.Println(err)
	}

	// DBから取得した値を格納する
	var resultCal *int
	fmt.Println(myLineRequest.Events[0].Message.Text)

	// DBから食品情報を取得する
	dbErr := db.QueryRow("SELECT calQuantityKcal FROM foodAppDB.foodCalorieInformation WHERE Name LIKE ?", myLineRequest.Events[0].Message.Text).Scan(&resultCal)
	fmt.Println("*** dbErr")
	if dbErr != nil {
		fmt.Println("データベース接続失敗")
		panic(err.Error())
	} else {
		fmt.Println("データベース接続成功")
	}

	// リプライ実施
	fmt.Println("*** reply")
	var tmpReplyMessage string
	if resultCal == nil {
		tmpReplyMessage = "🦜オウムさん：" + myLineRequest.Events[0].Message.Text + "、そいつのカロリーは不明だ。食い物だったら登録してくれ😎"
	} else {
		tmpReplyMessage = "🦜オウムさん：" + myLineRequest.Events[0].Message.Text + "、そいつは" + strconv.Itoa(*resultCal) + "kcalだ!!"
	}
	if _, err = bot.ReplyMessage(myLineRequest.Events[0].ReplyToken, linebot.NewTextMessage(tmpReplyMessage)).Do(); err != nil {
		log.Fatal(err)
	}

	// データベースの接続を閉じる
	db.Close()

	// 終了
	fmt.Println("*** end")
	return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}

func init() {
	// データベースの接続情報
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	// dbName := os.Getenv("DB_NAME")
	// dbName := os.Getenv("foodAppDB")

	// データベースに接続するための文字列
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/foodAppDB", user, password, host, port)
	fmt.Println("*** dataSourceName : " + dataSourceName)
	fmt.Println("*** DB")
	// データベースに接続する
	var err error
	db, err = sql.Open("mysql", dataSourceName)

	fmt.Println(db)
	if err != nil {
		fmt.Println("error -- ")
		fmt.Println(err)
		fmt.Println("error -- ")
	}
	fmt.Println("*** DB : Start")

	// データベースの接続を確認する
	dbErr := db.Ping()
	if dbErr != nil {
		fmt.Println("dbErr -- ")
		fmt.Println(dbErr)
		fmt.Println("dbErr -- ")

	}
	fmt.Println("Connected!")
	fmt.Println("*** DB END")

}
