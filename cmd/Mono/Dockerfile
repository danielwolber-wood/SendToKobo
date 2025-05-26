FROM --platform=linux/amd64 golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .
RUN chmod +x ./main

FROM --platform=linux/amd64 alpine:latest
RUN apk --no-cache add ca-certificates pandoc
WORKDIR /root/
COPY --from=builder /app/main ./main
EXPOSE 8080
CMD ["./main"]