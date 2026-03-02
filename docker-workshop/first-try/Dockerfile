FROM dhi.io/golang:1.26

WORKDIR /home

COPY ./main.go .

RUN go mod init main && go mod tidy

RUN go build -o server .

EXPOSE 8080

CMD ["/home/server"]