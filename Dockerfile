FROM golang:alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GOOS linux

RUN apk update --no-cache && apk add --no-cache tzdata

WORKDIR /build

ADD go.mod .
ADD go.sum .

RUN go mod download

COPY . ./

RUN go build -ldflags="-s -w" -o url-shortener ./cmd/

FROM alpine

RUN apk update --no-cache && apk add --no-cache ca-certificates

COPY --from=builder /build/url-shortener /app/url-shortener

WORKDIR /app

EXPOSE 8080

CMD ["/app/url-shortener"]