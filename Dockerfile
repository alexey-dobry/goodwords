FROM golang:alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED=0

RUN apk update --no-cache && apk add --no-cache tzdata

WORKDIR /app/service/build

COPY ./service/go.mod .
COPY ./service/go.sum .
RUN go mod download

COPY ./service/cmd .
COPY ./service/internal .
RUN go build -ldflags="-s -w" -o /app/service ./cmd

FROM scratch

WORKDIR /app/service

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /app/service /app/service

CMD ["./cmd"]