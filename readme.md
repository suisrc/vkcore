# 说明


同时这也是最小化server框架

## golang

``` bash
go mod init vkcc
go mod tidy
go run main.go
go build


```

## 自签名证书

``` bash
cd cert

# 生成私钥
openssl genrsa -out default.key.pem 2048

# 生成自签名证书
openssl req -new -x509 -key default.key -out default.crt.pem -days 36500
```
 ## palywright

 ``` bash
# use
go get github.com/playwright-community/playwright-go

# deps /etc/hosts playwright.azureedge.net  13.107.246.70
go run github.com/playwright-community/playwright-go/cmd/playwright install firefox --with-deps
go run github.com/playwright-community/playwright-go/cmd/playwright install --with-deps
# Or
go install github.com/playwright-community/playwright-go/cmd/playwright
playwright install --with-deps

chown root:root ~
```