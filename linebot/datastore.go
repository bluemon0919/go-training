package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/iterator"
)

type Fuga struct {
	projectID string
	dataType  string
}

type Fugafuga struct {
	Fuga
	data RequestData
}

func New(projectID string, dataType string) *Fugafuga {
	return &Fugafuga{
		Fuga: Fuga{
			projectID: projectID,
			dataType:  dataType,
		},
		data: RequestData{},
	}
}

// Put エンティティを登録する
func (f *Fugafuga) Put(ctx context.Context, key *datastore.Key) error {
	client, err := datastore.NewClient(ctx, f.projectID)
	if err != nil {
		return err
	}
	defer client.Close()

	_, err = client.Put(ctx, key, &f.data)
	return err
}

// Get sessionIDが一致するエンティティを取り出す
func (f *Fugafuga) Get(ctx context.Context, query *datastore.Query) (*datastore.Key, error) {
	client, err := datastore.NewClient(ctx, f.projectID)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	it := client.Run(ctx, query)
	key, err := it.Next(&f.data)
	if err == iterator.Done { // datastoreにデータがない場合はキーを発行する
		key = datastore.NameKey("RequestData", "", nil)
		err = nil
	}
	return key, err
}

func exampleFugafuga() {
	data := RequestData{
		SessionID: "testID",
		Firstname: "hoge",
		Lastname:  "fuga",
		State:     0, // INITと書きたい
	}

	projectID := os.Getenv("PROJECT_ID")
	fuga := New(projectID, "RequestData")
	fmt.Println("hoge:", fuga)
	var err error
	var key *datastore.Key
	// 登録済みのデータを取得する。データがない場合は空のデータと新しいキーが発行される。
	ctx := context.Background()
	query := datastore.NewQuery("RequestData").Filter("SessionID =", data.SessionID)
	if key, err = fuga.Get(ctx, query); err != nil {
		fmt.Println("1: ", err)
		return
	}

	if err = fuga.Put(ctx, key); err != nil {
		log.Print(err)
	}

	if _, err = fuga.Get(ctx, query); err != nil {
		fmt.Println("2:", err)
		return
	}
	fmt.Println(data)
}

// Hoge datastoreの操作情報
type Hoge struct {
	projectID string // ProjectID
	//	key       *datastore.Key // キー
	dataType string      // 登録するデータ型
	data     interface{} // 登録するデータ
}

// CreateHoge creates Hoge
func CreateHoge(projectID string, dataType string) *Hoge {
	return &Hoge{
		projectID: projectID,
		dataType:  dataType,
		data:      nil,
	}
}

// Put エンティティを登録する
func (h *Hoge) Put(ctx context.Context, key *datastore.Key) error {
	fmt.Println("put in")
	//	ctx := context.Background()
	client, err := datastore.NewClient(ctx, h.projectID)
	if err != nil {
		return err
	}
	defer client.Close()

	fmt.Println("key:", key)
	//	key := datastore.NameKey("RequestData", "", nil)
	_, err = client.Put(ctx, key, h.data)
	return err
}

// Get sessionIDが一致するエンティティを取り出す
func (h *Hoge) Get(ctx context.Context, query *datastore.Query) (*datastore.Key, error) {
	fmt.Println("get in")
	client, err := datastore.NewClient(ctx, h.projectID)
	if err != nil {
		fmt.Println("get err:", err)
		return nil, err
	}
	defer client.Close()

	//	q := datastore.NewQuery("RequestData").Filter("SessionID =", d.data.SessionID)
	it := client.Run(ctx, query)
	key, err := it.Next(h)    // rを更新してキーを取得する
	if err == iterator.Done { // datastoreにデータがない場合はキーを発行する
		key = datastore.NameKey("RequestData", "", nil)
		err = nil
	}
	return key, err
}

func exampleHoge() {
	data := RequestData{
		SessionID: "testID",
		Firstname: "hoge",
		Lastname:  "fuga",
		State:     0, // INITと書きたい
	}

	projectID := os.Getenv("PROJECT_ID")
	hoge := CreateHoge(projectID, "RequestData")
	fmt.Println("hoge:", hoge)
	var err error
	var key *datastore.Key
	// 登録済みのデータを取得する。データがない場合は空のデータと新しいキーが発行される。
	ctx := context.Background()
	query := datastore.NewQuery("RequestData").Filter("SessionID =", data.SessionID)
	if key, err = hoge.Get(ctx, query); err != nil {
		fmt.Println("1: ", err)
		return
	}

	if err = hoge.Put(ctx, key); err != nil {
		log.Print(err)
	}

	if _, err = hoge.Get(ctx, query); err != nil {
		fmt.Println("2:", err)
		return
	}
	fmt.Println(data)
}

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
