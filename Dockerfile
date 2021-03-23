FROM golang:1.16.2 AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
ENV MONGO_URI=${MONGO_DATABASE}
ENV DATABASE=${MONGO_COLLECTION}
WORKDIR /root/
COPY --from=builder /app .
CMD ["./app"]