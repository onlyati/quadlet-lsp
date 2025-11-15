/* ---------- Color tokens ---------- */

:root {
  /* Podman brand */
  --color-primary: #892ca0;
  --color-accent: #c42482;

  /* Light theme backgrounds */
  --color-bg: #f8f5fa;
  --color-bg-panel: #ffffff;
  --color-bg-hover: #f1e9f7;

  /* Text */
  --color-text: #2d2235;
  --color-text-muted: #6d5a77;

  /* Borders */
  --color-border: #d9cbe3;
  --color-border-light: #eae0f0;

  /* Semantic */
  --color-success: #2e9e6f;
  --color-warning: #dba42e;
  --color-error: #d64566;
  --color-info: #5b8ef4;

  /* Shadows */
  --shadow-soft: 0 4px 14px rgba(0, 0, 0, 0.08);
  --shadow-strong: 0 8px 20px rgba(0, 0, 0, 0.15);

  /* Radii */
  --radius-sm: 4px;
  --radius-md: 8px;
  --radius-lg: 12px;

  /* Transitions */
  --transition-fast: 120ms ease-in-out;
  --transition-slow: 250ms ease-in-out;

  /* Layout */
  --layout-max-width: 1200px;
  --layout-sidebar-width: 260px;
  --layout-header-height: 56px;
}

/* ---------- Reset / base ---------- */

*,
*::before,
*::after {
  box-sizing: border-box;
}

html,
body {
  margin: 0;
  padding: 0;
}

body {
  min-height: 100vh;
  background-color: var(--color-bg);
  color: var(--color-text);
  font-family:
    "Inter",
    system-ui,
    -apple-system,
    BlinkMacSystemFont,
    "Segoe UI",
    sans-serif;
  line-height: 1.6;
  text-rendering: optimizeLegibility;
}

/* Remove default margins on headings and paragraphs, we’ll control them */
h1,
h2,
h3,
h4,
h5,
h6,
p,
ul,
ol,
pre,
figure {
  margin: 0;
}

img {
  max-width: 100%;
  display: block;
}

/* ---------- Layout shell ---------- */

