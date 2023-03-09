# Initial stage: download modules
FROM golang:1.20-alpine as builder

ENV config=docker

WORKDIR /app

COPY ./ /app

RUN go mod download


# Intermediate stage: Build the binary
FROM golang:1.20-alpine as runner

COPY --from=builder ./app ./app

RUN go get github.com/githubnemo/CompileDaemon

WORKDIR /app
ENV config=docker

EXPOSE 5000

ENTRYPOINT CompileDaemon --build="go build cmd/app/main.go" --command=./main



