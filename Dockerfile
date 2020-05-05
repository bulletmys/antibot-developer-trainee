FROM golang:1.14 AS stage
ENV GO111MODULE=on
WORKDIR /go/src/rate_limit
COPY . /go/src/rate_limit
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

FROM alpine
WORKDIR /app
COPY --from=stage  /go/src/rate_limit /app
RUN chmod +x  ./main
EXPOSE 8080/tcp
ENTRYPOINT [ "/app/main" ]