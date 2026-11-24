#!/bin/sh
# mita 服务端启动脚本：
# 1. 启动 mita daemon（后台 run，通过 UDS 提供 RPC）。
# 2. 首次启动（持久配置目录里还没有已应用的配置时）应用部署时渲染的 server.json；
#    server.json 模板保留在持久目录中，用户可随时编辑后
#    `mita apply config /lzcapp/var/mita/server.json && mita reload` 生效。
# 注意：mita 的真实状态文件是 /etc/mita/server.conf.pb（lzc-manifest 已把
# /lzcapp/var/mita/etc 挂载到 /etc/mita），重启容器配置不丢失。
set -e

VAR_DIR=/lzcapp/var/mita
CONFIG_FILE="$VAR_DIR/server.json"

mkdir -p "$VAR_DIR" "$VAR_DIR/etc" /var/run/mita
# mita 镜像内置用户 mita（uid 1000）需要对状态目录有写权限
chown -R mita:mita "$VAR_DIR/etc" /var/run/mita 2>/dev/null || true
chmod 775 "$VAR_DIR/etc" /var/run/mita 2>/dev/null || true

# 启动 daemon（后台），RPC 通过 UDS
mita run &
MITA_PID=$!

# 等 UDS 就绪（最多 15s）
i=0
while [ $i -lt 30 ]; do
  if [ -S /var/run/mita/mita.sock ]; then
    break
  fi
  i=$((i+1))
  sleep 0.5
done

# 首次启动：把渲染好的配置应用到 mita（密码由 mita 转为哈希保存在 /etc/mita/server.conf.pb）
NEED_APPLY=0
if mita describe config 2>/dev/null | grep -q '"portBindings"'; then
  NEED_APPLY=0
else
  NEED_APPLY=1
fi
if [ "$NEED_APPLY" = "1" ] && [ -f "$CONFIG_FILE" ]; then
  echo "[entrypoint] applying initial server config"
  if mita apply config "$CONFIG_FILE"; then
    mita start || true
  else
    echo "[entrypoint] WARNING: failed to apply initial config; will retry on next boot"
  fi
else
  # 已有配置：直接启动代理服务
  mita start || true
fi

# 前台等待 daemon 进程
wait $MITA_PID
