FROM golang:1.14 AS build

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
ENV GOFLAGS="-mod=vendor"

ARG version="undefined"

WORKDIR /build/balerter

ADD . /build/balerter

RUN go build -o /balerter -ldflags "-X main.version=${version} -s -w"  ./cmd/balerter

# -----

FROM  debian:stretch-slim
COPY --from=build /balerter /
COPY --from=build /build/balerter/modules /modules

RUN apt-get update \
     && apt-get install -y --no-install-recommends ca-certificates

RUN update-ca-certificates

ENTRYPOINT ["/balerter"]

CMD ["/balerter"]
