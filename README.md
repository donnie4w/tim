## TIM   IM Engine    [[中文]](https://github.com/donnie4w/tim/blob/master/README_zh.md "[中文]")

> preface：tim2.x  vs  tim1.x

> There is no direct relationship between tim2. x and tim1. x. tim1. x is an early implementation of im, while tim2. x is a redesigned and developed program. The implementation of tim1. x has many fundamental design issues that prevent further expansion and optimization based on the program. All code and design solutions for tim1. x can only be discarded, and redesigned and implemented.

> Tim2. x achieved the expected results and functional efficiency. It can solve all the IM needs currently encountered. And it has strong scalability, and the project will continue to carry out corresponding functional expansion and performance optimization in the future.

#### [Tim Official Website](https://tlnet.top/timen "Tim Official Website")

------------

### The characteristics of Tim

Tim is a decentralized distributed instant messaging engine.

Tim implements a completely decentralized cluster mode and distributed data storage, supporting millions of databases for distributed data storage. Therefore, Tim can support ultra large clusters and online user numbers can reach ten billions.

Tim has built a basic communication mode, namely 1:1, 1: N, N: N mode, to achieve underlying stream data communication. Developers can flexibly combine and use tim's communication mode according to business requirements to meet various instant messaging needs. IM communications such as WeChat, QQ, Tiktok, etc. can be implemented with tim. Similar to live streaming rooms, real-time audio and video, multiplayer videos, multiplayer audio and video conferences, and other functions, using Tim is very simple to implement. Similar to message recall, message burning upon reading, etc., Tim supports implementation from the underlying message communication type.

The use of Tim mainly relies on the terminal calling the service interface through the Tim client, and all communication logic is completed on the Tim server. The Tim client adopts a minimalist mode to implement interface calls, and almost all Turing complete programming languages can quickly implement the Tim client.



### Tim Function Introduction

1. Decentralized distributed architecture, supporting Nat penetrate to connect to cluster nodes, natural distributed architecture, no special configuration required, no minimum node limit, zero dependency, supporting large-scale clusters.
2. Supports multiple databases: TLDB, MySQL, PostgreSQL, SQL Server, Oracle, Oceanbase, etc
3. Highly inductive IM communication mode, supporting 1:1, 1: N, N: N communication modes from the bottom layer
4. Supporting streaming data transmission from the bottom layer, easily supporting the development of live streaming, real-time audio and video functions, etc
5. Emphasize data security , from account to communication data, are converted or encrypted to ensure the security of user information
6. Implement distributed storage of data, solve the problem of massive data storage, and support dynamic expansion of database nodes.
7. Featuring high performance of a stand-alone machine, high protocol serialization efficiency, and small size.
8. Support various communication types from the bottom level, including regular messages, recall messages, and messages that are burned upon reading, etc
9. Support various custom user statuses and benchmark the status function of QQ and other types of apps
10. Support group, benchmarking QQ, WeChat and other types of APP group functions
11. It supports multi person real-time streaming data transmission, benchmarking Tiktok live broadcast, video connection, or multi person real-time video conference ,etc.
12. Supports simultaneous login of multiple terminals with the same account, and supports restriction through configuration, benchmarking QQ and WeChat multiple terminal login functions
13. Support client access using JSON protocol.


### Tim's database

Tim's database can use databases such as TLDB, MySQL, PostgreSQL, SQL Server, Oracle, Oceanbase, etc.

The use of TLDB can refer to "[TLDB High Performance Distributed Database](https://tlnet.top/tldben "TLDB High Performance Distributed Database")"

TLDB is the default database for Tim's built-in user system. Through Tim's own data distributed storage design, TLDB no longer needs to build a distributed system, only needs to start a standalone machine mode service node. Tim can hash and store data in multiple standalone TLDB service nodes through data repository configuration.

The high-performance read and write data and support for a large number of client connections of TLDB make it relatively more suitable as a Tim database

If you do not use Tim's built-in user system, but need to access existing or self built user systems, you can access other databases such as MySQL and configure SQL to access external business data. Tim's core interface supports reading relevant external data



### Tim's protocol

Tim's custom communication protocol has significant advantages over common serialization frameworks in terms of serialization efficiency and compression ratio. For more details, please refer to the "Tim Practice Series - Comparison of Tim Protocol with Other Format Agreements"

Tim not only supports custom Thrift compression protocol, but also supports JSON protocol; The advantages and disadvantages of JSON itself are quite obvious. The biggest advantage of JSON is its versatility, while its serialization efficiency and protocol package size are its disadvantages. Timjs is a tim client implemented in JSON format, please refer to the timjs implementation source code for details



Tim zero dependency
Tim's deployment and startup do not rely on any third-party components or services. Tim supports both data mode and no data mode. In data mode, running Tim requires starting the database service first. In no data mode, it can be run directly.



### [TIM development and usage documentation](https://tlnet.top/timdoc "TIM development and usage documentation")

------------

- **Tim can greatly reduce the cost of developing IM and improve the efficiency of IM development.**

- **I believe Tim can easily and quickly solve IM related problems. If you have any questions, please email: donnie4w@gmail.com**