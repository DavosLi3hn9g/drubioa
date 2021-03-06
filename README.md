
<img src="https://iqiar.com/uploads/logo_git.png" width=300 />

以其人之道，还治其人之身，

用骚扰电话之武器来反击骚扰电话，

彻底解决电话骚扰问题的AI解决方案。

## 项目简介

这是一个软硬件结合的开源项目，如果是普通使用者，需要有初级的硬件动手能力（比如接电源网线，插个板子总得会）。

初版主要针对手机自带的呼叫转移功能，您可以将手机设置白名单策略，将所有陌生电话呼叫转移到AI接听。

可全程使用您自己可控的云服务和硬件，无需担心隐私泄露。

### 已实现功能：

**智能接听：** 识别对方的语音，配合您的话术设置判断是否骚扰行为（比如设置拒接"贷款、中介、炒股等"骚扰来电）

**多种挂断策略：** 发现骚扰行为并挂断后，自动发送短信通知到您的主号。

**短信代收代发：** 不仅通知信息，还会将副号收到的短信全部转发给您的主号，从此告别备用手机。

**人性化交互：** 几乎零代码部署，自带WebUI用户交互界面。

**AT控制台：** 管理中心可以直接使用AT命令与SIM扩展板交互。（针对高级用户）

**新手向导：** 初次使用也能轻易上手，玩转AI就跟设置路由器一样简单。（针对初级用户）

**接入阿里云：** 已经适配阿里云，稍作配置即可接入智能语音服务。

**一键更新：** 自带版本更新功能，一劳永逸。


## 项目原理：

<img src="https://iqiar.com/uploads/yuanli.jpg" width=75% />

## 准备材料

**1 . 一个树莓派(推荐2代以上版本，ZERO版暂未测试) 或者 一台有USB接口的PC（不推荐使用PC，仅供尝鲜）**

系统方面推荐刷写官方的Raspbian，Windows版还没有完整测试。

**2 . 一块SIM扩展板，和一张备用的SIM电话卡**

注意，SIM扩展板需支持PCM电话语音读写，否则只能收发短信，无法智能语音。

目前确定可用的扩展板，可在淘宝搜索 微雪SIM7600X，其他扩展板还有待测试。

**3 . 一个已经实名的阿里云开发者帐号**

根据反馈，再看是否要适配其他云服务。

## 硬件接线图

按照下图所示连接好SIM扩展板和树莓派：

<img src="https://iqiar.com/uploads/SIM7600X.jpg" width=600  />

## 开始使用

1 . ssh登录树莓派设备

`sudo raspi-config`
`然后选择Interfacing Options ->Serial ->no -> yes，关闭串口调试功能。`

2 . 下载zip包并解压：

`$ wget -O qiar.zip https://iqiar.com/qiar_armv6.zip && unzip -o ./qiar.zip -d ./QiarAI && cd ./QiarAI` 

3 . 赋予执行权限，开始运行：

`$ chmod 777 ./start.sh`

`$ ./start.sh`

停止运行：

`$ ./stop.sh`

## 访问WEB控制台

地址：http://您的设备IP:5000

## 后续版本规划

探索纯软件方案，比如借助目前骚扰电话常用的虚拟号码策略。

增加语音通话实时转发功能，实现更彻底的智能电话托管功能。


## 常见问题

反馈讨论QQ群：112997264

#### 无法启动？请检查日志：

后台运行的日志：./nohup.log

#### 开发者要注意什么？

编辑 ./configs/.env

`dev_mode = true` //启用开发者模式，日志里面会看到更多信息

`allow_origin = http://开发机IP:5000` 
//开启调试模式后，注意要设置跨域白名单

### 本项目主要依赖以下第三方开源库，感谢他们！

[serial](http://github.com/tarm/serial), [gin](http://github.com/gin-gonic/gin), [go-wav](http://github.com/youpy/go-wav), [websocket](http://github.com/gorilla/websocket), [go-update](http://github.com/inconshreveable/go-update), [vuejs](http://cn.vuejs.org/), [ant-design](http://github.com/vueComponent/ant-design-vue) ...

## 项目截图

<img src="https://iqiar.com/uploads/pic_1.jpg" width=75% />
