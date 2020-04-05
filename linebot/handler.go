package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"cloud.google.com/go/datastore"
	"github.com/line/line-bot-sdk-go/linebot"
)

// Roster is 名簿クラス
type Roster struct {
	//projectID string
	regist   *Registration
	question *Questionnaire
}

// NewRoster creates roster
func NewRoster(projectID string) *Roster {
	return &Roster{
		regist:   NewRegistration(projectID, "RegistrationData"),
		question: NewQuestionnaire(projectID, "QuestionnaireData"),
	}
}

// Handler handles linebot
type Handler struct {
	client *linebot.Client
	roster *Roster
}

// NewHandler creates handler
func NewHandler(client *linebot.Client, roster *Roster) *Handler {
	return &Handler{
		client: client,
		roster: roster,
	}
}

func (h *Handler) callback(w http.ResponseWriter, req *http.Request) {
	events, err := h.client.ParseRequest(req)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}
	for _, event := range events {
		h.LinebotMessageExec(event)
	}
}

// LinebotMessageExec Linebotへのメッセージを実行する
func (h *Handler) LinebotMessageExec(event *linebot.Event) {
	if event.Type != linebot.EventTypeMessage {
		return
	}

	switch message := event.Message.(type) {
	case *linebot.TextMessage:
		h.replyMessageExec(event, message)

	case *linebot.StickerMessage:
		replyMessage := fmt.Sprintf(
			"sticker id is %s, stickerResourceType is ...", message.StickerID)
		if _, err := h.client.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
			log.Print(err)
		}
	}
}

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

// 参考文献
// https://blog.kazu634.com/labs/golang/2019-02-23-line-sdk-go/
func (h *Handler) replyMessageExec(event *linebot.Event, message *linebot.TextMessage) {
	//datastoreから入力データを取得する
	var question QuestionnaireData
	question.SessionID = event.Source.UserID

	ctx := context.Background()
	query := datastore.NewQuery("QuestionnaireData").Filter("SessionID =", event.Source.UserID)
	key, err := h.roster.question.Get(ctx, query, &question)
	if err != nil {
		log.Print("Get失敗", err)
		return
	}

	_ = h.Exec(event, &question, message.Text)

	err = h.roster.question.Put(ctx, key, &question)
	if err != nil {
		log.Print("Put失敗", err)
	}
}

// LinebotTextMessage Linebotへのテキストメッセージを送信する
func (h *Handler) LinebotTextMessage(event *linebot.Event, message string) error {
	resp := linebot.NewTextMessage(message)
	_, err := h.client.ReplyMessage(event.ReplyToken, resp).Do()
	if err != nil {
		log.Print(err)
	}
	return err
}

// Exec Linebotへの入力指示を順に実行する
func (h *Handler) Exec(event *linebot.Event, request *QuestionnaireData, text string) error {
	var err error
	switch request.State {
	case INIT:
		err = h.LinebotTextMessage(event, "苗字を入れてください")
		request.State = LASTNAME

	case LASTNAME:
		request.Lastname = text
		err = h.LinebotTextMessage(event, "名前を入れてください")
		request.State = FIRSTNAME

	case FIRSTNAME:
		request.Firstname = text
		resp := linebot.NewTemplateMessage(
			"this is a confirm template",
			linebot.NewConfirmTemplate(
				"Are you sure?",
				linebot.NewMessageAction("Yes", "yes"),
				linebot.NewMessageAction("No", "no"),
			),
		)
		_, err := h.client.ReplyMessage(event.ReplyToken, resp).Do()
		if err != nil {
			log.Print(err)
		}
		request.State = RESULT

	case RESULT:
		// datastoreから一時情報を削除
		if text == "yes" {
			// datastoreに保存
			regist := RegistrationData{
				Firstname: request.Firstname,
				Lastname:  request.Lastname,
			}
			h.roster.regist.Put(context.Background(), datastore.NameKey("RegistrationData", "", nil), &regist)
			// メッセージを送信
			err = h.LinebotTextMessage(event, "登録しました")
		} else {
			// メッセージを送信
			err = h.LinebotTextMessage(event, "登録をキャンセルしました")
		}
		request.State = INIT
	}
	return err
}
