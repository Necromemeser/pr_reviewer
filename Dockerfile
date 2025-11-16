FROM golang:1.25.4

WORKDIR /app

COPY . .

RUN go mod tidy && go build -o app ./cmd/app

EXPOSE 8080

CMD ["./app"]