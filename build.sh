#! /bin/bash

WORKSPACE=`pwd`
cd $WORKSPACE
mkdir bin

# build cli
cd $WORKSPACE/bin
mkdir cli
cd cli
touch main.go
echo "package main
import \"github.com/mebiusashan/beaker\"
func main(){
    beaker.RunCli()
}" > main.go

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o beaker_linux main.go
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o beaker_mac main.go
rm main.go

# build linux server & admin
cd $WORKSPACE/bin
mkdir linux
cd linux
touch main.go
echo "package main
import \"github.com/mebiusashan/beaker\"
func main(){
    beaker.RunServer(true)
}" > main.go

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o beaker_server main.go

echo "package main
import \"github.com/mebiusashan/beaker\"
func main(){
    beaker.RunAdmin(true)
}" > main.go

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o beaker_admin main.go
rm main.go
cp $WORKSPACE/install/* ./
cp -r $WORKSPACE/static ./
cp -r $WORKSPACE/temp ./


# build darwin server & admin
cd $WORKSPACE/bin
mkdir darwin
cd darwin
touch main.go
echo "package main
import \"github.com/mebiusashan/beaker\"
func main(){
    beaker.RunServer(true)
}" > main.go

CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o beaker_server main.go

echo "package main
import \"github.com/mebiusashan/beaker\"
func main(){
    beaker.RunAdmin(true)
}" > main.go

CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o beaker_admin main.go
rm main.go
cp $WORKSPACE/install/* ./
cp -r $WORKSPACE/static ./
cp -r $WORKSPACE/temp ./

cd $WORKSPACE/bin
tar -czvf beaker-v0.1.1-cli.tar.gz ./cli
tar -czvf beaker-v0.1.1--linux-amd64.tar.gz ./linux
tar -czvf beaker-v0.1.1-darwin-amd64.tar.gz ./darwin
