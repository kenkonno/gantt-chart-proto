FROM golang:1.19

# コンテナ作業ディレクトリの変更
WORKDIR /go/src/github.com/kenkonno/gantt-chart-proto/backend/models/
COPY ./models/go.mod ./models/go.sum ./

RUN go mod download