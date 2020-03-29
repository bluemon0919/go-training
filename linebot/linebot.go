package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/line/line-bot-sdk-go/linebot"
)

const (
	INIT = iota // iotaの初期値は0
	FIRSTNAME
	LASTNAME
	RESULT
)

type Request struct {
	firstname string // 名前
	lastname  string // 苗字
	state     int
}
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
	ds := CreateDatastore(projectID, "RequestData")
	/*
		reqdatas, err := ds.Get(event.Source.UserID)
		var req Request
		if len(reqdatas) > 0 {
			log.Println("データが見つかりました")
			req = Request{
				firstname: reqdatas[0].Firstname,
				lastname:  reqdatas[0].Lastname,
				state:     reqdatas[0].State,
			}
		}
	*/

	var requestData RequestData
	requestData.SessionID = event.Source.UserID
	ctx := context.Background()
	key, err := ds.Get(ctx, &requestData)
	if err != nil {
		log.Println("失敗")
		log.Print(err)
		return
	}
	ds.key = key
	request := Request{
		firstname: requestData.Firstname,
		lastname:  requestData.Lastname,
		state:     requestData.State,
	}
	reqManager := CreateRequestManager(event, request)
	_ = reqManager.Exec(message.Text)

	requestData = RequestData{
		SessionID: event.Source.UserID,
		Firstname: reqManager.request.firstname,
		Lastname:  reqManager.request.lastname,
		State:     reqManager.request.state,
	}
	err = ds.Put(ctx, &requestData)
	if err != nil {
		log.Print(err)
	}
}

func LinebotTextMessage(event *linebot.Event, message string) error {
	resp := linebot.NewTextMessage(message)
	_, err := bot.ReplyMessage(event.ReplyToken, resp).Do()
	if err != nil {
		log.Print(err)
	}
	return err
}

func CreateRequestManager(e *linebot.Event, request Request) *RequestManager {
	return &RequestManager{
		request: request,
		event:   e,
	}
}

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
