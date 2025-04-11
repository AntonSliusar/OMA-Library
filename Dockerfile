FROM golang:latest AS builder

WORKDIR /app

COPY . .

RUN go build -o main .

CMD [ "./main" ]


