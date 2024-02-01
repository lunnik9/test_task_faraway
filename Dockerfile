FROM golang:1.18 AS builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o myapp

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/myapp .

CMD ["./myapp"]
