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
          <div class="sidebar-title">Quadlets</div>
          <ul class="sidebar-nav">
          {{ range $key, $value := .All.Quadlets }}
            <li>
              <a class="sidebar-link" href="{{ $key }}.html">{{ $key }}</a>
            </li>
          {{ end }}
          </ul>
        </aside>

        <section class="content">
          <div class="content-header">
            <h1>{{ .Q.Name }}</h1>
          </div>
          {{ if .Q.Header }}
          <div class="card">
            <p>
              {{ range .Q.Header }}
              {{ . }}
              {{ end }}
            </p>
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
          <div class="content-header"><h1>{{ .Directory }}/{{ .FileName }}</h1></div>

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
