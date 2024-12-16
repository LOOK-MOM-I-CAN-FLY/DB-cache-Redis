FROM golang:1.20 as builder

WORKDIR /app

COPY . .

RUN go mod tidy && go build -o main .

FROM debian:bullseye-slim

WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/config ./config

EXPOSE 8080
CMD ["./main"]