.app {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

/* Header */

.site-header {
  height: var(--layout-header-height);
  display: flex;
  align-items: center;
  justify-content: center;
  border-bottom: 1px solid var(--color-border);
  background: linear-gradient(90deg, #f8f5fa 0%, #f3e6fb 40%, #f1d9f7 100%);
  padding: 0 1.5rem;
  position: sticky;
  top: 0;
  z-index: 20;
}

.site-header-inner {
  width: 100%;
  max-width: var(--layout-max-width);
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1.5rem;
}

.site-logo {
  display: flex;
  align-items: center;
  gap: 0.65rem;
  font-weight: 600;
  font-size: 1.1rem;
  color: var(--color-text);
  text-decoration: none;
}

.site-logo-mark {
  width: 28px;
  height: 28px;
  border-radius: 999px;
  background: conic-gradient(
    from 140deg,
    var(--color-primary),
    var(--color-accent),
    #e89cdc,
    var(--color-primary)
  );
  box-shadow: 0 0 18px rgba(137, 44, 160, 0.35);
}

.site-header-nav {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.site-header-nav a {
  font-size: 0.9rem;
  color: var(--color-text-muted);
  text-decoration: none;
  padding: 0.25rem 0.5rem;
  border-radius: var(--radius-sm);
  transition:
    background-color var(--transition-fast),
    color var(--transition-fast);
}

.site-header-nav a:hover {
  background-color: var(--color-bg-hover);
  color: var(--color-text);
}

/* Main layout: sidebar + content */

.site-main {
  flex: 1;
  display: flex;
  justify-content: center;
  padding: 1.25rem 1.5rem 2rem;
}

.site-main-inner {
  width: 100%;
  max-width: var(--layout-max-width);
  display: grid;
  grid-template-columns: minmax(0, var(--layout-sidebar-width)) minmax(0, 1fr);
  gap: 1.5rem;
}

/* Sidebar */

.sidebar {
  background-color: var(--color-bg-panel);
  border-radius: var(--radius-lg);
  border: 1px solid var(--color-border);
  padding: 1rem 1rem 1.25rem;
  box-shadow: var(--shadow-soft);
  align-self: flex-start;
  position: sticky;
  top: calc(var(--layout-header-height) + 1rem);
  max-height: calc(100vh - var(--layout-header-height) - 2rem);
  overflow: auto;
}

.sidebar-title {
  font-size: 0.85rem;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: var(--color-text-muted);
  margin-bottom: 0.75rem;
}

.sidebar-nav {
  list-style: none;
  padding: 0;
  margin: 0;
}

.sidebar-nav li + li {
  margin-top: 0.25rem;
}

.sidebar-link {
  display: block;
  padding: 0.4rem 0.6rem;
  border-radius: var(--radius-md);
  text-decoration: none;
  font-size: 0.9rem;
  color: var(--color-text-muted);
  transition:
    background-color var(--transition-fast),
    color var(--transition-fast);
}

.sidebar-link:hover {
  background-color: var(--color-bg-hover);
  color: var(--color-text);
}

.sidebar-link-active {
  background: linear-gradient(90deg, var(--color-primary), var(--color-accent));
  color: #ffffff;
}

/* Content */

.content {
  background-color: transparent;
}

/* Footer */

.site-footer {
  padding: 1rem 1.5rem 1.5rem;
  border-top: 1px solid var(--color-border);
  color: var(--color-text-muted);
  font-size: 0.85rem;
}

.site-footer-inner {
  max-width: var(--layout-max-width);
  margin: 0 auto;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
}

/* ---------- Typography ---------- */

.content-header {
  margin-bottom: 1.5rem;
}

h1 {
  font-size: 2rem;
  font-weight: 650;
  letter-spacing: 0.01em;
  margin-bottom: 0.35rem;
}

h2 {
  font-size: 1.4rem;
  font-weight: 600;
  margin-top: 1.75rem;
  margin-bottom: 0.5rem;
}

h3 {
  font-size: 1.15rem;
  font-weight: 550;
  margin-top: 1.3rem;
  margin-bottom: 0.5rem;
}

h4 {
  font-size: 1rem;
  font-weight: 550;
  margin-top: 1.15rem;
  margin-bottom: 0.5rem;
}

.content-header p,
p {
  color: var(--color-text-muted);
}

p + p {
  margin-top: 0.5rem;
}

p + h2,
ul + h2,
pre + h2 {
  margin-top: 1.75rem;
}

a {
  color: var(--color-accent);
  text-decoration: none;
  transition:
    color var(--transition-fast),
    text-shadow var(--transition-fast);
}

a:hover {
  color: #e053a7;
  text-shadow: 0 0 4px rgba(196, 36, 130, 0.3);
}

strong {
  font-weight: 600;
}

/* ---------- Lists ---------- */

ul,
ol {
  padding-left: 1.25rem;
  margin-top: 0.75rem;
  margin-bottom: 0.75rem;
}

li + li {
  margin-top: 0.25rem;
}

/* ---------- Cards / panels ---------- */

.card {
  background-color: var(--color-bg-panel);
  border-radius: var(--radius-lg);
  border: 1px solid var(--color-border);
  padding: 1rem 1.25rem;
  box-shadow: var(--shadow-soft);
  margin-bottom: 1rem;
}

.card + .card {
  margin-top: 0.75rem;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 0.75rem;
  margin-bottom: 0.75rem;
}

.card-title {
  font-size: 1rem;
  font-weight: 550;
}

/* ---------- Buttons ---------- */

.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 0.4rem;
  padding: 0.45rem 0.9rem;
  border-radius: var(--radius-md);
  font-size: 0.9rem;
  font-weight: 500;
  border: 1px solid transparent;
  cursor: pointer;
  text-decoration: none;
  transition:
    background-color var(--transition-fast),
    border-color var(--transition-fast),
    color var(--transition-fast),
    transform var(--transition-fast),
    box-shadow var(--transition-fast);
}

.btn:active {
  transform: translateY(1px);
}

.btn-primary {
  background: linear-gradient(
    135deg,
    var(--color-primary),
    var(--color-accent)
  );
  color: #ffffff;
  border-color: var(--color-accent);
  box-shadow: 0 4px 12px rgba(137, 44, 160, 0.35);
}

.btn-primary:hover {
  background: linear-gradient(135deg, #9e38c0, #e46ab5);
}

.btn-secondary {
  background-color: var(--color-bg-panel);
  color: var(--color-text);
  border-color: var(--color-border-light);
}

.btn-secondary:hover {
  background-color: var(--color-bg-hover);
}

.btn-ghost {
  background: transparent;
  color: var(--color-text-muted);
  border-color: transparent;
}

.btn-ghost:hover {
  background-color: var(--color-bg-hover);
  color: var(--color-text);
}

/* ---------- Code blocks ---------- */

code,
pre {
  font-family:
    "JetBrains Mono", "Fira Code", Menlo, Monaco, Consolas, "Liberation Mono",
    "Courier New", monospace;
  font-size: 0.85rem;
}

code {
  background-color: #f3ecf8;
  padding: 0.1rem 0.3rem;
  border-radius: var(--radius-sm);
  border: 1px solid #e3d5ec;
}

pre {
  background: #f3ecf8;
  border-radius: var(--radius-lg);
  border: 1px solid var(--color-border);
  padding: 0.85rem 1rem;
  margin: 0.85rem 0;
  overflow-x: auto;
  box-shadow: var(--shadow-soft);
}

pre code {
  background: transparent;
  border: none;
  padding: 0;
}

/* Optional “terminal” style */
.code-terminal {
  position: relative;
  padding-top: 1.6rem;
}

.code-terminal::before {
  content: "Quadlet file";
  position: absolute;
  top: 0.35rem;
  left: 1rem;
  font-size: 0.75rem;
  color: var(--color-text-muted);
}

/* ---------- Tables ---------- */

table {
  width: 100%;
  border-collapse: collapse;
  margin: 0.75rem 0;
  font-size: 0.9rem;
}

th,
td {
  border: 1px solid var(--color-border);
  padding: 0.4rem 0.6rem;
}

th {
  text-align: left;
  background-color: var(--color-bg-hover);
  font-weight: 550;
}

tbody tr:nth-child(even) {
  background-color: #fdfbff;
}

/* ---------- Alerts / callouts ---------- */

.alert {
  border-radius: var(--radius-lg);
  border: 1px solid var(--color-border);
  padding: 0.6rem 0.8rem;
  display: flex;
  gap: 0.6rem;
  align-items: flex-start;
  margin: 0.75rem 0;
  background-color: var(--color-bg-panel);
  font-size: 0.9rem;
}

.alert-icon {
  margin-top: 0.1rem;
  font-size: 1rem;
}

.alert-title {
  font-weight: 550;
  margin-bottom: 0.1rem;
}

.alert-content p + p {
  margin-top: 0.3rem;
}

.alert-info {
  border-color: rgba(91, 142, 244, 0.7);
  background-color: #eef3ff;
  box-shadow: 0 0 18px rgba(91, 142, 244, 0.15);
}

.alert-success {
  border-color: rgba(46, 158, 111, 0.7);
  background-color: #e8f5ef;
}

.alert-warning {
  border-color: rgba(219, 164, 46, 0.75);
  background-color: #fff6e3;
}

.alert-error {
  border-color: rgba(214, 69, 102, 0.8);
  background-color: #ffe9ef;
}

/* ---------- Badges / labels ---------- */

.badge {
  display: inline-flex;
  align-items: center;
  padding: 0.1rem 0.45rem;
  border-radius: 999px;
  font-size: 0.7rem;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  border: 1px solid var(--color-border-light);
  color: var(--color-text-muted);
  background-color: #fdfbff;
}

.badge-primary {
  border-color: rgba(137, 44, 160, 0.6);
  color: #7a2b9b;
  background-color: #f3e4fb;
}

/* ---------- Forms (basic) ---------- */

input[type="text"],
input[type="search"],
input[type="email"],
input[type="password"],
textarea,
select {
  width: 100%;
  padding: 0.4rem 0.55rem;
  border-radius: var(--radius-md);
  border: 1px solid var(--color-border-light);
  background-color: var(--color-bg-panel);
  color: var(--color-text);
  font: inherit;
  transition:
    border-color var(--transition-fast),
    box-shadow var(--transition-fast),
    background-color var(--transition-fast);
}

input:focus,
textarea:focus,
select:focus {
  outline: none;
  border-color: var(--color-accent);
  box-shadow: 0 0 0 2px rgba(196, 36, 130, 0.25);
}

/* ---------- Utilities ---------- */

.mt-1 {
  margin-top: 0.25rem;
}
.mt-2 {
  margin-top: 0.5rem;
}
.mt-3 {
  margin-top: 0.75rem;
}
.mt-4 {
  margin-top: 1rem;
}

.mb-1 {
  margin-bottom: 0.25rem;
}
.mb-2 {
  margin-bottom: 0.5rem;
}
.mb-3 {
  margin-bottom: 0.75rem;
}
.mb-4 {
  margin-bottom: 1rem;
}

.flex {
  display: flex;
}

.flex-center {
  display: flex;
  align-items: center;
  justify-content: center;
}

.gap-1 {
  gap: 0.25rem;
}
.gap-2 {
  gap: 0.5rem;
}
.gap-3 {
  gap: 0.75rem;
}
.gap-4 {
  gap: 1rem;
}

/* ---------- Responsive ---------- */

@media (max-width: 920px) {
  .site-main-inner {
    grid-template-columns: minmax(0, 1fr);
  }

  .sidebar {
    position: static;
    max-height: none;
    order: -1;
  }
}

@media (max-width: 640px) {
  .site-header-inner {
    flex-direction: row;
  }

  .site-header-nav {
    display: none; /* Keep it simple; later you can add JS toggle for mobile nav */
  }

  .site-main {
    padding: 1rem;
  }

  .site-footer-inner {
    flex-direction: column;
    align-items: flex-start;
  }
}

/* ---------- Quadlet sections (Section / Property / Value) ---------- */

.quadlet-section {
  background-color: var(--color-bg-panel);
  border-radius: var(--radius-lg);
  border: 1px solid var(--color-border);
  padding: 1rem 1.25rem;
  box-shadow: var(--shadow-soft);
  margin-bottom: 1rem;
}

.quadlet-section-title {
  font-size: 1.2rem;
  font-weight: 600;
  color: var(--color-primary);
  margin-bottom: 0.75rem;
}

/* Definition list format */
.prop-list {
  margin: 0;
  padding: 0;
}

.prop-row {
  display: grid;
  grid-template-columns: 160px minmax(0, 1fr);
  gap: 0.5rem 1rem;
  padding: 0.5rem 0;
  border-bottom: 1px solid var(--color-border-light);
}

.prop-row:last-child {
  border-bottom: none;
}

.prop-row dt {
  font-weight: 600;
  color: var(--color-text);
}

.prop-row dd {
  margin: 0;
  color: var(--color-text-muted);
}

.prop-row code {
  background-color: #f3ecf8;
  padding: 0.15rem 0.35rem;
  border-radius: var(--radius-sm);
  border: 1px solid #e3d5ec;
  font-size: 0.85rem;
}

@media (max-width: 600px) {
  .prop-row {
    grid-template-columns: 1fr;
  }

  .prop-row dt {
    font-size: 0.75rem;
    text-transform: uppercase;
    letter-spacing: 0.06em;
    color: var(--color-text-muted);
  }
}
