# spreadsheet-api-on-golang
Spread Sheet API v4をgolangで実行する

## usage
src/main.go

const spreadsheetId = "1XVXj4VCvnaI_QBCJFO_ikJfLKqXr9vhofYjCfQsVLbk"
const sheetId = int64(0)
spreadsheetId、sheetIdを自身がアクセス出来るスプレッドシートのIDに変えてください。

ともに対象のスプレッドシートを開いた際のURLに記載されています。
例)https://docs.google.com/spreadsheets/d/1XVXj4VCvnaI_QBCJFO_ikJfLKqXr9vhofYjCfQsVLbk/edit#gid=0

```console
cd src/main
go run *.go
```


