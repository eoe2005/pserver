FROM alpine
MAINTAINER eoe2005
WORKDIR /worker
COPY main.go /worker
# RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && apk update && apk add go && go build -o /worker/main /worker/main.go
RUN apk update && apk add go && go build -o /worker/main /worker/main.go


EXPOSE 8888
ENTRYPOINT ["/worker/main"]
