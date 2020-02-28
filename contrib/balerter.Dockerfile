FROM golang:1.13 AS build

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
ENV GOFLAGS="-mod=vendor"

ARG version="undefined"

WORKDIR /build/balerter

ADD . /build/balerter
RUN go test ./...
RUN go build -o /balerter -ldflags "-X main.version=${version} -s -w"  ./cmd/balerter

# -----

FROM alpine:3.11.3
COPY --from=build /balerter /

ENTRYPOINT ["/balerter"]

CMD ["/balerter"]
