FROM golang:bookworm AS builder

WORKDIR /work
COPY . /work/

RUN go build .

FROM ghcr.io/yude/use-groff-docker:master AS runner

WORKDIR /work
COPY --from=builder /work/groffer /usr/local/bin/groffer
COPY do.sh /usr/local/bin/do-render
RUN chmod +x /usr/local/bin/do-render

ENTRYPOINT ["groffer"]