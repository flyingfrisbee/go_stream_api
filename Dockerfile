FROM golang:alpine as base
# RUN apk add build-base

WORKDIR /builder

COPY . /builder
RUN go build -o /start-go .

FROM alpine:latest
WORKDIR /app

COPY ./.env /app
COPY ./repository/database/relational/init.sql /app
COPY ./repository/database/prod-ca-2021.crt /app

COPY --from=base /start-go /app/main

CMD ["/app/main"]