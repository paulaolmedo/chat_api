FROM golang:1.16-alpine as builder
RUN apk add --update gcc musl-dev
WORKDIR $GOPATH/chat_api
COPY . .
ENV GO111MODULE=on
RUN GOOS=linux go build cmd/server.go

FROM alpine:latest
ENV GO_REST_API /home
WORKDIR $GO_REST_API
COPY --from=builder /go/chat_api/server .
COPY --from=builder /go/chat_api/cmd/chat.properties cmd/chat.properties
COPY --from=builder /go/chat_api/pkg/auth/auth0.properties pkg/auth/auth0.properties
CMD ./server