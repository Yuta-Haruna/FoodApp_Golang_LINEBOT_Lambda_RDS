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

// 接続先DB情報
var db *sql.DB

// 取得したLINEの入力内容JSONを整形
func UnmarshalLineRequest(data []byte) (LineRequest, error) {
	var r LineRequest
	err := json.Unmarshal(data, &r)
	return r, err
}

// JSON構造体
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

	fmt.Println(myLineRequest.Events[0].Message.Text)

	// データベースの接続情報
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// データベースに接続するための文字列
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, dbName)
	fmt.Println("*** dataSourceName : " + dataSourceName)
	fmt.Println("*** DB")
	// データベースに接続する
	var dbErr error
	db, dbErr = sql.Open("mysql", dataSourceName)

	// 処理の最後にデータベースの接続を閉じる
	defer func() {
		if cerr := db.Close(); cerr != nil {
			log.Printf("Failed to close database: %v", cerr)
		}
	}()

	fmt.Println(db)
	if dbErr != nil {
		fmt.Println("error -- ")
		fmt.Println(dbErr)
		fmt.Println("error -- ")
	}
	fmt.Println("*** DB : Start")

	// データベースの接続を確認する
	dbErr = db.Ping()
	if dbErr != nil {
		fmt.Println("dbErr -- ")
		fmt.Println(dbErr)
		fmt.Println("dbErr -- ")

	}
	fmt.Println("Connected!")
	fmt.Println("*** DB END")

	// 最終表記
	var tmpReplyMessage string

	// DBから食品情報を取得する
	query := "SELECT calQuantityKcal FROM foodCalorieInformation WHERE Name LIKE ? LIMIT 1"
	rows, dbErr := db.Query(query, myLineRequest.Events[0].Message.Text)

	fmt.Println("*** dbErr")
	if dbErr != nil {
		fmt.Println("データベース接続失敗")
		fmt.Println(dbErr)
	}
	// 処理終了後 SQLの停止
	defer rows.Close()

	// クエリの結果を処理する
	if !rows.Next() {
		// クエリの結果が空だった場合の処理
		fmt.Println("あああああ")
		tmpReplyMessage = "🦜オウムさん：" + myLineRequest.Events[0].Message.Text + "、そいつのカロリーは不明だ。食い物だったら登録してくれ😎"
	} else {
		// クエリの結果があった場合の処理
		// var column1 string
		var column1 int
		if err := rows.Scan(&column1); err != nil {
			fmt.Println("データベース接続成功 : クエリ失敗")
			fmt.Println(err)
		}
		fmt.Printf("column1: %d", column1)
		tmpReplyMessage = "🦜オウムさん：" + myLineRequest.Events[0].Message.Text + "、そいつは" + strconv.Itoa(column1) + "kcalだ!!"

	}

	// リプライ実施
	fmt.Println("*** reply")
	fmt.Println(tmpReplyMessage)

	fmt.Println(myLineRequest.Events[0].ReplyToken)

	// LINEBOTへの連絡内容を設定
	replyMessage, err := bot.ReplyMessage(myLineRequest.Events[0].ReplyToken, linebot.NewTextMessage(tmpReplyMessage)).Do()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(replyMessage)
	}

	// 終了
	fmt.Println("*** end")
	return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}
