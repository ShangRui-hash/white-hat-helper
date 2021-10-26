# 基于Nmap网络资产扫描

- 题目名称： 基于Nmap网络资产扫描
- 题目方向： 计科,软工
- 导师姓名： 杨力
- 导师电话： 13992822998
- 导师邮件： yangli@xidian.edu.cn

### 工程设计题目内容介绍：

为加强网络资产管理，需对网络资产进行扫描探测和拓扑描绘。具体要求如下：

##### 1.能针对网络设备进行主动和被动探测；

主动探测：
nmap基本都设计好了，主要涉及：

- 主机扫描，扫描一下网段里有哪些主机。扫描方法很多：比如发ping报文，如果是在内网，可以发送arp报文
- 端口扫描，扫描一下主机都开放了哪些端口。具体的扫描方法也很多：比如syn包扫描，全连接扫描，fin扫描，这些nmap都实现了，咱们只要fork一个子进程调用nmap就可以
- 端口服务扫描：扫描一下端口上跑的是什么服务。比如80端口跑的http,22端口跑的ssh。主要是通过建立tcp连接后 响应报文的banner信息来判断。这些nmap也替我们做好了

关于扫描的意义：扫描出来的主机、端口、服务都是我们的攻击面，比如说一旦发现对方开了ssh服务，我们就可以尝试比如ssh弱口令爆破之类的攻击手段。最好还能扫描出来服务的版本号，然后根据版本号把该版本的n day漏洞给找出来。这样我们就可以利用n day来尝试攻击对应的系统。唉，这种效果msf已经做了。

被动探测：

- 针对于内网中的目标：被动探测主要就是凭借嗅探。比如我们可以在混杂模式下监听流经我们的网卡的所有数据包。从这些数据包中我们可以知道内网中有哪些主机。
- 对于公网上的目标：被动探测可以是 公开情报搜集。比如调用shodan，fofa 这些站点的api。拿到他们已经收集好的目标的相关信息。

被动探测好处在与我们没有直接与目标发生交互。目标并不知道我们的存在。但是我们已经拿到的目标的相关信息。

##### 2.构建数据库用于网络设备资产信息的存取，数据库类型不做要求；

数据库可以考虑mysql，redis。缺点是：导致软件便携型下降。可以考虑封装成一个docker镜像。
也可以考虑 blot.db 这种文件型的数据库。优点是软件易于部署，便携型很高。

##### 3.设计实现前端界面，可视化探测过程，包括但不限于网段选择、探测结果显示、目标设备资产显示。

前端gui可以选用react或者vue 或者qt。为了可视化方便，可以考虑采用react和vue这类的web前端。然后前后端用websocket建立长连接。实现数据的实时交互。


##### 语言选择：
- 后端：考虑go或者python
- 前端：考虑vue 或者react ，可以结合Echarts做数据可视化


