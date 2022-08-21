FROM golang:1.16.5-alpine3.12 as builder

WORKDIR /app

COPY . .

RUN go build -ldflags="-w -s" -o main

FROM alpine

COPY --from=builder /app/main  /app/main
COPY --from=builder /app/config /app/config

WORKDIR /app

ENTRYPOINT ["./main"]
