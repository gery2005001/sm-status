### spacemesh node 监控

一个简单的 web 监控 spacemesh 节点和 post 的工具。
定时读取 node 的 gRPC 服务，通过 ws 推送到页面。

编译：

```
  git clone https://github.com/gery2005001/sm-status.git
  cd sm-status
  go mod tidy

  Window下：
  .\scripts\makefile.bat
  .\sm-status.exe

  Linux下：
  sh ./scripts/makefile.sh
  ./sm-status

```

使用：http://localhost:8008

config.json 文件配置说明：

```
{
  port: 8008,       // web 服务的端口
  refresh: 120,     // 刷新post service的http状态间隔时间(s)
  interval: 10,    // ws推送间隔时间(s)
  timeout: 5,       // 获取状态超时时间(s)
  reload: true,     // 每次刷新状态时是否重新加载config文件
  "node":[
    {
      "name": "node01",        //Node主机名称
      "ip": "192.168.31.220",        //Node主机IP地址
      "grpc-public-listener": 9092,  //grpc-public-listener端口
      "grpc-private-listener": 9092, //grpc-private-listener端口
      "grpc-post-listener": 9092,    //grpc-post-listener端口
      "grpc-json-listener": 9092,    //grpc-json-listener端口
      "enable": true,                //Node是否启用
      "node-type": "multi",          //节点类型：multi=1:n节点，alone=go-spacemesh启动的节点,smapp=SMAPP启动的节点
      "post": [                      //挂在node下的Post server盘
        {
          "enable": true,           //是否启用
          "title": "67180",         //名称标识
          "capacity": "0.89TB",     //容量
          "operator-address": "http://192.168.31.220:50050/status"  //--operator-address参数定义的http状态地址
        },
        {
          "enable": true,
          "title": "DE61F",
          "capacity": "3.50TB",
          "operator-address": "http://192.168.31.220:50051/status"
        }
      ]
    },
    {
      "name": "node02-alone",
      "ip": "192.168.31.220",
      "grpc-public-listener": 9092,  //grpc-public-listener端口
      "grpc-private-listener": 9092, //grpc-private-listener端口
      "grpc-post-listener": 9092,    //grpc-post-listener端口
      "grpc-json-listener": 9092,    //grpc-json-listener端口
      "enable": true,
      "node-type": "alone",          //节点类型：multi=1:n节点，alone=go-spacemesh启动的节点,smapp=SMAPP启动的节点
      "post": [
        {
          "enable": true,
          "title": "67180",
          "capacity": "0.89TB",
          "address": "",
          "status": ""
        }
      ]
    },
    {
      "name": "node03-smapp",
      "grpc-public-listener": 9092,  //grpc-public-listener端口
      "grpc-private-listener": 9092, //grpc-private-listener端口
      "grpc-post-listener": 9092,    //grpc-post-listener端口
      "grpc-json-listener": 9092,    //grpc-json-listener端口
      "enable": true,
      "node-type": "smapp",          //节点类型：multi=1:n节点，alone=go-spacemesh启动的节点,smapp=SMAPP启动的节点
      "post": [
        {
          "enable": true,
          "title": "node03",
          "capacity": "1.75TB",
          "address": "",
          "status": ""
        }
      ]
    }
  ]
}

```
