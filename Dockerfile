FROM golang

RUN mkdir /app

ADD . /app

WORKDIR /app

RUN go build -o ./bin ./client
RUN go build -o ./bin ./server

EXPOSE 8080
CMD [ "/app/bin/server" ]