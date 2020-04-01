package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/iterator"
)

// RequestData datastoreに登録するデータの型
type RequestData struct {
	SessionID string
	Firstname string // 名前
	Lastname  string // 苗字
	State     int
}

// Datastore datastoreを操作するための情報
type Datastore struct {
	projectID string         // datastoreのProjectID
	nameKey   string         // 登録するキー
	key       *datastore.Key // datastore Key
}

// CreateDatastore datastore操作を生成する
func CreateDatastore(projectID string, nameKey string) *Datastore {
	return &Datastore{
		projectID: projectID,
		nameKey:   nameKey,
		key:       nil,
	}
}

// Put エンティティを登録する
func (d *Datastore) Put(ctx context.Context, req *RequestData) error {
	//	ctx := context.Background()
	client, err := datastore.NewClient(ctx, d.projectID)
	if err != nil {
		return err
	}
	defer client.Close()

	//	key := datastore.NameKey("RequestData", "", nil)
	_, err = client.Put(ctx, d.key, req)
	return err
}

// Get sessionIDが一致するエンティティを取り出す
func (d *Datastore) Get(ctx context.Context, req *RequestData) (*datastore.Key, error) {
	client, err := datastore.NewClient(ctx, d.projectID)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	q := datastore.NewQuery("RequestData").Filter("SessionID =", req.SessionID)
	it := client.Run(ctx, q)
	key, err := it.Next(req)  // rを更新してキーを取得する
	if err == iterator.Done { // datastoreにデータがない場合はキーを発行する
		key = datastore.NameKey("RequestData", "", nil)
		err = nil
	}
	return key, err
}

func example() {
	reqdata := RequestData{
		SessionID: "testID",
		Firstname: "hoge",
		Lastname:  "fuga",
		State:     0, // INITと書きたい
	}

	projectID := os.Getenv("MY_PROJECT_ID")
	nameKey := "RequestData"

	ds := CreateDatastore(projectID, nameKey)

	ctx := context.Background()
	if err := ds.Put(ctx, &reqdata); err != nil {
		log.Print(err)
	}

	_, err := ds.Get(ctx, &reqdata)
	if err != nil {
		return
	}
	fmt.Println(reqdata)
}