##### 可以参考的开源项目：
- [NmapScaner](https://github.com/fuzz-security/NmapScaner/blob/master/scaner.sh)
这个项目用shell脚本实现 扫描+攻击
- [portMonitor](https://github.com/wantongtang/portMonitor)
这个基于nmap实现的一个端口监控程序，用于监控公司的端口
- [linglong](https://github.com/awake1t/linglong)
  一款甲方资产巡航扫描系统。系统定位是发现资产，进行端口爆破。帮助企业更快发现弱口令问题。主要功能包括: 资产探测、端口爆破、定时任务、管理后台识别、报表展示
  ![avatar](index.gif)
##### 可以参考的资料：
- [Python+Django+AnsiblePlaybook自动化运维项目实战](https://coding.imooc.com/class/160.html)
- [诸神之眼——Nmap网络安全审计技术揭秘](https://item.jd.com/12165817.html?cu=true&utm_source=www.baidu.com&utm_medium=tuiguang&utm_campaign=t_2016327531_&utm_term=879f6bb2e77d4041aa459e049bb24c86)
- 

gui界面：nmap官方有zenmap，可以参考：
![avatar](https://img-blog.csdnimg.cn/20211008175010946.png?x-oss-process=image/watermark,type_ZHJvaWRzYW5zZmFsbGJhY2s,shadow_50,text_Q1NETiBA5peg5Zyo5peg5LiN5Zyo,size_20,color_FFFFFF,t_70,g_se,x_16)
![avatar](https://img-blog.csdnimg.cn/20211008175058296.png?x-oss-process=image/watermark,type_ZHJvaWRzYW5zZmFsbGJhY2s,shadow_50,text_Q1NETiBA5peg5Zyo5peg5LiN5Zyo,size_20,color_FFFFFF,t_70,g_se,x_16)
![avatar](https://img-blog.csdnimg.cn/20211008175154645.png?x-oss-process=image/watermark,type_ZHJvaWRzYW5zZmFsbGJhY2s,shadow_50,text_Q1NETiBA5peg5Zyo5peg5LiN5Zyo,size_20,color_FFFFFF,t_70,g_se,x_16)

##### idea:
1. 我们可以做成那种任务式的扫描器。用户可以添加多个扫描任务。让多个任务同时去跑,后端提供一个接口，可以添加扫描任务。添加完扫描任务后，这个任务就在后台去跑。前端可以继续添加更多的扫描任务。后端向前实时反馈每个任务进行的进度。 
2. 可以让这个项目为SRC漏洞挖掘服务。SRC应急响应平台一般会向白帽提供一个域名列表作为授权渗透测试的资产列表。如果我们的工具可以以公司为单元，来管理扫描到的资产，效果一定很赞。
   例如：联想的src 给的域名列表如下：
   ![avatar](https://img-blog.csdnimg.cn/20211008234136815.png?x-oss-process=image/watermark,type_ZHJvaWRzYW5zZmFsbGJhY2s,shadow_50,text_Q1NETiBA5peg5Zyo5peg5LiN5Zyo,size_20,color_FFFFFF,t_70,g_se,x_16)   
   我们将该域名列表输入工具，工具能自动根据这些域名扫描出联想公司具体有哪些资产，例如：
   - 根据域名列表通过多种手段查出来所有子域名。（比如通过https证书、通过字典爆破、通过公开信息收集）
   - 哪些域名用了cdn
   - 域名的真实ip是多少
   - 哪些站点看似独立，但是实际部署在同一个主机上。

##### 概要设计
目标：设计成一款为红队服务的资产梳理工具，为针对特定目标的渗透测试服务。

1.公司管理：
- 添加公司
- 删除公司
- 修改公司
- 查看所有公司列表
- 查看公司资产列表：点击公司可以查看该公司的资产列表（做成类似于fofa，shodan的显示效果）
- 查看资产详情：点击资产列表的某个资产进入到资产详情页，显示具体的资产信息。具体的资产信息包括：ip地址、域名、开放端口、端口上的服务信息、如果是web服务，显示子目录列表（title+http响应码）
- 渗透记录：可以给某个资产添加渗透测试记录  
参考效果：  
1.censys
![输入图片说明](https://images.gitee.com/uploads/images/2021/1026/183757_3bdb4565_8582605.png "屏幕截图.png")
2.shodan
shodan 列表页
![输入图片说明](https://images.gitee.com/uploads/images/2021/1026/184235_ecd7ed58_8582605.png "屏幕截图.png")
shodan 详情页
![输入图片说明](https://images.gitee.com/uploads/images/2021/1026/184350_83c3bf01_8582605.png "屏幕截图.png")

2.任务管理：对任务进行增删改查
- 添加扫描任务
    - 输入：所属公司、资产列表（域名列表或者ip列表）、主动扫描还是被动扫描
    - 输出：扫描进度，将扫描到的信息存储到数据库中
    - 扫描机制：扫描出所有子域名->扫描子域名的ip->扫描ip开放的端口->扫描端口运行的服务（服务版本)->如果是web服务，扫描所有的web目录，并获取title信息
- 删除扫描任务
- 暂停扫描任务
- 查看所有扫描任务（任务名，公司名，进度）
3.其他功能：
为了便于团队协作，应该再加一个登录功能，进行访问控制


团队分工：
张：前端
尚：后端
姜：shell脚本调度nmap