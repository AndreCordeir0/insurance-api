FROM golang:1.16-alpine

WORKDIR "/app/insurance-api"

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o insurance-api cmd/main.go


CMD [ "./insurance-api" ]
