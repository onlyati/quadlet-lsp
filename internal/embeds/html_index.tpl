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
          {{ range $key, $value := .Quadlets }}
            <li>
              <a class="sidebar-link" href="{{ $key }}.html">{{ $key }}</a>
            </li>
          {{ end }}
          </ul>
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
