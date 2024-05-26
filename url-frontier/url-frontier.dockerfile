# base go image
FROM golang:1.18-alpine as builder

RUN mkdir /cmd

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go _build -o ./_build/urlLoaderApp ./cmd

RUN chmod +x ./_build/urlLoaderApp

# _build a tiny docker image
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/_build/urlLoaderApp /app

CMD ["/app/urlLoaderApp"]