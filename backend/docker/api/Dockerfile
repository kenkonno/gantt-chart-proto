FROM golang:1.19

# コンテナ作業ディレクトリの変更
WORKDIR /go/src/github.com/kenkonno/gantt-chart-proto/backend/api/

#TODO:環境変数問題
COPY ./.env ../
COPY ./api/ ./
COPY ./models/ ../models/
COPY ./repository/ ../repository/

RUN go get github.com/githubnemo/CompileDaemon
RUN go install github.com/githubnemo/CompileDaemon

RUN go mod download
CMD ["CompileDaemon","-command=./api","-include=*.go","-polling"]

