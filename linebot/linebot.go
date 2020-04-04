package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/datastore"
	"github.com/line/line-bot-sdk-go/linebot"
)

const (
	// INIT 初期状態
	INIT = iota
	// FIRSTNAME 名前入力中
	FIRSTNAME
	// LASTNAME 苗字入力中
	LASTNAME
	// RESULT 結果確認中
	RESULT
)

// Request 入力情報を管理する
type Request struct {
	firstname string // 名前
	lastname  string // 苗字
	state     int
}

// RequestManager Linebotへの要求を管理する
type RequestManager struct {
	request Request
	event   *linebot.Event
}

var bot *linebot.Client

func init() {
	var err error
	bot, err = linebot.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}
}

// LinebotMessageExec Linebotへのメッセージを実行する
func LinebotMessageExec(event *linebot.Event) {
	if event.Type != linebot.EventTypeMessage {
		return
	}

	switch message := event.Message.(type) {
	case *linebot.TextMessage:
		//log.Println(event.Source.UserID)
		replyMessageExec(event, message)

	case *linebot.StickerMessage:
		replyMessage := fmt.Sprintf(
			"sticker id is %s, stickerResourceType is ...", message.StickerID)
		if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
			log.Print(err)
		}
	}
}

// 参考文献
// https://blog.kazu634.com/labs/golang/2019-02-23-line-sdk-go/
func replyMessageExec(event *linebot.Event, message *linebot.TextMessage) {
	//datastoreから入力データを取得する
	projectID := os.Getenv("PROJECT_ID")
	fuga := New(projectID, "RequestData")
	//	fuga.data.SessionID =

	var requestData RequestData
	requestData.SessionID = event.Source.UserID
	ctx := context.Background()
	query := datastore.NewQuery("RequestData").Filter("SessionID =", requestData.SessionID)
	key, err := fuga.Get(ctx, query)
	if err != nil {
		log.Print("Get失敗", err)
		return
	}
	fuga.data.SessionID = event.Source.UserID // ここだけ処理がまとまっていない
	request := Request{
		firstname: fuga.data.Firstname,
		lastname:  fuga.data.Lastname,
		state:     fuga.data.State,
	}
	reqManager := CreateRequestManager(event, request)
	_ = reqManager.Exec(message.Text)

	fuga.data = RequestData{
		SessionID: event.Source.UserID,
		Firstname: reqManager.request.firstname,
		Lastname:  reqManager.request.lastname,
		State:     reqManager.request.state,
	}
	err = fuga.Put(ctx, key)
	if err != nil {
		log.Print("Put失敗", err)
	}
}

// LinebotTextMessage Linebotへのテキストメッセージを送信する
func LinebotTextMessage(event *linebot.Event, message string) error {
	resp := linebot.NewTextMessage(message)
	_, err := bot.ReplyMessage(event.ReplyToken, resp).Do()
	if err != nil {
		log.Print(err)
	}
	return err
}

// CreateRequestManager RequestManagerを生成する
func CreateRequestManager(e *linebot.Event, request Request) *RequestManager {
	return &RequestManager{
		request: request,
		event:   e,
	}
}

// Exec Linebotへの入力指示を順に実行する
func (m *RequestManager) Exec(text string) error {
	var err error
	switch m.request.state {
	case INIT:
		err = LinebotTextMessage(m.event, "苗字を入れてください")
		m.request.state = LASTNAME

	case LASTNAME:
		m.request.lastname = text
		err = LinebotTextMessage(m.event, "名前を入れてください")
		m.request.state = FIRSTNAME

	case FIRSTNAME:
		m.request.firstname = text
		resp := linebot.NewTemplateMessage(
			"this is a confirm template",
			linebot.NewConfirmTemplate(
				"Are you sure?",
				linebot.NewMessageAction("Yes", "yes"),
				linebot.NewMessageAction("No", "no"),
			),
		)
		_, err := bot.ReplyMessage(m.event.ReplyToken, resp).Do()
		if err != nil {
			log.Print(err)
		}
		m.request.state = RESULT

	case RESULT:
		// datastoreから一時情報を削除
		if text == "yes" {
			// datastore2に保存
			// メッセージを送信
			err = LinebotTextMessage(m.event, "登録しました")
		} else {
			// メッセージを送信
			err = LinebotTextMessage(m.event, "登録をキャンセルしました")
		}
		m.request.state = INIT
	}
	return err
}
