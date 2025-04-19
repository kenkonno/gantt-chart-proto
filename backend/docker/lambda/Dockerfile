FROM golang:1.19-alpine AS build_base
RUN apk add --no-cache git

# コンテナ作業ディレクトリの変更
WORKDIR /go/src/github.com/kenkonno/gantt-chart-proto/backend/api/

COPY ./api/ ./
COPY ./models/ ../models/
COPY ./repository/ ../repository/

RUN go mod download

RUN GOOS=linux CGO_ENABLED=0 go build -o bootstrap .
FROM alpine:3.9
RUN apk add ca-certificates
COPY --from=public.ecr.aws/awsguru/aws-lambda-adapter:0.9.0 /lambda-adapter /opt/extensions/lambda-adapter
COPY --from=build_base /go/src/github.com/kenkonno/gantt-chart-proto/backend/api/bootstrap /app/bootstrap

ENV PORT=8000 GIN_MODE=release

EXPOSE 8000

CMD ["/app/bootstrap"]
