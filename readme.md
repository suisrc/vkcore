# 说明

这是一个解决各种挑战的解决方案

## 特别声明

所有的内容仅限于测试，不可用于生产

## golang

``` bash
go mod init vkc
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


## voska

model默认位置 data/vosk

```bash
apt-get install -y unzip
cd data
&& curl -Lo model.zip https://alphacephei.com/vosk/models/vosk-model-small-en-us-0.15.zip \
&& unzip model.zip && mv vosk-model-small-en-us-0.15 vosk \
&& curl -Lo vosk.zip  https://github.com/alphacep/vosk-api/releases/download/v0.3.45/vosk-linux-x86_64-0.3.45.zip \
&& unzip vosk.zip && mv vosk-linux-x86_64-0.3.45 vosk-so \
&& cp vosk-so/vosk_api.h /usr/include/ && cp vosk-so/libvosk.so /usr/lib/ \
```