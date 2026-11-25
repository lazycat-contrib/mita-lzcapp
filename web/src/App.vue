<script setup>
import { ref, onMounted } from 'vue'

// 构建时由 Vite 注入：连接参数来自 mita server.json（见 vite.config.js / renderpage）
const conf = window.__MITA_CONF__ || {}

const lang = ref(localStorage.getItem('mita-lang') || 'zh')
const theme = ref(localStorage.getItem('mita-theme') || 'auto')
const copied = ref('')

const T = {
  zh: {
    title: 'mita 服务端正在运行',
    subtitle: '这是 mieru 代理体系的服务端。用 mieru / Clash 系客户端按下面的参数连接即可。',
    conn: '客户端连接参数',
    addr: '地址',
    addrval: '本机内网 IP（局域网直连）；外网用公网 IP / 域名并在路由器上转发该端口',
    port: '端口',
    portnote: 'TCP 与 UDP 均已开启',
    proto: '协议',
    user: '用户名',
    pass: '密码',
    seed: '流量特征种子',
    warn: '种子必须一致：客户端 profile 里的 trafficPattern.seed 必须填写同一个数字，否则无法建立会话。',
    clients: '推荐客户端',
    clDesktop: '桌面：mieru 官方客户端、Clash Verge Rev、Mihomo Party',
    clAndroid: '安卓：ClashMetaForAndroid、ClashMi、husi + mieru plugin',
    clIos: 'iOS：Karing、ClashMi',
    maint: '常用维护命令（应用终端）',
    m1: '查看服务状态',
    m2: '查看用户用量',
    m3: '编辑配置模板',
    m4: '应用修改',
    copied: '已复制',
    copy: '复制',
    auto: '自动', light: '亮色', dark: '暗色',
  },
  en: {
    title: 'mita server is running',
    subtitle: 'This is the server side of the mieru proxy system. Connect from any mieru / Clash client using the parameters below.',
    conn: 'Client connection parameters',
    addr: 'Address',
    addrval: 'LazyCat LAN IP (direct on local network); from the internet use your public IP / domain and forward the port on your router',
    port: 'Port',
    portnote: 'both TCP and UDP enabled',
    proto: 'Protocol',
    user: 'Username',
    pass: 'Password',
    seed: 'Traffic pattern seed',
    warn: 'Seed must match: the trafficPattern.seed in your client profile must be the same number, otherwise sessions cannot be established.',
    clients: 'Recommended clients',
    clDesktop: 'Desktop: mieru client, Clash Verge Rev, Mihomo Party',
    clAndroid: 'Android: ClashMetaForAndroid, ClashMi, husi + mieru plugin',
    clIos: 'iOS: Karing, ClashMi',
    maint: 'Maintenance commands (app terminal)',
    m1: 'Show service status',
    m2: 'Show user usage',
    m3: 'Edit config template',
    m4: 'Apply changes',
    copied: 'Copied',
    copy: 'Copy',
    auto: 'Auto', light: 'Light', dark: 'Dark',
  },
}

const tcpPort = conf.tcpPort ?? 28999
const username = conf.username ?? ''
const password = conf.password ?? ''
const seed = conf.seed ?? ''

