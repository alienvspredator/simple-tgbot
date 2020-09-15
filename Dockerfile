FROM golang:1.15 as builder

ARG SERVICE

RUN apt-get -qq update && apt-get -yqq install upx

ENV GO111MODULE=on \
    CGO_ENABLED=0

WORKDIR /src
COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY cmd cmd
COPY internal internal
COPY pkg pkg

RUN go build \
  -trimpath \
  -installsuffix cgo \
  -tags netgo \
  -o /bin/service \
  ./cmd/${SERVICE}

RUN strip /bin/service
RUN upx -q -9 /bin/service

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /bin/service /bin/service

ENTRYPOINT ["/bin/service"]
