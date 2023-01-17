FROM golang:alpine

WORKDIR /App

COPY go.mod .

RUN go mod tidy

COPY . .

RUN go build main.go

EXPOSE 60000

CMD [ "./main" ]