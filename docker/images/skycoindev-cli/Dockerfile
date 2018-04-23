# Creates an image for skycoin development
FROM golang:1.10.0-stretch

# Installs software
RUN apt-get update && apt-get install -y cmake gdbserver gdb vim ctags vim-scripts screen

# Installs Criterion
RUN git clone --recurse-submodules -j8 https://github.com/Snaipe/Criterion && \
    cd /go/Criterion && \
    cmake . && \
    make install && \
    rm -r /go/Criterion

# Installs 
RUN go get -u github.com/derekparker/delve/cmd/dlv  && \
    git clone https://github.com/fatih/vim-go /usr/share/vim/vim80/pack/dev/start/vim-go && \
    git clone https://github.com/tpope/vim-fugitive /usr/share/vim/vim80/pack/dev/start/vim-fugitive && \
    git clone https://github.com/Shougo/vimshell.vim /usr/share/vim/vim80/pack/dev/start/0vimshell && \
    git clone https://github.com/Shougo/vimproc.vim /usr/share/vim/vim80/pack/dev/start/0vimproc && \

WORKDIR $GOPATH/src/github.com/skycoin/skycoin
VOLUME $GOPATH/src/

ENV LD_LIBRARY_PATH=/usr/local/lib
