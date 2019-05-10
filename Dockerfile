FROM golang:1.12 AS builder
WORKDIR /plentiform
COPY . .
RUN CGO_ENABLED=0 go build -a -installsuffix cgo

FROM alpine
CMD ["/plentiform"]
COPY --from=builder /plentiform/public /public
COPY --from=builder /plentiform/templates /templates
COPY --from=builder /plentiform/plentiform /plentiform
