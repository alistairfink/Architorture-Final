FROM golang:1.16.3 as backend-builder
LABEL maintainer="AlistairFink <alistairfink@gmail.com>"

WORKDIR /go/src/github.com/alistairfink/Architorture-Backend
COPY ./Backend .
RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/Architorture-Backend .

FROM node:10.24.1-alpine3.11 as frontend-builder
LABEL maintainer="AlistairFink <alistairfink@gmail.com>"

COPY ./Frontend .
RUN npm install
RUN npm run build

FROM docker.io/library/nginx:latest
COPY ./Frontend/default.conf /etc/nginx/conf.d/default.conf
COPY --from=backend-builder /go/bin/Architorture-Backend .
COPY --from=frontend-builder /build /usr/share/nginx/html

CMD nginx && /Architorture-Backend