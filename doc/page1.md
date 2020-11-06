# 代理

讲一个故事,可能不太准确。

## 反向代理
> 不知道处理方是谁

小明想去旅行，订一个旅行团。旅行团安排了火车-飞机-汽车-轮船最终到达了目的地

> 添加这行代码 server.reverseHandler(req)
> curl http://localhost:8081/1


## 正向代理
> 知道我要去哪儿?

小明想去旅行，订一张机票直接飞到目的地

> curl -x http://localhost:8081 http://www.baidu.com 


## 参考

[https://github.com/fagongzi/manba](https://github.com/lafikl/liblb)

[https://github.com/fagongzi/manba](https://github.com/fagongzi/manba)

[https://github.com/panjf2000/goproxy](https://github.com/panjf2000/goproxy)



## 已完成

- 代理请求

```
curl http://localhost:8081/1


实际访问 http://localhost:8080/1
```


- 直接改写请求

```
curl http://localhost:8081/1
对请求的url进行修改就可以访问制定的后台服务器

curl http://localhost:9070/push/1
```









