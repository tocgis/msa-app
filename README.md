## HCI MICROSERVICE-APP

### 1. Introduction

This is Microservcie Architecture use GO-KIT & Gin

### 2. 开发流程
1. 开发服务端
   1. 书写 .proto 文件，定义请求和返回逻辑
   2. 编译 .proto 使用 plugins=grpc
   3. 书写 srv 代码 service.go & middleware.go
   4. 书写 入口文件 main.go
2. 开发客户端
   1. 书写 Client 代码，封装业务请求逻辑
   2. 书写 handler 封装 api 入口
3. 把客户端 handler 注册进 api 总入口


### 3. Feature
功能 | 介绍
--------|-----------------
api         |  注册app所有endpoint.
client      |  所有访问微服务的客户端, 供apiGateway调用. 提供服务发现,负载均衡,错误重试和故障降级等功能.
cmd         |  各个服务的启动命令.
docker      |  构建各个服务的docker镜像.
monitor     |  监控组件.
proto       |  服务间IPC方式采用grpc.
tracer      |  分布式跟踪.
