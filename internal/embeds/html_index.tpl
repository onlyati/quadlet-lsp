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
            <h1>Select Quadlet</h1>
          </div>
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
