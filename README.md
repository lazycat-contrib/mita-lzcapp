# mita — mieru 代理服务端（LazyCat 喵喵商店）

mieru【見える】是一款无流量特征、难以主动探测的 socks5 / HTTP / HTTPS 代理软件。
本应用部署其服务端 **mita【見た】**：跑在懒猫微服上作为代理服务器，
桌面与手机的 mieru / Clash 系客户端连回家里网络。

## 来源

- 上游：https://github.com/enfein/mieru （GPL-3.0）
- 镜像：上游官方镜像 `ghcr.io/enfein/mita`，按版本 tag 跟踪（当前 v3.40.0），经镜像站 `ghcr.1ms.run` 拉取
- 打包：`ca-x/lazycat-github-action`，上游发新版后自动跟随发布

## 工作方式

- 部署时设置代理端口、用户名、密码（密码随机生成，可在部署页改）
- `setup_script` 渲染 `/lzcapp/var/mita/server.json`，`entrypoint.sh` 首次启动时
  通过 `mita apply config` 写入（mita 内部转为哈希存储，明文随即删除）
- TCP ingress 映射代理端口；可选 UDP ingress（同一端口号）
- 无 Web 管理界面：应用入口页是连接说明（contentdir 静态页）
- 修改配置：应用终端编辑 `/lzcapp/var/mita/server.json` 后
  `mita apply config ... && mita reload`

## 客户端连接

地址 = 懒猫主机（局域网 IP 或公网域名），端口 = 部署设置的代理端口，
协议 = mieru（TCP/UDP），用户名密码 = 部署设置。
客户端推荐见应用内说明页。
