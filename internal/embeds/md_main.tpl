# Quadlets

{{ range $key, $value := .Quadlets }}
- [{{ $key }}](#{{ $key }})
{{- end }}

{{ range $key, $value := .Quadlets }}
## {{ $key }}

{{- range $value.Header }}
{{ . }}
{{- end }}

{{ if .References }}
### References
{{ range .References }}
- [{{ . }}](#{{ . }})
{{- end}}
{{- end}}

{{ if .Properties }}
### File

| Section | Property | Value |
| ------- | -------- | ----- |
{{- range $key, $value := .Properties }}
{{- range $value }}
| {{ $key }} | {{ .Property }} | `{{ .Value }}` |
{{- end }}

{{- end }}
{{- end}}

{{ if .Dropins }}
### Dropins

| Location | Section | Property | Value |
| -------- | ------- | -------- | ----- |
{{- range $dropins := .Dropins }}
{{- range $key, $value := .Properties }}
{{- range $value }}
| {{ $dropins.Directory }}/{{ $dropins.FileName }} | {{ $key }} | {{ .Property }} | `{{ .Value }}` |
{{- end }}
{{- end}}
{{- end}}
{{- end }}

{{ end }}
