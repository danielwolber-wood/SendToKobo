FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go mod download
RUN go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates pandoc
WORKDIR /root/
COPY --from=builder /app/main ./main
EXPOSE 8080
CMD ["./main"]