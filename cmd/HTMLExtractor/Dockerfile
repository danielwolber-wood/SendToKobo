FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go mod download
RUN go build -o main .
RUN chmod +x ./main

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main ./main
EXPOSE 80
CMD ["./main"]
