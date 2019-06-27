VERSION=1.12.5
OS=linux
ARCH=amd64
wget https://dl.google.com/go/go$VERSION.$OS-$ARCH.tar.gz
sudo tar -C /usr/local -xzf go$VERSION.$OS-$ARCH.tar.gz
export PATH=$PATH:/usr/local/go/bin
export GOPATH=$HOME/go
