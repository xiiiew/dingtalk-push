# dingtalk-push

钉钉消息推送服务。支持推送钉钉text、link、markdown、actionCard及feedCard消息类型。为了防止消息推送频率超过限制，支持根据每个机器人限制推送频率和合并推送消息。

### Docker部署

```shell
docker build -t dingtalk-push:latest .
docker run -d -p 8080:8080 --restart always -v /data/log/dingtalk-push/logs:/app/logs -v /etc/localtime:/etc/localtime:ro -e "TZ=Asia/Shanghai" --name dingtalk-push dingtalk-push:latest
```

### SDK

`dingtalk-push`为了方便调用，支持如下版本SDK：

[Golang](https://github.com/xiiiew/dingtalk-push-golang-sdk)

### 配置文件

```yaml
app:
  app_name: dingtalk-push

log:
  # 日志文件路径
  logs_dir: ./logs
  # 日志文件分割频率(小时)
  logs_rotation_time: 1
  # 日志文件保存个数
  logs_rotation_count: 100

http:
  # 监听端口
  server_listen_port: 28080
  # 跨域
  cors:
    # 启用跨域配置
    enable: false
    access_control_allow_origin: '*'
    access_control_allow_methods: POST,GET,OPTIONS,DELETE,HEAD,PUT
    access_control_allow_headers: Authorization,Content_Type,Accept,Origin,User_Agent,DNT,Cache_Control,X_Mx_ReqToken,X_Data_Type,X_Requested_With, X_Data_Type,X_Auth_Token,token,language,Pragma

dingtalk:
  # 发消息超时时间
  time_out: 10
  # 消息间隔时间(单位：秒)
  message_duration: 3
  # 多条消息分界线
  boundary: '-----一条华丽的分割线-----'
```

`dingtalk-push`支持修改消息发送间隔时间以避免超过钉钉限制，被拉黑。截止目前，每个钉钉机器人消息频次每分钟限制在60次，建议`dingtalk.message_duration`配置不要低于`3`。
若同一个机器人在`dingtalk.message_duration`间隔内收到多条消息，多条消息会被合并为一条消息发送。`dingtalk.boundary`设置消息合并为同一条消息后，多条消息之间的分界线，用以方便查看。

### 消息限频

所有类型的消息都将执行限频策略。每个secret对应的机器人在单位频率区间内只会发送一条消息。单位频率区间收到的同一可合并类型的消息将会合并成一条消息，并会在某个频率区间统一发送到钉钉。