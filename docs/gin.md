https://gin-gonic.com/ja/docs/quickstart/

gin.Defaultでルーティング
r.GETでルーティング

c.JSONでjsonレスポンス

git.H構造体でkeyValueのレスポンスができる。

## HMR
go install https://github.com/air-verse/air@latest
airをいれてgo initすればOK
tomlが作られる。
デフォルト設定でOK

あとはairコマンド実行するとmain.goが動く

## CRUD実装
Router/Controller/Service/Repository/resourceの構成とする
各層は必ずIFに依存することで依存性の逆転を行う

