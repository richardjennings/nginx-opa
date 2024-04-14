
{{/*
Expand the name of the chart.
*/}}
{{- define "opa-nginx.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "opa-nginx.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Generate certificates for opa-nginx server
*/}}
{{- define "opa-nginx.gen-certs" -}}
{{- $altNames := list ( printf "%s.%s" (include "opa-nginx.name" .) .Release.Namespace ) ( printf "%s.%s.svc" (include "opa-nginx.name" .) .Release.Namespace ) -}}
{{- $ca := genCA "opa-nginx-ca" 365 -}}
{{- $cert := genSignedCert ( include "opa-nginx.name" . ) nil $altNames 365 $ca -}}
tls.crt: {{ $cert.Cert | b64enc }}
tls.key: {{ $cert.Key | b64enc }}
{{- end -}}