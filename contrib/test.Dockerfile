FROM golang:1.18 AS build

ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64
ENV GOFLAGS="-mod=vendor"

ARG version="undefined"

WORKDIR /build/test

ADD . /build/test

RUN go build -o /test -ldflags "-X main.version=${version} -s -w"  ./cmd/test

# -----

FROM  debian:stretch-slim
COPY --from=build /test /
COPY --from=build /build/test/modules /modules

ENTRYPOINT ["/test"]

CMD ["/test"]
