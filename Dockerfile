FROM golang:1.17-rc-alpine3.13
 

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

ENTRYPOINT ["/app/txbot"]
