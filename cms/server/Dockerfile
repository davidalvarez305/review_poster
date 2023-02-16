# syntax=docker/dockerfile:1

FROM golang:1.19-alpine

# This is done so that it matches with several paths in ENV file.
WORKDIR /content_go/server

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

ENV APP_ENV=production

RUN go build -o /exec

CMD [ "/exec" ]