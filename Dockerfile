###############################
FROM golang:1.18-alpine AS build

RUN mkdir -p /go/src/github.com/subchen/frep
COPY . /go/src/github.com/subchen/frep
WORKDIR /go/src/github.com/subchen/frep

RUN apk add --no-cache make git
RUN make clean build-linux

###############################
FROM alpine:3.15

COPY --from=build /go/src/github.com/subchen/frep/_releases/frep-* /usr/local/bin/frep

ENTRYPOINT [ "/usr/local/bin/frep" ]
CMD [ "--help" ]
