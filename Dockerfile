FROM golang:1.22-bookworm AS build-env

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . ./

ENV GOARCH=amd64

## Deploy
FROM gcr.io/distroless/base-debian11

EXPOSE 8081

USER nonroot:nonroot

CMD [ "/app" ]