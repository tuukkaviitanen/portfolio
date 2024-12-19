FROM golang:1.23.3

WORKDIR /usr/src/app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build src/server.go

CMD [ "./server" ]
