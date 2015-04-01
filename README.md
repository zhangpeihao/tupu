# tupu
图普科技开放平台API接口SDK

## 安装
go get github.com/zhangpeihao/tupu

## 使用
从[图普开放平台](http://open.tuputech.com)取得接口调用URL、secretId和modelId。
准备好PKCS8格式的私钥。
参考demo加载图片和私钥。
```
req := tupu.NewRequest(TUPU_API_URL, SECRET_ID, MODEL_ID, private_key)
resp, err := req.CheckSingleImage(imgBuf, imgName)
```
