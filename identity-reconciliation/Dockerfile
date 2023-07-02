FROM golang:1.18.4-alpine3.16
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go build -o main .
EXPOSE 8080
ENTRYPOINT ["/app/main", "api", "development"]





