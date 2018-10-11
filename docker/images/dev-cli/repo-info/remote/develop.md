## `skycoin/skycoindev-cli:develop`

```console
$ docker pull skycoin/skycoin@sha256:7c1d1d2acfc098ae172b565e8b1975d37a3b8d945b82c78010c47c950cc020b7
```

- Manifest MIME: `application/vnd.docker.distribution.manifest.list.v2+json`
- Platforms:
	- linux; amd64

- Layers:
    - sha256:cc1a78bfd46becbfc3abb8a74d9a70a0e0dc7a5809bbd12e814f9382db003707
    - sha256:d2c05365ee2a2245bb9f6786bc88aa12bf64da676a52668424437826d0f0cb92
    - sha256:231cb0e216d30ea48044d44d37fba016eb67eca9b19b29a741d95775359d3533
    - sha256:3d2aa70286b89febc36370098220c9b2960cc67c03375c9df4e82736519f1e8a
    - sha256:6b387538c1ae7468ba635e703c1f55e68ef1634290958da056185f5ce529dba4
    - sha256:39f9c30eddc9211649a148492c59b81dfe6ecff769bd96462220d1329fb62718
    - sha256:f4e296016ab64df5cf193e6e2c3391e2f3e8689de8765dc966f14bcc042ce7a3
    - sha256:dd412a416bc36702f39e14277ae2e43103c7251f49c27b369c2de5dc4bc8d209
    - sha256:2d071012d44f8781afad94e13709c0d3ea6ac66dc9bdef0a9d35c893d045d13e
    - sha256:88c53c392054ea0bc1405502a0af2d7439b1d340690312ef6d5fae630ec025e5
    - sha256:748b5ee860b04de9f71702e12108249c4ed29ad7226a2d659b3cf14417c4a618
    - sha256:b6942952f9935f0f0bd813a2d379c1e6b22da183315cda30fd0b590806b7610a
    - sha256:22084e1c017fce593f31454bc27238c267953d6fc162eccef0928b204cf67593
    - sha256:9ae2147d285fb3720eab6350aa36c0611bd713ef51345e11b61617099f6abce4
    - sha256:8f210d9cfd5af953f8daff98f74a29b21ec115fa50155cb9f76a2f549f40ff81
    - sha256:a00e162b8cbd633a07d10b4cd105234db559f14bf4b87124b8bfa6faf0ea4d4e
    - sha256:92b5df9979d594a665691e680c469deca3726c3f6203d1b084f56ae5d9b31b96

- Expose Ports:
	- 6000
    - 6420

```dockerfile
# 2018-07-02T22:18:08.316631342Z
VOLUME [/go/src/]
# 2018-07-02T22:18:08.316631342Z
WORKDIR /go/src/github.com/skycoin\
# 2018-07-02T22:18:08.316631342Z
cd /tmp/; \
wget http://prdownloads.sourceforge.net/swig/swig-3.0.12.tar.gz && \
tar -zxf swig-3.0.12.tar.gz ;\
cd swig-3.0.12 ;\
./configure --prefix=/usr && \
make && \
make install && \
rm -rf /tmp/swig-*
# 2018-07-02T22:18:08.316631342Z
git clone https://github.com/fatih/vim-go /usr/share/vim/vim80/pack/dev/start/vim-go && \
git clone https://github.com/tpope/vim-fugitive /usr/share/vim/vim80/pack/dev/start/vim-fugitive && \
git clone https://github.com/Shougo/vimshell.vim /usr/share/vim/vim80/pack/dev/start/0vimshell && \
git clone https://github.com/Shougo/vimproc.vim /usr/share/vim/vim80/pack/dev/start/0vimproc && \
git clone https://github.com/w0rp/ale.git /usr/share/vim/vim80/pack/dev/start/ale && \
cd /usr/share/vim/vim80/pack/dev/start/0vimproc && \
make ; \
git clone https://github.com/iberianpig/tig-explorer.vim.git /tmp/tig-explorer; \
cp /tmp/tig-explorer/autoload/tig_explorer.vim /usr/share/vim/vim80/autoload; \
cp /tmp/tig-explorer/plugin/tig_explorer.vim /usr/share/vim/vim80/plugin;  \
rm -rf /tmp/tig-explorer
# 2018-07-02T22:18:08.316631342Z
go get -u github.com/derekparker/delve/cmd/dlv && \
go get -u github.com/FiloSottile/vendorcheck && \
go get -u github.com/alecthomas/gometalinter && \
gometalinter --vendored-linters --install && \
go get -u github.com/zmb3/gogetdoc && \
go get -u golang.org/x/tools/cmd/guru && \
go get -u github.com/davidrjenni/reftools/cmd/fillstruct && \
go get -u github.com/rogpeppe/godef && \
go get -u github.com/fatih/motion && \
go get -u github.com/nsf/gocode && \
go get -u github.com/jstemmer/gotags && \
go get -u github.com/josharian/impl && \
go get -u github.com/fatih/gomodifytags && \
go get -u github.com/dominikh/go-tools/cmd/keyify && \
go get -u golang.org/x/tools/cmd/gorename && \
go get -u github.com/klauspost/asmfmt/cmd/asmfmt && \
go get -u github.com/vektra/mockery/.../ && \
curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
```
