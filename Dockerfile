FROM alpine
MAINTAINER eoe2005
WORKDIR /worker
COPY . /worker
# RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && apk update && apk add go && go env -w GO111MODULE=on && go env -w GOPROXY=https://goproxy.cn,direct && go build -o /usr/bin/goproxyagent /worker/server/main.go && apk del go && rm -Rf /root/go && rm -Rf /worker
RUN apk update && apk add go && go build -o /usr/bin/goproxyagent /worker/server/main.go  && apk del go && rm -Rf /root/go && rm -Rf /worker


EXPOSE 8888
ENTRYPOINT ["/usr/bin/goproxyagent"]
