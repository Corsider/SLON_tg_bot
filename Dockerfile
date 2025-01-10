FROM golang:latest
LABEL authors="Corsider"

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build main.go

CMD ["./main"]