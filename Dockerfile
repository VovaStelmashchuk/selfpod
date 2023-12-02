FROM golang:1.20-alpine

RUN apk update & apk upgrade
RUN apk --no-cache add ca-certificates wget bash xz-libs git
WORKDIR /tmp
RUN wget https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-amd64-static.tar.xz
RUN tar -xvJf  ffmpeg-release-amd64-static.tar.xz
RUN cd ff* && mv ff* /usr/local/bin

WORKDIR /

ENV GO111MODULE=on

WORKDIR /application

COPY . /application

RUN go build -o application .

EXPOSE 5000

CMD ["./application"]
