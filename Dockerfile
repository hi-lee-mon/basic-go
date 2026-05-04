# 動作確認
## docker build -t basic-go:local .
## docker run -p 8080:8080 basic-go:local

# alpine上にgoをインストール済みのものを使用
FROM golang:1.25-alpine AS builder
# イメージの中にappフォルダを作成してそこをルートとする
WORKDIR /app
# go.modをルートにコピー(/app/go.modになる)
COPY go.mod ./
RUN go mod download
COPY . .
# Pure Goバイナリ生成、Linux用バイナリ、serverというファイル名で生成、.をビルド対象に指定
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

# goはbuilderステージでバイナリにビルド済みのためコンパイラ不要のlinuxイメージのみを選択
FROM alpine:3.21
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]