Mac環境建置
---------------
Mac環境，可以直接使用Homebrew安裝
```
brew install go
```
設定環境變數
```
export PATH=$PATH:/usr/local/opt/go/libexec/bin
export GOROOT=$HOME/go1.X
```

檢查一下是否安裝成功
```
go version
```
Ubuntu環境建置
---------------
更新套件包
```
sudo apt-get update  
sudo apt-get -y upgrade  
```

download go binary
```
wget https://dl.google.com/go/go1.17.7.linux-amd64.tar.gz 
```

設定環境變數
```
sudo tar -xvf go1.17.7.linux-amd64.tar.gz
sudo mv go /usr/local 
export GOROOT=/usr/local/go 
export GOPATH=$HOME
export PATH=$PATH:/usr/local/opt/go/libexec/bin
```

檢查一下是否安裝成功
```
go version
```

執行程式
---------------
使用打包好的執行檔 bin/demo

step 1 : (給執行檔任何人都可執行權限)
```
sudo chmod +x demo
```
step 2:(執行)
```
./demo 
```

使用souce code 執行
```
cd {project_path}
go run main.go 

```


打包成執行檔
---------------
打包在liunx x64 可執行檔（ubuntu 16.04 下可執行）
```
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/demo  forSkywatch
```
GOOS：目標平臺的作業系統（darwin、freebsd、linux、windows）
GOARCH：目標平臺的體系架構（386、amd64、arm）
