## TIM   IM Engine    [[中文]](https://github.com/donnie4w/tim/blob/master/README_zh.md "[中文]")

> The decentralized distributed architecture of Tim's instant messaging engine features decentralization, distributed data storage, support for large-scale users, instant messaging, security and privacy protection, high availability and fault tolerance, as well as scalability and flexibility. It can effectively solve the design and implementation problems of large-scale distributed systems and improve system performance, availability, and scalability. The IM engine of Tim is a decentralized distributed architecture, whose main features are detailed in the following content

------------

#### Tim's open source project webtim

Webtim online access address: https://tim.tlnet.top

------------

### Features of Tim's architecture

1. Decentralization: Tim adopts a decentralized distributed architecture without a central node or control unit. Each node is independent and has a highly autonomous characteristic. This architecture approach can reduce the risk of single point failure and improve the reliability and security of the system.
2. Distributed data storage: Tim uses a distributed storage algorithm to store data on multiple database nodes in a decentralized manner. This storage method can improve data reliability and security, and better protect against data loss or corruption issues.
3. Support for large-scale users: Tim can support large-scale users online at the same time and ensure timely delivery of messages. Through optimized protocols and serialization techniques, Tim can efficiently process massive amounts of data and requests, ensuring fast transmission and reliable delivery of messages.
4. Instant messaging: Tim as an instant messaging IM engine emphasizes the timeliness and accessibility of information. It uses efficient transmission protocols and serialization techniques to optimize the delivery and reception of messages, ensuring that users can communicate and interact instantly.
5. Security and privacy protection: Tim attaches great importance to the privacy and security of user information. Through encryption technology and decentralized authentication, it ensures the security and privacy of user data from being infringed upon. At the same time, internal developers cannot directly query user and communication information, which increases data security.
6. High availability and fault tolerance: Tim has high availability and fault tolerance. Due to the decentralized distributed architecture, even if a node fails, other nodes can continue to work, ensuring the continuity of the overall service. This design approach improves the reliability and stability of the system.
7. Scalability and flexibility: Tim's distributed architecture enables it to have good scalability and flexibility. As business needs grow, more nodes can be added to improve the system's processing power and storage capacity. At the same time, due to the free connection and interaction between nodes, Tim can quickly adapt to changing needs and environments.

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


### [Tim Practice Series](https://github.com/donnie4w/Tim-Practical-Article "Tim Practice Series") (articles are continuously updated...)

