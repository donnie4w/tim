## TIM即时通讯引擎    [[English]](https://github.com/donnie4w/tim/blob/master/README.md "[English]")

> Tim即时通讯引擎的去中心化分布式架构具有去中心化、分布式数据存储、支持大规模用户、即时通讯、安全性和隐私保护、高可用性和容错性以及可扩展性和灵活性等特点。能够有效地解决大规模分布式系统的设计和实现问题，并提高系统的性能、可用性和扩展性。Tim即时通讯IM引擎是一个去中心化分布式架构，其主要特点在以下内容中详细介绍

### Tim实现的开源项目 webtim

- webtim在线访问地址：https://tim.tlnet.top

------------
### Tim的架构特点

1. 去中心化：Tim采用去中心化的分布式架构，没有中心节点或控制单元。每个节点都是独立的，具有高度自治的特征。这种架构方式可以降低单点故障的风险，提高系统的可靠性和安全性。
2. 分布式数据存储：Tim采用分布式存储算法，将数据分散存储在多个数据库节点上。这种存储方式可以提高数据的可靠性和安全性，并且能够更好地抵御数据丢失或损坏的问题。
3. 支持大规模用户：Tim能够支持大规模用户同时在线，并保证消息的及时送达。通过优化的协议和序列化技术，Tim能够高效地处理海量数据和请求，确保消息的快速传输和可靠交付。
4. 即时通讯：Tim作为即时通讯IM引擎，强调信息的及时性和可达性。它采用高效的传输协议和序列化技术，优化消息的传递和接收，确保用户能够即时地交流和互动。
5. 安全性和隐私保护：Tim重视用户信息的隐私和安全。通过加密技术和去中心化身份验证等手段，确保用户数据的安全和隐私不受侵犯。同时，内部开发人员也无法直接查询用户及通讯信息，增加了数据的安全性。
6. 高可用性和容错性：Tim具有高可用性和容错性。由于采用去中心化分布式架构，即使某个节点出现故障，其他节点仍能继续工作，确保整体服务的连续性。这种设计方式提高了系统的可靠性和稳定性。
7. 可扩展性和灵活性：Tim的分布式架构使其具有良好的可扩展性和灵活性。随着业务需求的增长，可以增加更多的节点来提高系统的处理能力和存储容量。同时，由于节点间的自由连接和交互，Tim能够快速适应变化的需求和环境。

------------

### Tim的通讯模式

tim是一个去中心化的分布式即时通讯引擎。

tim实现完全无中心的集群模式，并实现分布式数据存储，支持百万台数据库分布式存储数据，所以tim可以支持超大规模的集群，支持在线用户量可以达到百亿级别。

tim 构建了基础通讯模式，即1:1，1:N，N:N 模式，实现了底层流数据通讯。开发者根据业务需求灵活组合并使用tim的通讯模式，可以实现各种即时通讯的需求，类似微信，QQ，抖音等等im通讯，都可以用tim实现。类似直播间，实时音视频，多人视频，多人音视频会议等等功能，使用tim，非常简单就可以实现。类似消息撤回，消息阅后即焚等等，tim从底层消息通信类型就支持实现。

tim的使用主要依靠终端通过tim客户端调用服务接口，所有通讯逻辑都在tim服务端完成，tim客户端是采用极简模式实现接口调用，基本所有图灵完备的编程语言都可以快速实现tim客户端。


### tim功能简介

1. 去中心化分布式架构，支持Nat穿透集群节点，天然分布式架构，无需特殊配置，无最小节点数限制，零依赖，支持大规模集群。
2. 支持多种数据库：TLDB，Mysql，PostgreSQL ，SQL Server，Oracle，Oceanbase等
3. 高度归纳IM通讯模式，从底层支持1:1，1:N，N:N 通讯模式
4. 从底层支持流数据发送，轻松支持直播，实时音视频等功能的开发
5. 重视数据安全，从账号到通讯数据，都进行换算或加密，保证用户信息安全
6. 实现分布式存储数据，解决海量数据存储的问题，支持动态扩容数据库节点。
7. 具备单机高性能特点，协议序列化效率高，体积小。
8. 从底层支持各种通讯类型，包括普通信息，撤回信息，信息阅后即焚等
9. 支持各种自定义用户状态，对标QQ等类型APP的状态功能
10. 支持群组，对标QQ，微信等类型APP群组功能
11. 支持多人实时流数据发送，对标抖音直播，视频连线，或多人实时视频会议等app的功能
12. 支持同账号多终端同时登录，并支持通过配置进行限制，对标QQ，微信多终端登录功能
13. 支持thrift压缩，json，big string，big binary等多种形式的客户端协议
 

