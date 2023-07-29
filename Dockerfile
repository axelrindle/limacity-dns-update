FROM golang:1.20-alpine as build

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/app


FROM golang:1.20-alpine

ENV CONTAINER=true

COPY --from=build /usr/local/bin/app /

HEALTHCHECK CMD sh -c "[ ! -f /tmp/failure ]"

CMD ["/app"]
