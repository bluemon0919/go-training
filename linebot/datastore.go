package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/iterator"
)

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
func (f *Questionnaire) Put(ctx context.Context, key *datastore.Key) error {
	client, err := datastore.NewClient(ctx, f.projectID)
	if err != nil {
		return err
	}
	defer client.Close()

	_, err = client.Put(ctx, key, &f.entity)
	return err
}

// Get sessionIDが一致するエンティティを取り出す
func (f *Questionnaire) Get(ctx context.Context, query *datastore.Query) (*datastore.Key, error) {
	client, err := datastore.NewClient(ctx, f.projectID)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	it := client.Run(ctx, query)
	key, err := it.Next(&f.entity)
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
