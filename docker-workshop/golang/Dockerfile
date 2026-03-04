FROM golang:1.24.0-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod init main
RUN go mod tidy
RUN go build -o server .

#The final image
FROM golang:1.24.0-alpine

COPY --from=builder /app/server /app/server

CMD ["/app/server"]