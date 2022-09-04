# problem1

バックエンドAPIの実装

## 事前準備

### インストールが必要なもの

- docker / docker-compose

## 動作確認

```
$ docker-compose up -d
```

を実行することでローカルのDocker上にサーバが起動

サーバ起動後、以下のURLにアクセスするか、ターミナル上でcurlコマンドを叩くことで動作確認が可能<br><br>
<http://localhost:1323/(APIパス)> <br><br>

```
$ curl -X <メソッド> --location "http://localhost:1323/(APIパス)" \ 
-H "accept: application/json, text/plain, */*"
```

## API仕様

API仕様はSwaggerUIを利用して閲覧

SwaggerUIサーバ起動後以下のURLからSwaggerUIへアクセス

SwaggerUI: <http://localhost:3000/><br>
定義ファイル: `./spec/openapi.yaml`
