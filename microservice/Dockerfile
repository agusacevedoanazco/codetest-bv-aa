FROM golang:alpine3.20 as build-stage

WORKDIR /app
COPY go.mod ./
RUN go mod download

COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY ui/ ./ui/

RUN go build -o microservice ./cmd/web

FROM alpine:3.20 as svc-stage

COPY --from=build-stage --chown=10000:10000 /app/microservice /microservice
USER 10000:10000
EXPOSE 32080

CMD ["/microservice"]
