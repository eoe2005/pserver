FROM alpine
MAINTAINER eoe2005
WORKDIR /worker
COPY . /worker
# RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && apk update && apk add go && go build -o /worker/main /worker/main.go
RUN apk update && apk add go && go build -o /worker/main /worker/server/main.go


EXPOSE 8888
ENTRYPOINT ["/worker/main"]
