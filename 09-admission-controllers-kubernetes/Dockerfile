FROM golang:1.20.3 AS base

ENV GOOS=linux \
    GOARCH=arm64 \
    CGO_ENABLED=0

WORKDIR /work
COPY . /work

RUN go build -o /admission .

FROM scratch

COPY --from=base /admission /usr/local/bin/

CMD ["admission"]