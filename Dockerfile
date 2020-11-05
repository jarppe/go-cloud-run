#
# Build:
#

FROM golang:1.15-buster AS build

WORKDIR /app

COPY ./src/go.mod .
# COPY ./src/go.sum .

RUN go mod download

COPY ./src/*.go ./

RUN go install

#
# Dist:
#

FROM debian:10-slim AS dist

WORKDIR /app

COPY --from=build /go/bin/go-cloud-run /app/go-cloud-run
COPY ./src/assets ./assets

CMD ["./go-cloud-run"]
