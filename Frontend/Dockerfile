FROM node:10.24.1-alpine3.11 as builder
LABEL maintainer="AlistairFink <alistairfink@gmail.com>"

COPY . .
RUN npm install
RUN npm run build

RUN ls -al /build

#FROM alpine:latest
#RUN apk --no-cache add ca-certificates
#WORKDIR /root/
#COPY --from=builder /go/bin/Architorture-Backend .
#
#EXPOSE 8080
#
#CMD ["./Architorture-Backend"]

FROM docker.io/library/nginx:latest
COPY ./default.conf /etc/nginx/conf.d/default.conf
COPY --from=builder /build /usr/share/nginx/html