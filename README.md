# spreadsheet-api-on-golang
Spread Sheet API v4を使ってgopherをスプレッドシートに描く

## usage
src/main/main.go

const spreadsheetId = "1XVXj4VCvnaI_QBCJFO_ikJfLKqXr9vhofYjCfQsVLbk"
const sheetId = int64(0)
spreadsheetId、sheetIdを自身がアクセス出来るスプレッドシートのIDに変えてください。

ともに対象のスプレッドシートを開いた際のURL内に含まれています。
例)https://docs.google.com/spreadsheets/d/1XVXj4VCvnaI_QBCJFO_ikJfLKqXr9vhofYjCfQsVLbk/edit#gid=0

```console 
go get golang.org/x/net/context golang.org/x/oauth2 google.golang.org/api/sheets/v4 google.golang.org/cloud/compute/metadata

cd src/main
go run main.go gclient.go sheetsService.go 
```

**初回起動時のみ**以下のようにURLが表示されるのでブラウザでアクセスして認可を行ってください。
認可するとコードが出てくるのでそれをコンソールに入力してください。

![step1](https://raw.githubusercontent.com/howdy39/spreadsheets-api-on-golang/master/screenshots/step1.png)

![step2](https://raw.githubusercontent.com/howdy39/spreadsheets-api-on-golang/master/screenshots/step2.png)

![step3](https://raw.githubusercontent.com/howdy39/spreadsheets-api-on-golang/master/screenshots/step3.png)

![step4](https://raw.githubusercontent.com/howdy39/spreadsheets-api-on-golang/master/screenshots/step4.png)
