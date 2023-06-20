FROM golang:1.20.4-alpine as builder
LABEL maintainer = "Chenghao Yu <yuchenghao0624@qq.com>"

WORKDIR /usr/src/edge-alert
COPY ./ /usr/src/edge-alert/
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct
RUN go mod tidy && go build

FROM alpine:3
RUN mkdir -p /etc/edge-alert
COPY --from=builder /usr/src/edge-alert/edge-alert /usr/bin/
COPY ./etc/edge-alert/conf /etc/edge-alert/conf
RUN chmod u+rw /etc/edge-alert/conf/conf.toml
EXPOSE 30000
CMD ["sh", "-c", "edge-alert"]