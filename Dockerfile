FROM golang:1.25.5-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify
COPY . .

RUN go build -o api-monitoring ./cmd/server/

EXPOSE 8080

CMD ["./api-monitoring"]