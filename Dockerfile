FROM golang:1.12 as builder
ADD . /decoherence
WORKDIR /decoherence
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' .

FROM scratch
COPY --from=builder /decoherence/decoherence /decoherence
WORKDIR /workdir
ENTRYPOINT ["/decoherence"]