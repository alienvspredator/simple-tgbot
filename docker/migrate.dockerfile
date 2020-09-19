FROM golang:1.15 AS installer

RUN go get github.com/cockroachdb/cockroach-go/crdb \
    && go get github.com/lib/pq \
    && go get github.com/golang-migrate/migrate \
    && go install -tags=cockroachdb github.com/golang-migrate/migrate/cmd/migrate

FROM alpine
COPY --from=installer /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=installer /go/bin/migrate /usr/local/bin/migrate

ENTRYPOINT [ "/usr/local/bin/migrate" ]
