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
RUN chmod +x /usr/src/image/bin/runc

RUN curl -fL -o /usr/src/containerd.tgz https://github.com/containerd/containerd/releases/download/v1.1.0/containerd-1.1.0.linux-amd64.tar.gz
RUN tar xvf /usr/src/containerd.tgz -C /usr/src/image/bin bin/containerd bin/containerd-shim --strip-components=1

RUN curl -fL -o /usr/src/cni.tgz https://github.com/containernetworking/plugins/releases/download/v0.7.1/cni-plugins-amd64-v0.7.1.tgz
RUN mkdir /usr/src/image/cni
RUN tar xvf /usr/src/cni.tgz -C /usr/src/image/cni ./bridge ./flannel ./host-local ./loopback ./portmap

RUN mkdir -p /usr/src/image/etc/containerd
COPY config.toml /usr/src/image/etc/containerd/config.toml
COPY init /usr/src/image/init

RUN mkdir -p /usr/src/image/etc/cni/net.d
COPY 10-flannel.conflist /usr/src/image/etc/cni/net.d

RUN rm -rf /usr/src/image/etc/init.d
COPY inittab /usr/src/image/etc/inittab

RUN mkdir -p /usr/src/image/etc/kube-flannel
COPY net-conf.json /usr/src/image/etc/kube-flannel/

COPY agent /usr/src/image/bin/agent

RUN cd /usr/src/image && \
    mkdir -p lib/modules lib/firmware && \
    rm -rf dev etc/network share/udhcpc

FROM scratch
COPY --from=ubuntu /usr/src/image/ /
CMD ["/bin/bash"]
