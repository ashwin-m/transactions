ARG GO_VERSION=1.22

FROM golang:${GO_VERSION}

WORKDIR /go/src/app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o main main.go

EXPOSE 8080

CMD ["./main"]