# base go image
FROM golang:1.18-alpine as builder

RUN mkdir /cmd

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o bin/urlFrontierApp ./cmd

RUN chmod +x ./bin/urlFrontierApp

# build a tiny docker image
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/bin/urlFrontierApp /app

CMD ["/app/urlFrontierApp"]