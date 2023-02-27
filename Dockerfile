FROM golang:1.19

COPY ./ ./

RUN go mod vendor

RUN go build app.go

EXPOSE 8080
CMD ["app"]