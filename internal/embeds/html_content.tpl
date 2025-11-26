<html lang="en">
  <head>
    <title>Quadlets</title>
    <meta
      name="viewport"
      content="width=device-width,initial-scale=1,shrink-to-fit=no"
    />
    <link rel="stylesheet" href="./style.css" />
  </head>
  <body class="app">
    <header class="site-header">
      <div class="site-header-inner">
        <a href="/" class="site-logo">
          <span>Quadlet Docs</span>
        </a>
        <nav class="site-header-nav">
          <a href="https://github.com/onlyati/quadlet-lsp">GitHub</a>
        </nav>
      </div>
    </header>

    <main class="site-main">
      <div class="site-main-inner">
        <aside class="sidebar">
          {{ range $key, $value := .MenuItems }}
          {{ if gt (len $value) 0 }}
          <details>
            <summary class="sidebar-title">{{ $key }}</summary>
            <ul class="sidebar-nav">
              {{ range $value }}
              <li>
                <a class="sidebar-link" href="{{ replaceSlash . }}.html">{{ . }}</a>
              </li>
              {{ end }}
            </ul>
          </details>
          {{ end }}
          {{ end }}
        </aside>

        <section class="content">
          <div class="content-header">
            <h1>{{ .Q.Name }}</h1>
          </div>
          {{ if .Q.HeaderHTML }}
          <div class="card">
            <p>{{ .Q.HeaderHTML }}</p>
          </div>
          {{ end }}
          {{ if .Q.PartOf }}
          <div class="card">
            <div class="card-header">
              <h2 class="quadlet-section-title">Part of these units</h2>
            </div>
            <div class="flex gap-4">
            {{ range .Q.PartOf }}
            <a class="btn btn-ghost" href="{{ . }}.html">{{ . }}</a>
            {{ end }}
            </div>
          </div>
          {{ end }}
          {{ if .Q.References }}
          <div class="card">
            <div class="card-header">
              <h2 class="quadlet-section-title">Other unit references</h2>
            </div>
            <div class="flex gap-4">
            {{ range .Q.References }}
            <a class="btn btn-ghost" href="{{ . }}.html">{{ . }}</a>
            {{ end }}
            </div>
          </div>
          {{ end }}
          {{ range $key, $value := .Q.Properties }}
          <section class="quadlet-section">
            <h2 class="quadlet-section-title">{{ $key }}</h2>
            <dl class="prop-list">
              {{ range $value }}
              <div class="prop-row">
                <dt>{{ .Property }}</dt>
                <dd><code>{{ .Value }}</code></dd>
              </div>
              {{ end }}
            </dl>
          </section>
          {{ end }}

          <div class="card">
            <div class="card-header">
              <h2 class="quadlet-section-title">Source</h2>
            </div>
            <pre class="code-terminal"><code>{{ .Q.SourceFile }}</code></pre>
          </div>

          {{ if .Q.Dropins }}
          {{ range .Q.Dropins }}
          <div class="content-header"><h1>{{ .FileName }}</h1></div>

          {{ range $key, $value := .Properties }}
          <section class="quadlet-section">
            <h2 class="quadlet-section-title">{{ $key }}</h2>
            <dl class="prop-list">
              {{ range $value }}
              <div class="prop-row">
                <dt>{{ .Property }}</dt>
                <dd><code>{{ .Value }}</code></dd>
              </div>
              {{ end }}
            </dl>
          </section>
          {{ end }}

          <div class="card">
            <div class="card-header">
              <h2 class="quadlet-section-title">Source</h2>
            </div>
            <pre class="code-terminal"><code>{{ .SourceFile }}</code></pre>
          </div>
          {{ end }}
          {{ end }}
        </section>
      </div>
    </main>

    <footer class="site-footer">
      <div class="site-footer-inner">
        <span>© Quadlet LSP - docs</span>
      </div>
    </footer>
  </body>
</html>
