FROM golang:1.17-rc-alpine3.13 AS builder

WORKDIR /app
COPY . . 

RUN apk -v --update --no-cache add \
	curl \
	git \
	groff \
	less \
	mailcap \
	gcc \
	libc-dev \
	bash  \
	leveldb-dev  

RUN go build

FROM alpine:3.13

WORKDIR /app 

COPY --from=builder /app/txbot /app/txbot

ENTRYPOINT ["/app/txbot"]
