FROM golang:1.23-alpine as build

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -ldflags="-w -s" -o /usr/local/bin/app


FROM alpine:3

ENV CONTAINER=true

COPY --from=build /usr/local/bin/app /

HEALTHCHECK CMD sh -c "[ ! -f /tmp/failure ]"

USER 1001

CMD ["/app"]
