linebot
====

linebotを使ったアンケート収集アプリのサンプル。苗字と名前を入力してもらい登録する。

## Description
- linebotからwebhookで接続するアプリ
- 苗字と名前を別々に入力してもらい、datastoreに登録する
- linebotのユーザーIDを利用してユーザーごとに入力を受け付けることができる

## Usage
linebot＋GCPで動作させる設定を紹介します。
### 1.LINE Developerの設定
1. LINE DeveloperでMessaging APIを作成する
2. チャネルシークレットとチャネルアクセストークンを取得する

### 2.GCPの設定
1. GCPでプロジェクトを作る
2. プロジェクトIDを取得する

### 3.app.yamlの設定
1. PROJECT_IDにGAEのプロジェクトIDを設定する
2. CHANNEL_SECRETにMessaging APIのチャネルシークレットを設定する
3. CHANNEL_TOKENにMessaging APIのチャネルアクセストークンを設定する

### 4.GCPにデプロイ
1. プロジェクトをGCPにデプロイする
2. GCP上で動作しているURLを取得する

### 5.LINE Developerの設定
1. LINE DeveloperでWebhook URLを上記のURLを登録する
