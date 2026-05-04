basic-goという名前、localというタグ名でimageを作成
```bash
docker build -t basic-go:local .
```

8080ポートでimageをrun
```bash
docker run -p 8080:8080 basic-go:local
```

ciで以下のように定義する。CI_COMMIT_SHORT_SHAでコミットのハッシュをタグにできる
タグ名からlatestを作成する。2回目以降はlatestは上書きされる。これは常に最新のタグが簡単に追えるようにする
```bash
- docker build -t basic-go:$CI_COMMIT_SHORT_SHA .
- docker tag basic-go:$CI_COMMIT_SHORT_SHA basic-go:latest
```
