FROM golang:1.23.3

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY ./src ./src

RUN go build src/server.go

VOLUME ["/app/portfolio.yaml"]

EXPOSE 8080

CMD [ "./server" ]
