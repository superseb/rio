FROM envoyproxy/envoy:latest AS envoy

FROM golang:1.10 AS cni
RUN mkdir -p $GOPATH/src/github.com/containernetworking && \
    cd $GOPATH/src/github.com/containernetworking && \
    git clone https://github.com/ibuildthecloud/plugins.git && \
    cd plugins && \
    git checkout 299b021fdcb5533d5babc40006d2965a8563066a
RUN go build -ldflags -s -o /usr/bin/cni github.com/containernetworking/plugins

FROM ubuntu:18.04 AS ubuntu
RUN yes | unminimize
RUN apt-get update
RUN apt-get install -y \
    build-essential \
    gcc \
    g++ \
    rsync \
    bc \
    wget \
    curl \
    ca-certificates \
    ncurses-dev \
    python \
    unzip
RUN mkdir /usr/src/buildroot
RUN curl -fL https://buildroot.org/downloads/buildroot-2018.02.2.tar.bz2 | tar xvjf - -C /usr/src/buildroot --strip-components=1
RUN curl -fL https://storage.googleapis.com/buildroot-cache/2018.02.2.tar.gz | tar xvzf - -C /usr/src/buildroot
WORKDIR /usr/src/buildroot
COPY config /usr/src/buildroot/.config
COPY busybox.config /usr/src/buildroot/package/busybox/busybox.config
RUN make oldconfig
RUN make source
RUN make

RUN mkdir /usr/src/image && \
    tar xvf /usr/src/buildroot/output/images/rootfs.tar -C /usr/src/image ./etc && \
    tar xvf /usr/src/buildroot/output/images/rootfs.tar -C /usr/src/image ./usr --strip-components=2

RUN curl -fL -o /usr/src/image/bin/runc https://github.com/opencontainers/runc/releases/download/v1.0.0-rc5/runc.amd64
RUN chmod +x /usr/src/image/bin/runc && strip --strip-debug /usr/src/image/bin/runc

RUN curl -fL -o /usr/src/containerd.tgz https://github.com/containerd/containerd/releases/download/v1.1.1/containerd-1.1.1.linux-amd64.tar.gz
RUN tar xvf /usr/src/containerd.tgz -C /usr/src/image/bin bin/containerd bin/containerd-shim --strip-components=1

RUN mkdir -p /usr/src/image/share/cni/bin
COPY --from=cni /usr/bin/cni /usr/src/image/share/cni/bin
RUN cd /usr/src/image/share/cni/bin && \
    for i in ./bridge ./flannel ./host-local ./loopback ./portmap; do \
        ln -s cni $i \
    ;done

RUN mkdir -p /usr/src/image/etc/containerd
COPY config.toml /usr/src/image/etc/containerd/config.toml
COPY init /usr/src/image/init

RUN mkdir -p /usr/src/image/etc/cni/net.d
COPY 10-flannel.conflist /usr/src/image/etc/cni/net.d

RUN rm -rf /usr/src/image/etc/init.d
COPY inittab /usr/src/image/etc/inittab

RUN mkdir -p /usr/src/image/etc/kube-flannel
COPY net-conf.json /usr/src/image/etc/kube-flannel/

COPY envoy.yaml /etc/envoy/envoy.yaml
COPY --from=envoy /usr/local/bin/envoy /usr/src/image/bin/envoy

COPY agent /usr/src/image/bin/agent

RUN cd /usr/src/image && \
    mkdir -p lib/modules lib/firmware && \
    rm -rf dev etc/network share/udhcpc

FROM scratch
COPY --from=ubuntu /usr/src/image/ /
CMD ["/bin/bash"]
