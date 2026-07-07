/*
 * Copyright 2026 Jonas Kaninda
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package storefront

// baseHTML is the shared document shell. Each page provides a "content" block.
const baseHTML = `{{define "page"}}<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Miabi Marketplace</title>
<style>
  :root { --p:#9333ea; --p-50:#faf5ff; --bg:#fff; --fg:#1e1b2e; --muted:#6b7280; --border:#e5e7eb; }
  * { box-sizing:border-box; }
  body { margin:0; font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,sans-serif; color:var(--fg); background:#f9fafb; }
  a { color:var(--p); text-decoration:none; }
  a:hover { text-decoration:underline; }
  header.site { background:var(--bg); border-bottom:1px solid var(--border); padding:16px 24px; }
  header.site .brand { font-size:18px; font-weight:700; color:var(--fg); }
  header.site .brand b { color:var(--p); }
  header.site .tag { color:var(--muted); font-size:13px; margin-left:10px; }
  main { max-width:1100px; margin:0 auto; padding:24px; }
  .badge { display:inline-block; font-size:11px; font-weight:600; padding:2px 8px; border-radius:999px; }
  .badge-official { background:#ecfdf5; color:#047857; }
  .badge-community { background:#eff6ff; color:#1d4ed8; }
  .badge-cat { background:#f3f4f6; color:var(--muted); }
  .grid { display:grid; grid-template-columns:repeat(auto-fill,minmax(260px,1fr)); gap:16px; }
  .card { background:var(--bg); border:1px solid var(--border); border-radius:12px; padding:16px; transition:box-shadow .15s,transform .15s; display:block; }
  .card:hover { box-shadow:0 6px 20px rgba(0,0,0,.08); transform:translateY(-2px); text-decoration:none; }
  .card .head { display:flex; justify-content:space-between; align-items:flex-start; margin-bottom:10px; }
  .icon { width:40px; height:40px; border-radius:8px; background:var(--p-50); color:var(--p); display:flex; align-items:center; justify-content:center; overflow:hidden; font-weight:700; }
  .icon img { width:26px; height:26px; object-fit:contain; }
  .name { font-size:15px; font-weight:600; margin:0; color:var(--fg); }
  .desc { font-size:13px; color:var(--muted); margin:4px 0 0; min-height:2.4em; }
  .meta { font-size:12px; color:var(--muted); margin-top:12px; display:flex; justify-content:space-between; }
  form.search { display:flex; gap:8px; flex-wrap:wrap; margin-bottom:20px; }
  form.search input, form.search select { padding:8px 10px; border:1px solid var(--border); border-radius:8px; font-size:14px; }
  form.search input[type=text] { flex:1; min-width:200px; }
  form.search button { background:var(--p); color:#fff; border:0; border-radius:8px; padding:8px 16px; font-weight:600; cursor:pointer; }
  .pager { display:flex; justify-content:center; align-items:center; gap:16px; margin-top:24px; color:var(--muted); font-size:14px; }
  .btn { border:1px solid var(--border); background:#fff; border-radius:8px; padding:6px 14px; color:var(--fg); }
  .empty { text-align:center; color:var(--muted); padding:48px 0; }
  .detail-head { display:flex; gap:16px; align-items:center; margin-bottom:16px; }
  .detail-head .icon { width:56px; height:56px; }
  .detail-head .icon img { width:36px; height:36px; }
  .section { background:var(--bg); border:1px solid var(--border); border-radius:12px; padding:16px 20px; margin-top:16px; }
  .section h3 { margin:0 0 10px; font-size:14px; }
  ul.bare { list-style:none; padding:0; margin:0; }
  ul.bare li { padding:4px 0; font-size:14px; border-bottom:1px solid var(--border); }
  ul.bare li:last-child { border-bottom:0; }
  code { background:#f3f4f6; padding:1px 6px; border-radius:4px; font-size:13px; }
  pre.readme { white-space:pre-wrap; background:#0f0a1e0a; padding:14px; border-radius:8px; font-size:13px; line-height:1.5; overflow:auto; }
</style>
</head>
<body>
<header class="site"><span class="brand">Mia<b>b</b>i Marketplace</span><span class="tag">official &amp; community app templates</span></header>
<main>{{template "content" .}}</main>
</body></html>{{end}}`

const homeHTML = `{{define "content"}}
<form class="search" method="get" action="/">
  <input type="text" name="q" value="{{.Q}}" placeholder="Search templates…">
  <select name="source">
    <option value=""{{if eq .Source ""}} selected{{end}}>All sources</option>
    <option value="official"{{if eq .Source "official"}} selected{{end}}>Official</option>
    <option value="community"{{if eq .Source "community"}} selected{{end}}>Community</option>
  </select>
  <select name="category">
    <option value="">All categories</option>
    {{range .Categories}}<option value="{{.Category}}"{{if eq $.Category .Category}} selected{{end}}>{{.Category}} ({{.Count}})</option>{{end}}
  </select>
  <button type="submit">Search</button>
</form>

{{if eq .Result.Total 0}}
  <div class="empty">No templates found.</div>
{{else}}
  <div class="grid">
    {{range .Result.Items}}
    <a class="card" href="/templates/{{.Name}}">
      <div class="head">
        <span class="icon">{{if isURL .Icon}}<img src="{{.Icon}}" alt="{{.DisplayName}}">{{else}}{{slice .DisplayName 0 1}}{{end}}</span>
        <span class="badge {{badge .Source}}">{{.Source}}</span>
      </div>
      <p class="name">{{.DisplayName}}</p>
      <p class="desc">{{.Description}}</p>
      <div class="meta"><span>{{provision .}}</span><span>v{{.Version}}</span></div>
    </a>
    {{end}}
  </div>
  <div class="pager">
    {{if .PrevURL}}<a class="btn" href="{{.PrevURL}}">← Prev</a>{{end}}
    <span>Page {{.Result.Page}} of {{.Result.TotalPages}} · {{.Result.Total}} templates</span>
    {{if .NextURL}}<a class="btn" href="{{.NextURL}}">Next →</a>{{end}}
  </div>
{{end}}
{{end}}`

const detailHTML = `{{define "content"}}
<p><a href="/">← All templates</a></p>
<div class="detail-head">
  <span class="icon">{{if isURL .Listing.Icon}}<img src="{{.Listing.Icon}}" alt="{{.Listing.DisplayName}}">{{else}}{{slice .Listing.DisplayName 0 1}}{{end}}</span>
  <div>
    <h1 style="margin:0">{{.Listing.DisplayName}}</h1>
    <div style="margin-top:6px">
      <span class="badge {{badge .Listing.Source}}">{{.Listing.Source}}</span>
      {{if .Listing.Category}}<span class="badge badge-cat">{{.Listing.Category}}</span>{{end}}
      <span style="color:var(--muted);font-size:13px;margin-left:6px">v{{.Listing.Version}}</span>
    </div>
  </div>
</div>
<p style="color:var(--muted)">{{.Listing.Description}}</p>
<p>
  {{if .Listing.Homepage}}<a href="{{.Listing.Homepage}}" rel="nofollow noopener">Homepage ↗</a>{{end}}
  {{if .T.Meta.SourceRepo}} &nbsp;·&nbsp; <a href="{{.T.Meta.SourceRepo}}" rel="nofollow noopener">Source ↗</a>{{end}}
</p>

<div class="section">
  <h3>What gets created — {{.Provision}}</h3>
  <ul class="bare">
    {{range .M.Applications}}<li>🧊 App <b>{{.Name}}</b> — <code>{{.Image}}{{if .Tag}}:{{.Tag}}{{end}}</code></li>{{end}}
    {{range .M.Databases}}<li>🗄️ Database <b>{{.Name}}</b> — {{.Engine}}{{if .Version}} {{.Version}}{{end}} ({{.Placement}})</li>{{end}}
    {{range .M.Volumes}}<li>💾 Volume <b>{{.Name}}</b></li>{{end}}
  </ul>
</div>

{{if .M.Inputs}}
<div class="section">
  <h3>Install inputs</h3>
  <ul class="bare">
    {{range .M.Inputs}}<li><code>{{.Key}}</code>{{if .Label}} — {{.Label}}{{end}}{{if .Required}} <span style="color:#b91c1c">required</span>{{end}}</li>{{end}}
  </ul>
</div>
{{end}}

<div class="section">
  <h3>Versions</h3>
  <ul class="bare">
    {{range .T.Versions}}<li>v{{.Version}} — <code>{{.Digest}}</code></li>{{end}}
  </ul>
</div>

<div class="section">
  <h3>Install</h3>
  <p style="font-size:14px;color:var(--muted)">In Miabi, open <b>Marketplace → {{.Listing.Source}}</b>, find <b>{{.Listing.DisplayName}}</b>, and click Install. Requires <code>MIABI_MARKETPLACE_URL</code> to point at this registry.</p>
</div>

{{if .T.Readme}}
<div class="section">
  <h3>README</h3>
  <pre class="readme">{{.T.Readme}}</pre>
</div>
{{end}}
{{end}}`

// notFoundHTML is the 404 page content, rendered through the shared base layout.
const notFoundHTML = `{{define "content"}}
<div class="empty">
  <h1 style="font-size:20px;margin:0 0 8px">404 — template not found</h1>
  <p><a href="/">← Browse all templates</a></p>
</div>
{{end}}`
