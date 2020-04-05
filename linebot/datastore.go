package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/iterator"
)

// DatastoreOperator datastore操作インターフェース
type DatastoreOperator interface {
	Put(ctx context.Context, key *datastore.Key) error
	Get(ctx context.Context, query *datastore.Query) (*datastore.Key, error)
}

// RegistrationData 登録する内容
type RegistrationData struct {
	Firstname string
	Lastname  string
}

// Registration ユーザーへのアンケート結果を登録する
type Registration struct {
	projectID  string
	entityType string           // datastore entity type
	entity     RegistrationData // datastore entity
}

// NewRegistration 登録を実行するクラス
func NewRegistration(projectID string, entityType string) *Registration {
	return &Registration{
		projectID:  projectID,
		entityType: entityType,
		entity:     RegistrationData{},
	}
}

// Put エンティティに登録する
func (r *Registration) Put(ctx context.Context, key *datastore.Key) error {
	client, err := datastore.NewClient(ctx, r.projectID)
	if err != nil {
		return err
	}
	defer client.Close()

	_, err = client.Put(ctx, key, &r.entity)
	return err
}

// Get sessionIDが一致するエンティティを取り出す
func (r *Registration) Get(ctx context.Context, query *datastore.Query) (*datastore.Key, error) {
	client, err := datastore.NewClient(ctx, r.projectID)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	it := client.Run(ctx, query)
	key, err := it.Next(&r.entity)
	if err == iterator.Done { // datastoreにデータがない場合はキーを発行する
		key = datastore.NameKey("RegistrationData", "", nil)
		err = nil
	}
	return key, err
}

// QuestionnaireData アンケートを取る内容. datastore entityに該当する
type QuestionnaireData struct {
	SessionID string
	Firstname string
	Lastname  string
	State     int
}

// Questionnaire ユーザーからアンケートをとる
type Questionnaire struct {
	projectID  string
	entityType string            // datastore entity type
	entity     QuestionnaireData // datastore entity
}

// NewQuestionnaire アンケートクラス生成する
func NewQuestionnaire(projectID string, entityType string) *Questionnaire {
	return &Questionnaire{
		projectID:  projectID,
		entityType: entityType,
		entity:     QuestionnaireData{},
	}
}

// Put エンティティに登録する
func (q *Questionnaire) Put(ctx context.Context, key *datastore.Key) error {
	client, err := datastore.NewClient(ctx, q.projectID)
	if err != nil {
		return err
	}
	defer client.Close()

	_, err = client.Put(ctx, key, &q.entity)
	return err
}

// Get sessionIDが一致するエンティティを取り出す
func (q *Questionnaire) Get(ctx context.Context, query *datastore.Query) (*datastore.Key, error) {
	client, err := datastore.NewClient(ctx, q.projectID)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	it := client.Run(ctx, query)
	key, err := it.Next(&q.entity)
	if err == iterator.Done { // datastoreにデータがない場合はキーを発行する
		key = datastore.NameKey("QuestionnaireData", "", nil)
		err = nil
	}
	return key, err
}

func exampleQuestionnaire() {
	entity := QuestionnaireData{
		SessionID: "testID",
		Firstname: "hoge",
		Lastname:  "Questionnaire",
		State:     0, // INITと書きたい
	}

	projectID := os.Getenv("PROJECT_ID")
	Questionnaire := NewQuestionnaire(projectID, " QuestionnaireData")
	var err error
	var key *datastore.Key
	// 登録済みのデータを取得する。データがない場合は空のデータと新しいキーが発行される。
	ctx := context.Background()
	query := datastore.NewQuery(" QuestionnaireData").Filter("SessionID =", entity.SessionID)
	if key, err = Questionnaire.Get(ctx, query); err != nil {
		return
	}

	if err = Questionnaire.Put(ctx, key); err != nil {
		log.Print(err)
	}

	if _, err = Questionnaire.Get(ctx, query); err != nil {
		return
	}
	fmt.Println(entity)
}
