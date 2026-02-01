{{- define "weather-client.name" -}}
weather-client
{{- end -}}

{{- define "weather-client.fullname" -}}
{{- printf "%s" (include "weather-client.name" .) -}}
{{- end -}}

{{- define "weather-client.labels" -}}
app.kubernetes.io/name: {{ include "weather-client.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}
