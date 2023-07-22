FROM golang:1.19

# コンテナ作業ディレクトリの変更
WORKDIR /go/src/github.com/kenkonno/twitch_clips_project/backend/api/
COPY ./api/go.mod ./api/go.sum ./
COPY ./models/go.mod ./models/go.sum ../models/
COPY ./repository/go.mod ./repository/go.sum ../repository/
RUN go get github.com/githubnemo/CompileDaemon
RUN go install github.com/githubnemo/CompileDaemon

RUN go mod download