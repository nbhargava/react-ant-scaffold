FROM golang:1.11-alpine as go-build

RUN mkdir -p /go/src/app
WORKDIR /go/src/app

RUN apk add --no-cache git gcc musl-dev

RUN go get github.com/NYTimes/gziphandler && \
    go get github.com/lib/pq && \
    go get github.com/golang-migrate/migrate && \
    go get github.com/gorilla/sessions

COPY ./backend/ /go/src/app

RUN go build -o run-app main.go



FROM node:9.9-alpine as yarn-build

WORKDIR /usr
COPY ./package.json ./yarn.lock ./tsconfig.json ./tsconfig.prod.json ./tslint.json ./config-overrides.js /usr/
RUN yarn

COPY ./src /usr/src
COPY ./public /usr/public

RUN yarn build



FROM alpine:3.7

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=yarn-build /usr/build ./static
COPY --from=go-build /go/src/app/run-app .


CMD ["./run-app", "--production"]
