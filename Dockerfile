FROM golang:1.23.4-bullseye

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o /godocker

EXPOSE 8080

CMD [ "/godocker" ]