function setLang(l) {
  lang.value = l
  localStorage.setItem('mita-lang', l)
  document.documentElement.lang = l === 'zh' ? 'zh-CN' : 'en'
}
function setTheme(t) {
  theme.value = t
  localStorage.setItem('mita-theme', t)
  applyTheme(t)
}
function applyTheme(t) {
  const eff = t === 'auto'
    ? (window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light')
    : t
  document.documentElement.dataset.theme = eff
}
async function copy(label, text) {
  try {
    await navigator.clipboard.writeText(text)
  } catch {
    const ta = document.createElement('textarea')
    ta.value = text; document.body.appendChild(ta); ta.select()
    document.execCommand('copy'); ta.remove()
  }
  copied.value = label
  setTimeout(() => { if (copied.value === label) copied.value = '' }, 1500)
}
onMounted(() => {
  applyTheme(theme.value)
  document.documentElement.lang = lang.value === 'zh' ? 'zh-CN' : 'en'
  window.matchMedia('(prefers-color-scheme: dark)')
    .addEventListener('change', () => applyTheme(theme.value))
})
</script>

<template>
  <div class="toolbar">
    <div class="seg">
      <button :class="{ on: lang === 'zh' }" @click="setLang('zh')">中文</button>
      <button :class="{ on: lang === 'en' }" @click="setLang('en')">English</button>
    </div>
    <div class="seg">
      <button :class="{ on: theme === 'auto' }" @click="setTheme('auto')">{{ T[lang].auto }}</button>
      <button :class="{ on: theme === 'light' }" @click="setTheme('light')">{{ T[lang].light }}</button>
      <button :class="{ on: theme === 'dark' }" @click="setTheme('dark')">{{ T[lang].dark }}</button>
    </div>
  </div>

  <div class="wrap">
    <div class="card">
      <h1><span class="dot"></span>{{ T[lang].title }}</h1>
      <p class="sub">{{ T[lang].subtitle }}</p>
    </div>

    <div class="card">
      <h2>{{ T[lang].conn }}</h2>
      <table>
        <tr>
          <td class="k">{{ T[lang].addr }}</td>
          <td><span class="val soft">{{ T[lang].addrval }}</span></td>
        </tr>
        <tr>
          <td class="k">{{ T[lang].port }}</td>
          <td>
            <span class="val big">{{ tcpPort }}</span>
            <span class="note">{{ T[lang].portnote }}</span>
          </td>
        </tr>
        <tr>
          <td class="k">{{ T[lang].proto }}</td>
          <td><span class="val">mieru</span></td>
        </tr>
        <tr>
          <td class="k">{{ T[lang].user }}</td>
          <td>
            <span class="val big">{{ username }}</span>
            <button class="cp" @click="copy('user', username)">{{ copied === 'user' ? T[lang].copied : T[lang].copy }}</button>
          </td>
        </tr>
        <tr>
          <td class="k">{{ T[lang].pass }}</td>
          <td>
            <span class="val big">{{ password }}</span>
            <button class="cp" @click="copy('pass', password)">{{ copied === 'pass' ? T[lang].copied : T[lang].copy }}</button>
          </td>
        </tr>
        <tr>
          <td class="k">{{ T[lang].seed }}</td>
          <td>
            <span class="val big">{{ seed }}</span>
            <button class="cp" @click="copy('seed', seed)">{{ copied === 'seed' ? T[lang].copied : T[lang].copy }}</button>
          </td>
        </tr>
      </table>
      <div class="warn">⚠️ {{ T[lang].warn }}</div>
    </div>

    <div class="card">
      <h2>{{ T[lang].clients }}</h2>
      <p class="cli">{{ T[lang].clDesktop }}</p>
      <p class="cli">{{ T[lang].clAndroid }}</p>
      <p class="cli">{{ T[lang].clIos }}</p>
    </div>

    <div class="card">
      <h2>{{ T[lang].maint }}</h2>
      <p class="cli"><b>mita status</b> — {{ T[lang].m1 }}</p>
      <p class="cli"><b>mita get users</b> — {{ T[lang].m2 }}</p>
      <p class="cli"><b>vi /lzcapp/var/mita/server.json</b> — {{ T[lang].m3 }}</p>
      <p class="cli"><b>mita apply config /lzcapp/var/mita/server.json &amp;&amp; mita reload</b> — {{ T[lang].m4 }}</p>
    </div>

    <p class="foot">mieru · github.com/enfein/mieru · GPL-3.0</p>
  </div>
</template>

<style>
:root {
  --bg: #f6f7f9; --card: #ffffff; --text: #1a1a1a; --muted: #6b7280;
  --accent: #4a90d9; --border: #e3e6ea; --chip: #f2f4f7;
  --warnbg: #fff7e6; --warnline: #e6a23c;
}
html[data-theme='light'] { --bg:#f6f7f9; --card:#fff; --text:#1a1a1a; --muted:#6b7280; --accent:#4a90d9; --border:#e3e6ea; --chip:#f2f4f7; --warnbg:#fff7e6; --warnline:#e6a23c; }
html[data-theme='dark']  { --bg:#14161a; --card:#1d2026; --text:#e8eaed; --muted:#9aa3af; --accent:#6aa8e8; --border:#2c313a; --chip:#262b33; --warnbg:#2e2a1d; --warnline:#d19a3a; }

* { box-sizing: border-box; }
body {
  font-family: -apple-system, 'PingFang SC', 'Microsoft YaHei', system-ui, sans-serif;
  background: var(--bg); color: var(--text); margin: 0; padding: 32px 16px 48px;
  line-height: 1.7; transition: background .3s ease, color .3s ease;
}
.wrap { max-width: 660px; margin: 0 auto; }

@keyframes fadeUp { from { opacity: 0; transform: translateY(12px); } to { opacity: 1; transform: none; } }
.card { animation: fadeUp .45s ease both; }
.card:nth-of-type(2) { animation-delay: .08s; }
.card:nth-of-type(3) { animation-delay: .16s; }
.card:nth-of-type(4) { animation-delay: .24s; }

.card {
  background: var(--card); border: 1px solid var(--border); border-radius: 12px;
  padding: 20px 24px; margin-bottom: 16px;
}
h1 { font-size: 1.35em; margin: 0 0 4px; display: flex; align-items: center; gap: 10px; }
h1 .dot { width: 10px; height: 10px; border-radius: 50%; background: #34c759; box-shadow: 0 0 0 3px rgba(52,199,89,.18); animation: pulse 2s ease infinite; }
@keyframes pulse { 50% { box-shadow: 0 0 0 6px rgba(52,199,89,.06); } }
.sub { color: var(--muted); font-size: .92em; margin: 0; }
h2 { font-size: 1.02em; margin: 0 0 12px; color: var(--accent); letter-spacing: .02em; }

table { border-collapse: collapse; width: 100%; }
td { padding: 10px 8px; border-bottom: 1px solid var(--border); vertical-align: top; }
tr:last-child td { border-bottom: none; }
td.k { color: var(--muted); white-space: nowrap; width: 8.5em; }

.val {
  font-family: ui-monospace, Menlo, Consolas, monospace;
  background: var(--chip); padding: 3px 8px; border-radius: 6px;
  word-break: break-all; user-select: all; font-size: .95em;
  transition: background .2s;
}
.val.big { font-size: 1.15em; font-weight: 600; }
.val.soft { font-family: inherit; }
.note { color: var(--muted); font-size: .88em; margin-left: 6px; }

.cp {
  border: 1px solid var(--border); background: transparent; color: var(--accent);
  border-radius: 6px; padding: 2px 10px; margin-left: 8px; cursor: pointer;
  font-size: .82em; font-family: inherit; transition: all .2s;
}
.cp:hover { background: var(--chip); }

.warn { background: var(--warnbg); border-left: 4px solid var(--warnline); padding: 10px 14px; border-radius: 8px; margin-top: 14px; font-size: .92em; }

.toolbar { display: flex; justify-content: space-between; align-items: center; max-width: 660px; margin: 0 auto 16px; }
.seg { display: flex; border: 1px solid var(--border); border-radius: 8px; overflow: hidden; background: var(--card); }
.seg button {
  border: none; background: transparent; color: var(--muted); padding: 5px 12px;
  cursor: pointer; font-size: .85em; transition: all .2s; font-family: inherit;
}
.seg button.on { background: var(--accent); color: #fff; }
.foot { color: var(--muted); font-size: .85em; text-align: center; margin-top: 8px; }
.cli { font-family: ui-monospace, Menlo, Consolas, monospace; font-size: .88em; color: var(--muted); margin: 4px 0; }
.cli b { color: var(--text); font-weight: 500; }
</style>
