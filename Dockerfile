FROM golang:alpine

WORKDIR /usr/src/app

RUN go install github.com/air-verse/air@latest

COPY . .

RUN go mod tidy

EXPOSE 8080

CMD ["air"]