FROM golang:alpine

WORKDIR /app
COPY . .

RUN go mod download

ENV AUTH_SERVER_REDIS_PORT="6379"
ENV REDIS_USER1_PASS=""
ENV REDIS_USER1_NAME=""
ENV LOGIN_SERV_PORT=8080

RUN go build

EXPOSE 8080

CMD [ "./robinservicologin" ]