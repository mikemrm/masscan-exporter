{{/*
  port-protocols takes a list of int/string ports and returns a yaml map of tcp and udp ports.

  Output:
    tcp:
      - "80"
      - "443"
    udp:
      - "53"
*/}}
{{- define "port-protocols" }}
  {{- $tcp := list }}
  {{- $udp := list }}

  {{- range default . list }}
    {{- if hasSuffix "/tcp" (toString .) }}
      {{- $tcp = append $tcp (trimSuffix "/tcp" (toString .)) }}
    {{- else if hasSuffix "/udp" (toString .) }}
      {{- $udp = append $udp (trimSuffix "/udp" (toString .)) }}
    {{- else }}
      {{- $tcp = append $tcp (toString .) }}
      {{- $udp = append $udp (toString .) }}
    {{- end }}
  {{- end }}
  {{- toYaml (dict
      "tcp" $tcp
      "udp" $udp
    )
  }}
{{- end }}
