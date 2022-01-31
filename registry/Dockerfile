FROM golang:1.17.5-alpine as builder
WORKDIR /build

COPY ./ /build/
RUN CGO_ENABLED=0 GOOS=linux go build -a -o registry ./cmd

FROM alpine:3.14.3
COPY --from=builder /build/registry ./

EXPOSE 80

ENTRYPOINT [ "./registry" ]
