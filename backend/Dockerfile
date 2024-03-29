FROM golang:1.15-buster as base
WORKDIR /go/src/app
ENV GOOS=linux
ENV GOARCH=amd64
COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify
COPY . .

FROM base as build
RUN go build main.go

FROM alpine:3.14 as prod
WORKDIR /home/app
RUN apk update && \
    apk add --no-cache nginx supervisor libc6-compat && \
    rm /etc/nginx/http.d/default.conf && \
    mkdir logs
COPY --from=build /go/src/app/main ./
COPY ./deploy/nginx/app.conf /etc/nginx/http.d/
COPY ./deploy/supervisord.conf /etc/supervisord.conf
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisord.conf"]