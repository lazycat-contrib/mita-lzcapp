package main

// renderpage: 读取 mita 的 server.json，把连接参数注入 Vite 构建的前端页面，
// 输出最终的 index.html。
//
// 用法: renderpage <server.json> <out.html>
//
// 前端产物位于 web/dist（Vite + Vue 构建），通过 go:embed 打进本二进制；
// 修改前端后执行 tools/renderpage/build.sh 重新构建。

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"embed"
)

//go:embed all:web/dist
var distFS embed.FS

type portBinding struct {
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
}

type trafficPattern struct {
	Seed      json.Number `json:"seed"`
	UnlockAll bool        `json:"unlockAll"`
}

type serverConf struct {
	PortBindings []portBinding `json:"portBindings"`
	Users        []struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	} `json:"users"`
	LoggingLevel   string         `json:"loggingLevel"`
	TrafficPattern trafficPattern `json:"trafficPattern"`
}

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "usage: renderpage <server.json> <out.html>")
		os.Exit(2)
	}
	raw, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "read config:", err)
		os.Exit(1)
	}
	var conf serverConf
	if err := json.Unmarshal(raw, &conf); err != nil {
		fmt.Fprintln(os.Stderr, "parse config:", err)
		os.Exit(1)
	}

	tcpPort, udpPort := 0, 0
	for _, pb := range conf.PortBindings {
		switch pb.Protocol {
		case "TCP":
			tcpPort = pb.Port
		case "UDP":
			udpPort = pb.Port
		}
	}
	if tcpPort == 0 {
		tcpPort = 28999
	}
	username, password := "", ""
	if len(conf.Users) > 0 {
		username = conf.Users[0].Name
		password = conf.Users[0].Password
	}

	confJSON, _ := json.Marshal(map[string]any{
		"tcpPort":  tcpPort,
		"udpPort":  udpPort,
		"username": username,
		"password": password,
		"seed":     conf.TrafficPattern.Seed.String(),
	})

	// 找到 index.html 与其引用的静态资源
	indexHTML, err := distFS.ReadFile("web/dist/index.html")
	if err != nil {
		fmt.Fprintln(os.Stderr, "embed index.html:", err)
		os.Exit(1)
	}

	outDir := os.Args[2]
	if err := os.MkdirAll(filepath.Dir(outDir), 0755); err != nil {
		fmt.Fprintln(os.Stderr, "mkdir:", err)
		os.Exit(1)
	}

	// 写出 index.html，注入配置对象
	html := strings.Replace(string(indexHTML),
		"<head>",
		fmt.Sprintf("<head><script>window.__MITA_CONF__=%s;</script>", confJSON), 1)
	if err := os.WriteFile(outDir, []byte(html), 0644); err != nil {
		fmt.Fprintln(os.Stderr, "write output:", err)
		os.Exit(1)
	}

	// 写出带 hash 的静态资源（assets/*），保持内嵌路径层级
	entries, err := distFS.ReadDir("web/dist/assets")
	if err != nil {
		fmt.Fprintln(os.Stderr, "embed assets:", err)
		os.Exit(1)
	}
	baseDir := filepath.Dir(outDir)
	if err := os.MkdirAll(filepath.Join(baseDir, "assets"), 0755); err != nil {
		fmt.Fprintln(os.Stderr, "mkdir assets:", err)
		os.Exit(1)
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		data, err := distFS.ReadFile("web/dist/assets/" + e.Name())
		if err != nil {
			fmt.Fprintln(os.Stderr, "embed asset:", err)
			os.Exit(1)
		}
		if err := os.WriteFile(filepath.Join(baseDir, "assets", e.Name()), data, 0644); err != nil {
			fmt.Fprintln(os.Stderr, "write asset:", err)
			os.Exit(1)
		}
	}

	// favicon 等根级静态文件
	for _, root := range []string{"favicon.svg"} {
		if data, err := distFS.ReadFile("web/dist/" + root); err == nil {
			_ = os.WriteFile(filepath.Join(baseDir, root), data, 0644)
		}
	}
}
