# base go image
FROM golang:1.18-alpine as builder

RUN mkdir /cmd

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o _build/urlFrontierApp ./cmd

RUN chmod +x ./_build/urlFrontierApp

# _build a tiny docker image
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/_build/urlFrontierApp /app

CMD ["/app/urlFrontierApp"]