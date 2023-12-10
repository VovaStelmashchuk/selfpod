FROM arm64v8/golang:1.20-alpine

RUN apk update
RUN apk upgrade
RUN apk add --no-cache ffmpeg

WORKDIR /application

COPY . /application

ENV CGO_ENABLED=1
RUN go build -o application .

EXPOSE 5000

CMD ["./application"]
