<img src="https://iqiar.com/uploads/logo_git.png" width=300 />

## 项目简介

这是一个软硬件结合的开源项目，如果是普通使用者，需要有初级的硬件动手能力（比如接电源网线，插个板子总得会）。

控制台地址：http://您的设备IP:5000

## 准备材料

**1 . 一个树莓派(推荐2代以上版本，ZERO版暂未测试) 或者 一台有USB接口的PC（不推荐使用PC，仅供尝鲜）**

系统方面推荐刷写官方的Raspbian，Windows版还没有完整测试

**2 . 一块SIM扩展板，和一张备用的SIM电话卡**

注意，SIM扩展板需支持PCM电话语音读写，否则只能收发短信，无法智能语音。

目前确定可用的扩展板，可在淘宝搜索 微雪SIM7600X，其他扩展板还有待测试。

**3 . 按照下图所示连接好SIM扩展板和树莓派**

<img src="https://iqiar.com/uploads/SIM7600X.jpg" width=600  />


## 开始使用

1 . ssh登录树莓派设备

2 . 下载zip包并解压：

`$ wget -O qiar.zip https://iqiar.com/qiar_armv6.zip && unzip -o ./qiar.zip -d ./QiarAI && cd ./QiarAI` 

3 . 赋予执行权限，开始运行：

`$ chmod 777 ./start.sh`

`$ ./start.sh`

停止运行：

`$ ./stop.sh`

## 访问WEB控制台

地址：http://您的设备IP:5000

## 常见问题

#### 无法启动？请检查日志：

后台运行的日志：./nohup.log

#### 开发者要注意什么？

编辑 ./configs/.env

`dev_mode = true` //启用开发者模式，日志里面会看到更多信息

`allow_origin = http://开发机IP:5000` 
//如果您要修改VUE源代码，调试时注意要设置跨域白名单

### 本项目主要依赖以下第三方开源库，感谢他们！

[serial](http://github.com/tarm/serial), [gin](http://github.com/gin-gonic/gin), [go-wav](http://github.com/youpy/go-wav), [websocket](http://github.com/gorilla/websocket), [go-update](http://github.com/inconshreveable/go-update), [vuejs](http://cn.vuejs.org/), [ant-design](http://github.com/vueComponent/ant-design-vue) ...

##项目截图

<img src="https://iqiar.com/uploads/pic_1.jpg" width=75% />