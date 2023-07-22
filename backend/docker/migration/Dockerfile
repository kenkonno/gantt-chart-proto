FROM golang:1.19

# コンテナ作業ディレクトリの変更
WORKDIR /go/src/github.com/kenkonno/twitch_clips_project/backend/models/
COPY ./models/go.mod ./models/go.sum ./

RUN go mod download