### tim的数据库

tim的数据库可以使用 TLDB，Mysql，PostgreSQL ，SQL Server，Oracle，Oceanbase等数据库。

TLDB的使用可以参考《TLDB 高性能分布式数据库》

TLDB是tim内置用户系统的默认数据库。通过tim自身的数据分布式存储设计，tldb无需再搭建分布式系统，只需要启动单机模式的服务节点即可。tim通过数据分库配置，可以把数据散列存储到多个单机模式的tldb服务节点中。

TLDB的高性能读写数据与支持大量客户端连接的特点，相对来说更适合作为tim数据库

如果不使用tim内置用户系统，而是需要接入已经存在或自己构建的用户系统时，可以接入Mysql等其他数据库，通过配置SQL的方式来接入外部的业务数据。tim的核心接口支持读取相关的外部数据



### tim的协议

tim自定义的通讯协议，在序列化效率与压缩比例上，相对于常见序列化框架，都具备较大的优势，详细请参考《tim实践系列——tim协议与其他格式协议的比较》

tim除了支持自定义的thrift压缩协议，同时支持也支持json格式协议；json本身的优缺点都比较明显，json最大优势是其通用性，而json的序列化效率与协议包的体积是其劣势。timjs是采用json格式实现的客户端，详细参考timjs实现源码

tim除了支持有固定格式的协议，同时支持无格式的大字符串(big string),大字节数组(big binary)协议，该两种协议无需通过序列化与反序列化操作，没有固定的格式，适用所有终端。部分终端如浏览器无法使用大体积的json格式时，可以使用big string或big binary协议传输。而且big string或big binary支持协议包被拆包发送的复杂场景。




### tim零依赖

tim的部署与启动不依赖任何第三方组件或服务。tim支持有数据模式与无数据模式。在有数据模式下，运行tim需要先启动数据库服务，在无数据模式下，直接运行即可。

------------

