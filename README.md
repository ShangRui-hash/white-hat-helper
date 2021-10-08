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


#### 语言选择：
- 后端：考虑go或者python
- 前端：考虑vue 或者react ，可以结合Echarts做数据可视化


##### 可以参考的开源项目：
- [NmapScaner](https://github.com/fuzz-security/NmapScaner/blob/master/scaner.sh)
这个项目用shell脚本实现 扫描+攻击
- [portMonitor](https://github.com/wantongtang/portMonitor)
这个基于nmap实现的一个端口监控程序，用于监控公司的端口

gui界面：nmap官方有zenmap，可以参考：
![avatar](https://img-blog.csdnimg.cn/20211008175010946.png?x-oss-process=image/watermark,type_ZHJvaWRzYW5zZmFsbGJhY2s,shadow_50,text_Q1NETiBA5peg5Zyo5peg5LiN5Zyo,size_20,color_FFFFFF,t_70,g_se,x_16)
![avatar](https://img-blog.csdnimg.cn/20211008175058296.png?x-oss-process=image/watermark,type_ZHJvaWRzYW5zZmFsbGJhY2s,shadow_50,text_Q1NETiBA5peg5Zyo5peg5LiN5Zyo,size_20,color_FFFFFF,t_70,g_se,x_16)
![avatar](https://img-blog.csdnimg.cn/20211008175154645.png?x-oss-process=image/watermark,type_ZHJvaWRzYW5zZmFsbGJhY2s,shadow_50,text_Q1NETiBA5peg5Zyo5peg5LiN5Zyo,size_20,color_FFFFFF,t_70,g_se,x_16)

idea:
我们可以做成那种任务式的扫描器。用户可以添加多个扫描任务。让多个任务同时去跑
后端提供一个接口，可以添加扫描任务。  
添加完扫描任务后，这个任务就在后台去跑。  
前端可以继续添加更多的扫描任务。  
后端向前实时反馈每个任务进行的进度。  