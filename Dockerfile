FROM golang:1.19.4-alpine3.17 AS build
RUN apk add build-base
#ENV CGO_ENABLED=0
COPY . /app
WORKDIR /app
RUN go build -o app.bin ./cmd
RUN chmod +x ./app.bin

FROM alpine:3.17
COPY --from=build /app/app.bin /app.bin
COPY ./.env /.env

EXPOSE 7998

CMD ["/app.bin"]