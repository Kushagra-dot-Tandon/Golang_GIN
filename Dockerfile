FROM golang:1.16-alpine
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go mod download
RUN go build -o /docker-gin
EXPOSE 8080
CMD [ "/docker-gin" ]