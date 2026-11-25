#!/bin/bash
# 构建 renderpage：
# 1. pnpm build 前端（~/mita-lzcapp/web）→ web/dist
# 2. 把 web/dist 链接到本目录供 go:embed
# 3. 交叉编译 amd64 / arm64 二进制到 ~/mita-lzcapp/content/
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
cd "$ROOT/web"
pnpm install --frozen-lockfile 2>/dev/null || pnpm install
pnpm build

# go:embed 不跟随目录符号链接，复制真实目录
rm -rf "$ROOT/tools/renderpage/web"
mkdir -p "$ROOT/tools/renderpage/web"
cp -r "$ROOT/web/dist" "$ROOT/tools/renderpage/web/dist"

cd "$ROOT/tools/renderpage"
for arch in amd64 arm64; do
  GOOS=linux GOARCH=$arch go build -trimpath -ldflags "-s -w" -o "renderpage-$arch" .
done

cp renderpage-amd64 renderpage-arm64 "$ROOT/content/"
echo "built: $(ls -la renderpage-amd64 renderpage-arm64 | awk '{print $NF, $5}')"
