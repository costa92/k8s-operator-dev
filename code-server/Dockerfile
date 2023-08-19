FROM ubuntu:20.04
WORKDIR /.workspace

COPY code-server-4.11.0-linux-arm64.tar.gz .
COPY go1.19.4.linux-amd64.tar.gz .


RUN tar zxvf code-server-4.11.0-linux-arm64.tar.gz && \
    rm -f code-server-4.11.0-linux-arm64.tar.gz && \
    tar zxvf go1.19.4.linux-amd64.tar.gz -C /usr/local  && \
    rm -f  go1.19.4.linux-amd64.tar.gz && \
    mkdir -p /root/workspace /go/{src,pkg,bin} && \
    apt-get -y update && \
    apt-get -qq update  && \
    apt-get install libterm-readkey-perl -y && \
    apt-get -qq install -y --no-install-recommends ca-certificates curl  && \
    apt-get install git -y

ENV GO111MODULE on
ENV GOPROXY https://goproxy.cn,direct
ENV GOROOT /usr/local/go
ENV PATH /usr/local/go/bin:$PATH
ENV GOPATH /go
ENV PATH &GOPATH/bin:$PATH

WORKDIR /.workspace/code-server-4.11.0-linux-arm64

EXPOSE 9999

CMD ["./bin/code-server","--port","9999","--host","0.0.0.0","--auth","none","--disable-update-check","--open","/root/workspace"]