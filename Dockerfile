#
# Build:
#

FROM golang:1.15-buster AS build

ENV TINI_VERSION v0.19.0
ADD https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini /tini
RUN chmod +x /tini

WORKDIR /app

COPY ./src/go.mod  .
COPY ./src/go.sum  .

RUN go mod download

COPY ./src  ./

RUN CGO_ENABLED=0    \
    GOOS=linux       \
    GOARCH=amd64     \
    go build -a -o /server

#
# Dist:
#

FROM gcr.io/distroless/static-debian10 AS dist

COPY --from=build /tini     /tini
COPY --from=build /server   /server
COPY              ./assets  /assets

ENV HOST=0.0.0.0
ENV PORT=8080
ENV MODE=production
ENV ASSETS=/assets

ENTRYPOINT ["/tini", "--"]
CMD ["/server"]
