// Copyright 2016 LINE Corporation
//
// LINE Corporation licenses this file to you under the Apache License,
// version 2.0 (the "License"); you may not use this file except in compliance
// with the License. You may obtain a copy of the License at:
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/linebot"
)

var bot *linebot.Client

func main() {
	var err error
	bot, err = linebot.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Setup HTTP Server for receiving requests from LINE platform
	http.HandleFunc("/callback", callback)
	// This is just sample code.
	// For actual use, you must support HTTPS by using `ListenAndServeTLS`, a reverse proxy or something else.
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		log.Fatal(err)
	}
}

// UserSession managed user session
type UserSession struct {
	UserID string
	Count  int
}

var us UserSession

// Start starts session
func (us *UserSession) Start(userID string) int {
	if us.UserID == "" {
		us.UserID = userID
	}
	if us.UserID != userID {
		return -1
	}
	return us.Count
}

// Close closes session
func (us *UserSession) Close(userID string) {
	if us.UserID != userID {
		return
	}
	us.Count++
	if us.Count > 3 {
		us.Count = 0
	}
}

func callback(w http.ResponseWriter, req *http.Request) {
	events, err := bot.ParseRequest(req)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}
	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				log.Println(event.Source.UserID)
				replyMessageExec(event, message)

			case *linebot.StickerMessage:
				replyMessage := fmt.Sprintf(
					"sticker id is %s, stickerResourceType is ...", message.StickerID)
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
					log.Print(err)
				}
			}
		}
	}
}

// 参考文献
// https://blog.kazu634.com/labs/golang/2019-02-23-line-sdk-go/
func replyMessageExec(event *linebot.Event, message *linebot.TextMessage) {
	sessionCount := us.Start(event.Source.UserID)

	if "" != message.Text {
		switch sessionCount {
		case 0:
			resp := linebot.NewTextMessage("苗字を入れてください")
			_, err := bot.ReplyMessage(event.ReplyToken, resp).Do()
			if err != nil {
				log.Print(err)
			}
		case 1:
			resp := linebot.NewTextMessage("名前を入れてください")
			_, err := bot.ReplyMessage(event.ReplyToken, resp).Do()
			if err != nil {
				log.Print(err)
			}
		case 2:
			// 入力の確認
			resp := linebot.NewTemplateMessage(
				"this is a confirm template",
				linebot.NewConfirmTemplate(
					"Are you sure?",
					linebot.NewMessageAction("Yes", "yes"),
					linebot.NewMessageAction("No", "no"),
				),
			)
			_, err := bot.ReplyMessage(event.ReplyToken, resp).Do()
			if err != nil {
				log.Print(err)
			}
		case 3:
			if "yes" == message.Text {
				resp := linebot.NewTextMessage("登録しました")
				_, err := bot.ReplyMessage(event.ReplyToken, resp).Do()
				if err != nil {
					log.Print(err)
				}
			}
		}
	}
	us.Close(event.Source.UserID)
}