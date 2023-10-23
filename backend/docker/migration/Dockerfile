FROM golang:1.19

# コンテナ作業ディレクトリの変更
WORKDIR /go/src/github.com/kenkonno/gantt-chart-proto/backend/migration/

COPY ./migration/ ./
COPY ./models/ ../models/
COPY ./repository/ ../repository/
RUN go mod download
CMD ["go","run","./main.go"]