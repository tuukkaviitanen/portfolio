FROM golang:1.23.0

WORKDIR /usr/src/app

COPY . .

RUN go build src/server.go

CMD [ "./server" ]