- [Tim Practice Series - Comparison of Tim Protocol and Other Formats Protocol](https://tlnet.top/article/22425142 "Tim Practice Series - Comparison of Tim Protocol and Other Formats Protocol")
- [Tim Practice Series - Tim Design Source and Design Pattern](https://tlnet.top/article/22425137 "Tim Practice Series - Tim Design Source and Design Pattern")
- [Tim Practice Series - How to Customize Various Messages Using TimMessage](https://tlnet.top/article/22425173 "Tim Practice Series - How to Customize Various Messages Using TimMessage")
- [Tim Practice Series - How to Customize Various User States Using TimPrecence](https://tlnet.top/article/22425172 "Tim Practice Series - How to Customize Various User States Using TimPrecence")
- [Tim Practice Series - The Role of Virtual Rooms and How to Use Them](https://tlnet.top/article/22425182 "Tim Practice Series - The Role of Virtual Rooms and How to Use Them")
-[ Tim Practice Series - How Users Can Achieve Invisible, Online, Busy, and Other Statuses](https://tlnet.top/article/22425175 " Tim Practice Series - How Users Can Achieve Invisible, Online, Busy, and Other Statuses")
- [Tim Practice Series - Message Features and How to Use in Real Business](https://tlnet.top/article/22425174 "Tim Practice Series - Message Features and How to Use in Real Business")
- [Tim Practice Series - How Tim Restricts Multiple Terminal Logins for One Account](https://tlnet.top/article/22425168 "Tim Practice Series - How Tim Restricts Multiple Terminal Logins for One Account")
- [Tim Practice Series - Tim Information Security and Account System](https://tlnet.top/article/22425171 "Tim Practice Series - Tim Information Security and Account System")
- [Tim Practice Series - Built-in Friendships and Groups](https://tlnet.top/article/22425170 "Tim Practice Series - Built-in Friendships and Groups")
- [Tim Practice Series - How to Develop Functions such as Message Withdrawal and Burn After Reading](https://tlnet.top/article/22425181 "Tim Practice Series - How to Develop Functions such as Message Withdrawal and Burn After Reading")
-[ Tim Practice Series - Characteristics of Decentralized Distributed Architecture](https://tlnet.top/article/22425179 " Tim Practice Series - Characteristics of Decentralized Distributed Architecture")
- [Tim Practice Series - Distributed Data Storage and Dynamic Database Expansion](https://tlnet.top/article/22425176 "Tim Practice Series - Distributed Data Storage and Dynamic Database Expansion")
- [Tim Practice Series - How to Build an Instant Messaging System Supporting Billions of Online Users](https://tlnet.top/article/22425157 "Tim Practice Series - How to Build an Instant Messaging System Supporting Billions of Online Users")
- [Tim Practice Series - Accessing External Account Systems and Configuring Relational Databases](https://tlnet.top/article/22425156 "Tim Practice Series - Accessing External Account Systems and Configuring Relational Databases")
- [Tim Practice Series - Tim's Limiting Flow, Message Length, Connection Count, Request Frequency](https://tlnet.top/article/22425178 "Tim Practice Series - Tim's Limiting Flow, Message Length, Connection Count, Request Frequency")
- Tim Practice Series - Using the Backend Administrator Interface
- Tim Practice Series - Account Security Measures
- Tim Practice Series - How to Implement File Transfer Between Users
- Tim Practice Series - WeChat Real-time Audio and Video Development
- Tim Practice Series - Development of Tiktok Studio
- Tim Practice Series - Multi person Real time Video Conference, Tiktok Live Room Online Multi person Video Connection
- Tim Practice Series - How to Realize Real time Subscription of System Columns such as official account Subscription Number
- Tim Practice Series - Statistics and Suggestions on Tim Production
- Tim Practice Series - Tim Configuration System Parameters to Pay Attention to
- Tim Practice Series - Implementing Web IM Using Tim No-Database Mode

------------

- [TIM development and usage document](https://tlnet.top/timdoc "TIM development and usage document")
- [ TIM source code address](https://github.com/donnie4w/tim " TIM source code address")
- [ Online experience](https://tim.tlnet.top/ " Online experience")
- [Download address of TIM](https://tlnet.top/download "Download address of TIM")

------------
### Programs related to Tim

- go client            timgo： https://github.com/donnie4w/timgo
- java client          atim： https://github.com/donnie4w/atim
- js client              timjs：https://github.com/donnie4w/timjs
- Management interface Example      admintim：  https://github.com/donnie4w/admintim
- webtim project         https://github.com/donnie4w/webtim       Access address： https://tim.tlnet.top

------------
### tim brings many advantages and functions

1. Improve development efficiency: TIM provides rich functional modules and interfaces, simplifying the development process of instant messaging applications and reducing the time spent on reinventing the wheel and starting from scratch.
2. Security: TIM has a strong security mechanism that can ensure the security of user data and communication content, preventing data leakage and malicious attacks.
3. Provide stable services: TIM can provide stable and reliable instant messaging services to ensure the continuity and reliability of user communication.
4. Support for multiple platforms: TIM can be used across platforms and devices, supporting multiple operating systems and terminal devices to meet the needs of different users.
5. Rich scalability: Decentralized and distributed TIM has super scalability, which can be horizontally expanded according to the needs of the project, facilitating secondary development and functional upgrades.
6. Reduce maintenance costs: TIM can reduce a large amount of operation and maintenance costs. The TIM cluster and distributed data storage rely on the TIM algorithm itself to complete, without requiring manual deployment.

------------

- Tim can greatly reduce the cost of developing IM and improve the efficiency of IM development.
- I believe Tim can easily and quickly solve IM related problems. If you have any questions, please email: donnie4w@gmail.com