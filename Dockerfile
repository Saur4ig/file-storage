FROM golang:1.22.5 AS builder

WORKDIR /app
ENV CGO_ENABLED=0
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /file_storage ./app/main.go


FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /file_storage ./file_storage

ENTRYPOINT ["/root/file_storage"]
