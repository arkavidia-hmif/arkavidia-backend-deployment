FROM golang:1.13

EXPOSE 8080

WORKDIR /code
COPY . /code

RUN go build -o deploy

CMD ["/code/start.sh"]
