### tim是一个分布式聊天服务器

1. 支持im的基本功能。 <br/>
2. 支持群聊。 <br/>
3. 支持用户状态信息推送，"在线","离开"等。 <br/>
4. 支持消息回执.消息不丢失。 <br/>
5. 支持离线信息，聊天信息等存储与拉取。 <br/>
6. 通过协议拓展，可以支持视频，音频等通讯。 <br/>
7. 支持心跳检测异常断开的客户端，检测客户端验证超时等。 <br/>
8. 支持可配置的同一账号多客户端同时登陆。 <br/>
9. 可以配置关联其他数据库用户系统，tim可以独立于业务之外。 <br/>
10. 支持无数据库模式，此模式无法保存数据。 <br/>
11. 支持自定义消息类型，如发送位置，分享购物信息等。 <br/>
12. 支持分布式部署，服务器横向拓展。 <br/>
13. 支持TLS安全传输层协议。  <br/>
14. 支持hbase存储消息数据。  <br/>

目前实现的客户端有java，golang，Obj-C，由于使用thrift作为传输协议，开发相对简单许多，大部分只是简单接口调用即可，协议拓展也相对容易许多。 <br/>

实际应用场景：已在公司上线使用，原使用阿里悟空即时通讯，由于需求的增加，及部分消息类型需要定制，后来改用了tim即时通讯，基本可以满足当前以及以后的拓展需求。 <br/>

支持分布式部署是tim的重要功能也是许多使用im的开发者关心的问题。经过一段时间的思考，决定采用最简单有效的方式：tim的集群非常简单，只需额外提供一个redis服务即可，每个tim节点会从redis服务上读取其他tim节点的信息，使用thrift协议在服务器之间进行信息交互。tim与redis的交换也非常简单，存，读，删，除redis命令以外如果还有逻辑，则采用lua完成。集群的数量没有限制，理论上可以无限的扩展。 <br/>

tim服务器启动时需要指定配置文件 <br/>
如 tim  f  tim.xml  d  debug   c cluster.xml    <br/>
f后跟基本配置文件tim.xml ;  d后跟日志打印级别debug，info，warn，error等,可参考github.com/donnie4w/go-logger项目  ; c后跟集群文件cluster.xml <br/>

protocols文件夹中有几门语言的thrift协议，通过这协议可与服务器通信，通讯流程请在doc中查阅. <br/>

tim即时通讯项目虽然已经投入使用，但是目前的改动还是比较大，项目尚在不断完善当中！有任何建议或意见请随时email给我：donnie4w@gmail.com .谢谢！ <br/>
<br/>
体验demo <br/>
另外tim提供了windows环境的两个可执行文件，有兴趣的人可以玩一下，server.exe与client.exe。启动请看命令说明，主要是用户登录发送信息，状态信息等简单的体验。没有上下线的通知是由于用户关系没有建立。 <br/>
### v1.1
版本开始支持Hbase存储；对消息量比较大的服务，可以使用Hbase存储数据，hbase版本要求0.98以上,需开启hbase的thrift2服务。

### 客户端：
1. [java:https://github.com/donnie4w/atim](https://github.com/donnie4w/atim)
2. [go:https://github.com/donnie4w/timgo](https://github.com/donnie4w/timgo)
3. [Obj-C:https://github.com/3990995/tim-objc](https://github.com/3990995/tim-objc)
4. [kotlin:https://github.com/donnie4w/timkotlin](https://github.com/donnie4w/timkotlin)