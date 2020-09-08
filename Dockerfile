FROM golang:1.15-alpine as go-build

RUN mkdir -p /go/src/app
WORKDIR /go/src/app

RUN apk add --no-cache git gcc musl-dev

RUN go get github.com/NYTimes/gziphandler && \
    go get github.com/lib/pq && \
    go get github.com/golang-migrate/migrate && \
    go get github.com/hashicorp/go-multierror && \
    go get github.com/gorilla/sessions

COPY ./backend/ /go/src/app

RUN go build -o run-app main.go



FROM node:12.18-alpine as yarn-build

WORKDIR /usr
COPY ./frontend/package.json ./frontend/yarn.lock ./frontend/tsconfig.json ./frontend/config-overrides.js /usr/
RUN yarn

COPY ./frontend/src /usr/src
COPY ./frontend/public /usr/public

RUN yarn build



FROM alpine:3.7

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=yarn-build /usr/build ./static
COPY --from=go-build /go/src/app/run-app .


CMD ["./run-app", "--production"]
