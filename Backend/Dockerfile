FROM golang:1.16.3 as builder
LABEL maintainer="AlistairFink <alistairfink@gmail.com>"

WORKDIR /go/src/github.com/alistairfink/Architorture-Backend
COPY . .
RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/Architorture-Backend .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/bin/Architorture-Backend .

EXPOSE 8080

CMD ["./Architorture-Backend"] 