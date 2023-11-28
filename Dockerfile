FROM golang:1.20

RUN apt-get update && apt-get install -y ffmpeg

WORKDIR /application

COPY . /application

RUN go build -o application .

EXPOSE 5000

CMD ["./application"]
