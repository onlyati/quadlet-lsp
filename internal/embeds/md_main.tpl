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

{{- range $key, $value := .Properties }}
### {{ $key }}

{{- range $value }}
- {{ .Property }}: `{{ .Value }}`
{{- end }}

{{- end }}
{{- end}}

{{ if .Dropins }}

{{- range $dropins := .Dropins }}
### Dropins - {{ $dropins.Directory }}/{{ $dropins.FileName }}

{{- range $key, $value := .Properties }}
#### {{ $key }}

{{- range $value }}
- {{ .Property }}: `{{ .Value }}`
{{- end }}
{{- end }}
{{- end }}
{{- end }}

{{ end }}
