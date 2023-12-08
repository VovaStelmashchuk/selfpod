FROM arm64v8/golang:1.20-alpine

RUN apk update
RUN apk upgrade
RUN apk add --no-cache ffmpeg

WORKDIR /application

COPY . /application

RUN go build -o application .

EXPOSE 5000

CMD ["./application"]
