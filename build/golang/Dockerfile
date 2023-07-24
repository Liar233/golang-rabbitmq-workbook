FROM golang:1.20-buster AS build-dev-server

WORKDIR /app
COPY . .

RUN go mod download
RUN go mod verify

EXPOSE 8080

CMD ["sh", "/app/build/golang/daemon.sh"]