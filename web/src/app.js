const app = document.querySelector('#app');
const storedTheme = localStorage.getItem('wirelab-theme') || 'light';
document.documentElement.dataset.theme = storedTheme;

const esc = (value) => String(value ?? '').replace(/[&<>"']/g, (char) => ({'&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;'}[char]));

async function api(path) {
  const response = await fetch(path);
  if (!response.ok) throw new Error(`请求失败 ${response.status}`);
  return response.json();
}

function listMarkup(items) {
  return `<div class="article-grid">${items.map((article) => `<a class="article-card" href="#/article/${esc(article.id)}">
    <div class="card-top"><span class="eyebrow">${esc(article.category)}</span><span class="arrow">↗</span></div>
    <h3>${esc(article.title)}</h3><p>${esc(article.summary)}</p><span class="card-link">查看材料与步骤</span>
  </a>`).join('')}</div>`;
}

function shell(content, active = '') {
  document.querySelectorAll('.primary-nav a').forEach((link) => link.classList.toggle('active', link.getAttribute('href') === `#/${active}`));
  app.innerHTML = content;
  app.focus();
}

function dataError(message) {
  return `<section class="state state-error"><span class="state-code">!</span><h2>数据格式错误</h2><p>${esc(message)}</p><a class="button" href="#/">返回学习首页</a></section>`;
}

async function home(query = '') {
  shell(`<section class="hero">
    <div class="hero-copy"><p class="kicker">家庭 · 小办公室 · 可复查的实践</p><h1>把每一根线，<br><em>接成可靠的网络。</em></h1><p class="hero-lede">从选线、压接到弱电箱整理，把网络布线拆成可以照着做、测得出来的步骤。</p><form class="search-form" id="search-form"><label class="search-box"><span>⌕</span><input name="q" value="${esc(query)}" placeholder="搜索：网线、水晶头、交换机…" autocomplete="off"><button type="submit">搜索</button></label></form><div class="hero-meta"><span>04 个主题</span><span>·</span><span>材料 / 步骤 / 排错</span></div></div>
    <div class="hero-image"><img src="/assets/ethernet-switch.jpg" alt="交换机与彩色网线的工作台照片"><span class="image-caption">一条链路的起点，是清楚的线序</span></div>
  </section><section class="section-head"><div><p class="kicker">资料索引 / ${query ? `搜索“${esc(query)}”` : '从这里开始'}</p><h2 id="article-heading">布线笔记</h2></div><span class="section-count" id="article-count"></span></section><div id="article-list" class="loading">正在读取本地资料…</div>`,'');
  document.querySelector('#search-form').addEventListener('submit', (event) => { event.preventDefault(); const value = new FormData(event.currentTarget).get('q').trim(); location.hash = value ? `#/search/${encodeURIComponent(value)}` : '#/'; });
  try {
    const items = await api(`/api/articles${query ? `?q=${encodeURIComponent(query)}` : ''}`);
    if (!Array.isArray(items)) throw new Error('文章列表必须是数组');
    document.querySelector('#article-count').textContent = `${items.length} 篇资料`;
    document.querySelector('#article-list').innerHTML = items.length ? listMarkup(items) : `<section class="state"><span class="state-code">0</span><h2>没有匹配文章</h2><p>换一个关键词，或从网线、交换机开始浏览。</p><a class="button" href="#/">查看全部文章</a></section>`;
  } catch (error) { document.querySelector('#article-list').innerHTML = dataError(error.message); }
}

async function detail(id) {
  shell('<div class="loading page-loading">正在读取文章…</div>');
  try {
    const article = await api(`/api/articles/${encodeURIComponent(id)}`);
    shell(`<article class="detail"><a class="back-link" href="#/">← 返回文章列表</a><header class="detail-head"><span class="eyebrow">${esc(article.category)}</span><h1>${esc(article.title)}</h1><p>${esc(article.summary)}</p><span class="audience">适合：${esc(article.audience)}</span></header><div class="detail-layout"><aside><h2>材料清单</h2><ul>${article.materials.map((item) => `<li>${esc(item)}</li>`).join('')}</ul></aside><div class="detail-body"><section><div class="section-label"><span>01</span><h2>动手步骤</h2></div><ol>${article.steps.map((item) => `<li>${esc(item)}</li>`).join('')}</ol></section><section class="troubleshoot"><div class="section-label"><span>02</span><h2>排错清单</h2></div>${article.troubleshooting.map((item) => `<div class="trouble"><h3>${esc(item.symptom)}</h3><p>${esc(item.solution)}</p></div>`).join('')}</section></div></div></article>`,'');
  } catch (error) { shell(dataError(error.message)); }
}

async function about() {
  shell(`<section class="page-intro"><p class="kicker">工具台 / 从准备开始</p><h1>工具放对位置，<br><em>排错就有了秩序。</em></h1><p>一份不追求堆满工具箱的清单：每件工具都对应一次可验证的动作。</p></section><section class="tools-layout"><div class="tool-illustration"><div class="tool-rack"><span>RJ45</span><span>CAT6</span><span>TEST</span></div><div class="tool-cable cable-one"></div><div class="tool-cable cable-two"></div><div class="tool-cable cable-three"></div></div><div id="tools" class="tool-list loading">正在读取工具清单…</div></section>`,'about');
  try {
    const tools = await api('/api/tools');
    document.querySelector('#tools').innerHTML = tools.map((tool, index) => `<div class="tool-row"><span class="tool-number">0${index + 1}</span><div><h2>${esc(tool.name)}</h2><p>${esc(tool.use)}</p></div><span class="tool-check">✓</span></div>`).join('');
  } catch (error) { document.querySelector('#tools').innerHTML = dataError(error.message); }
}

function settings() {
  const current = document.documentElement.dataset.theme;
  shell(`<section class="page-intro compact"><p class="kicker">偏好设置 / 让阅读舒服一点</p><h1>设置</h1><p>主题只保存在当前浏览器，不会发送到服务器。</p></section><section class="settings-panel"><div><span class="eyebrow">界面主题</span><h2>选择阅读环境</h2><p>深色主题适合弱电箱旁的低光环境。</p></div><div class="theme-switch" role="group" aria-label="界面主题"><button data-theme-value="light" class="${current === 'light' ? 'selected' : ''}">浅色</button><button data-theme-value="dark" class="${current === 'dark' ? 'selected' : ''}">深色</button></div></section>`,'settings');
  document.querySelectorAll('[data-theme-value]').forEach((button) => button.addEventListener('click', () => { const theme = button.dataset.themeValue; document.documentElement.dataset.theme = theme; localStorage.setItem('wirelab-theme', theme); document.querySelectorAll('[data-theme-value]').forEach((item) => item.classList.toggle('selected', item === button)); }));
}

function route() {
  const parts = location.hash.replace(/^#\/?/, '').split('/');
  if (parts[0] === 'article' && parts[1]) return detail(parts[1]);
  if (parts[0] === 'search') return home(decodeURIComponent(parts.slice(1).join('/')));
  if (parts[0] === 'about') return about();
  if (parts[0] === 'settings') return settings();
  return home();
}

window.addEventListener('hashchange', route);
route();
