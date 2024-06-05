FROM golang:alpine3.20 as build-stage

WORKDIR /app
COPY go.mod ./
RUN go mod download

COPY cmd ./cmd
RUN go build -o mock-endpoint ./cmd/web

FROM alpine:3.20 as serve-stage

COPY --from=build-stage /app/mock-endpoint /mock-endpoint

EXPOSE 8080

CMD [ "/mock-endpoint" ]
