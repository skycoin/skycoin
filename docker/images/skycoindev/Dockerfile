# skycoin build
# reference https://github.com/skycoin/skycoin
FROM golang:1.10.0-stretch

# dirs
RUN apt-get update && apt-get install -y cmake gdbserver gdb vim ctags vim-scripts screen
RUN git clone --recurse-submodules -j8 https://github.com/Snaipe/Criterion
RUN cd /go/Criterion && cmake .  && make install
RUN go get -u github.com/derekparker/delve/cmd/dlv

WORKDIR $GOPATH/src/github.com/skycoin/skycoin
VOLUME $GOPATH/src/

ENV LD_LIBRARY_PATH=/usr/local/lib
