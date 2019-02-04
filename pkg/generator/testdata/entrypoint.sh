#!/bin/bash

set -e

COMMAND="fpm "
{{ range . }}
if [[ -n "${{ .EnvVar }}" ]]; then
    {{- if .HasInput }}
    COMMAND="$COMMAND {{ .Option }} ${{ .EnvVar }}"
    {{- else }}
    COMMAND="$COMMAND {{ .Option }}"
    {{- end }}
fi
{{ end }}
echo $COMMAND