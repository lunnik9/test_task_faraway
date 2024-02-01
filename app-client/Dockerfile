FROM golang:1.18 as builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -v -o myclient

FROM alpine:3.14
WORKDIR /root/

COPY --from=builder /app/myclient .

CMD ["./myclient"]