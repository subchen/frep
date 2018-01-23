###############################
FROM golang:1.9-alpine AS build

RUN mkdir -p /go/src/github.com/subchen/frep
COPY . /go/src/github.com/subchen/frep
WORKDIR /go/src/github.com/subchen/frep

RUN apk add --no-cache make
RUN make build-linux

###############################
FROM alpine:3.7

COPY --from=build /go/src/github.com/subchen/frep/releases/frep-* /usr/local/bin/frep

ENTRYPOINT [ "/usr/local/bin/frep" ]
CMD [ "--help" ]