- [TIM开发使用文档](https://tlnet.top/timdoc "TIM开发使用文档")
- [TIM源码地址](https://github.com/donnie4w/tim "TIM源码地址")
- [在线体验](https://tim.tlnet.top/ "在线体验")
- [TIM下载地址](https://tlnet.top/download "TIM下载地址")

------------
### tim实践系列文章 (文章持续更新中...)

- [《tim实践系列——tim协议与其他格式协议的比较》](https://tlnet.top/article/22425142 "《tim实践系列——tim协议与其他格式协议的比较》")
- [《tim实践系列——tim设计来源与设计模式》](https://tlnet.top/article/22425137 "《tim实践系列——tim设计来源与设计模式》")
- [《tim实践系列——如何使用TimMessage自定义各种消息》](https://tlnet.top/article/22425173 "《tim实践系列——如何使用TimMessage自定义各种消息》")
- [《tim实践系列——如何使用TimPrecence自定义各种用户状态》](https://tlnet.top/article/22425172 "《tim实践系列——如何使用TimPrecence自定义各种用户状态》")
- [《tim实践系列——虚拟房间的作用和如何使用》](https://tlnet.top/article/22425182 "《tim实践系列——虚拟房间的作用和如何使用》")
- [《tim实践系列——用户如何实现 隐身，在线，忙碌等状态》](https://tlnet.top/article/22425175 "《tim实践系列——用户如何实现 隐身，在线，忙碌等状态》")
- [《tim实践系列——消息特点和如何在实际业务中使用》](https://tlnet.top/article/22425174 "《tim实践系列——消息特点和如何在实际业务中使用》")
- 《[tim实践系列——tim如何限制一个账号多个终端登录》](https://tlnet.top/article/22425175 "tim实践系列——tim如何限制一个账号多个终端登录》")
- 《[tim实践系列——tim信息安全与账号系统》](https://tlnet.top/article/22425171 "tim实践系列——tim信息安全与账号系统》")
- [《tim实践系列——内置的好友关系和群组》](https://tlnet.top/article/22425170 "《tim实践系列——内置的好友关系和群组》")
- [《tim实践系列——消息撤回，阅后即焚 等功能如何开发》](https://tlnet.top/article/22425181 "《tim实践系列——消息撤回，阅后即焚 等功能如何开发》")
- [《tim实践系列——去中心化分布式架构特点》](https://tlnet.top/article/22425179 "《tim实践系列——去中心化分布式架构特点》")
- [《tim实践系列——分布式数据存储与动态数据库扩容》](https://tlnet.top/article/22425176 "《tim实践系列——分布式数据存储与动态数据库扩容》")
- [《tim实践系列——如何架构支持百亿级别在线用户的即时消息系统》](https://tlnet.top/article/22425157 "《tim实践系列——如何架构支持百亿级别在线用户的即时消息系统》")
- [《tim实践系列——接入外部账号系统，配置关系型数据库》](https://tlnet.top/article/22425156 "《tim实践系列——接入外部账号系统，配置关系型数据库》")
- [《tim实践系列——tim的限流，报文长度，连接数，请求频率》](https://tlnet.top/article/22425178 "《tim实践系列——tim的限流，报文长度，连接数，请求频率》")
- 《tim实践系列——后台管理员接口的使用》
- 《tim实践系列——账号安全措施》
- 《tim实践系列——用户之间如何实现文件传输》
- 《tim实践系列——微信实时音视频开发》
- 《tim实践系列——抖音直播间的开发》
- 《tim实践系列——多人实时视频会议，抖音直播间在线多人视频连线》
- 《tim实践系列——如何实现对公众号订阅号等系统栏目的实时订阅》
- 《tim实践系列——tim产生的数据统计与建议》
- 《tim实践系列——tim配置系统参数要注意的地方》
- 《tim实践系列——使用tim无数据库模式实现web IM》

------------
### tim 相关工程项目

- go客户端            timgo： https://github.com/donnie4w/timgo
- java客户端          atim： https://github.com/donnie4w/atim
- js客户端              timjs：https://github.com/donnie4w/timjs
- 后台接口示例      admintim：  https://github.com/donnie4w/admintim
- webtim项目        https://github.com/donnie4w/webtim      访问地址： https://tim.tlnet.top
- tim客户端协议  [tim-protocol](https://github.com/donnie4w/tim-protocol "tim-protocol")
- [《tim实践系列文章》](https://github.com/donnie4w/Tim-Practical-Article "《tim实践系列文章》")


------------

### tim带来多方面的优势和作用

1. 提高开发效率：TIM提供了丰富的功能模块和接口，简化了即时通讯应用的开发过程，减少了重复造轮子和从头开始摸索的时间。
2. 保障安全性能：TIM具备强大的安全机制，能够保障用户数据和通信内容的安全性，防止数据泄露和恶意攻击。
3. 提供稳定服务：TIM能够提供稳定可靠的即时通讯服务，确保用户沟通的连续性和可靠性。
4. 支持多种平台：TIM可以跨平台、跨设备使用，支持多种操作系统和终端设备，满足不同用户的需求。
5. 丰富的扩展性：去中心化分布式的TIM具备超强扩展性，能够根据项目的需求进行水平扩展，方便进行二次开发和功能升级。
6. 降低维护成本：TIM可以降低大量运维成本，TIM集群，与分布式数据存储都依赖TIM本身算法完成，无需人为部署。

------------

- tim可以极大降低使用IM的成本，提高im开发效率。

- 相信tim必定可以更加简单快捷的解决im相关问题，有任何疑问请email：donnie4w@gmail.com