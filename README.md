# spreadsheet-api-on-golang
Spread Sheet API v4をgolangで実行する

## usage
src/main/main.go

const spreadsheetId = "1XVXj4VCvnaI_QBCJFO_ikJfLKqXr9vhofYjCfQsVLbk"
const sheetId = int64(0)
spreadsheetId、sheetIdを自身がアクセス出来るスプレッドシートのIDに変えてください。

ともに対象のスプレッドシートを開いた際のURLに記載されています。
例)https://docs.google.com/spreadsheets/d/1XVXj4VCvnaI_QBCJFO_ikJfLKqXr9vhofYjCfQsVLbk/edit#gid=0

```console
cd src/main
go run main.go gclient.go sheetsService.go 
```

**初回起動時のみ**以下のようにURLが表示されるのでブラウザでアクセスして認可を行ってください。
認可するとコードが出てくるのでそれをコンソールに入力してください。
```conosle
```



