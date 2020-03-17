FROM golang:1.14 AS build

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
ENV GOFLAGS="-mod=vendor"

WORKDIR /build/tgtool

ADD . /build/tgtool
RUN go build -o /tgtool ./cmd/tgtool

# -----

FROM alpine:3.11.3
COPY --from=build /tgtool /

ENTRYPOINT ["/tgtool"]

CMD ["/tgtool"]
