FROM golang:1.19

# コンテナ作業ディレクトリの変更
WORKDIR /go/src/github.com/kenkonno/gantt-chart-proto/backend/migration/
COPY ./migration/go.mod ./migration/go.sum ./
COPY ./models/go.mod ./models/go.sum ../models/
COPY ./repository/go.mod ./repository/go.sum ../repository/

RUN go mod download