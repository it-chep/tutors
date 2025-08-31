{{- $short := (shortname .Name "err" "res" "sqlstr" "db" "XOLog" "item") -}}
{{- $table := (schema .Schema .Table.TableName) -}}
{{- if .Comment -}}
// {{ .Comment }}
{{- else -}}
// {{ .Name }} represents a row from '{{ $table }}'.
{{- end }}
type {{ .Name }} struct {
{{- range .Fields }}
	{{ if eq (retype .Type) "custom.Jsonb" -}}{{ .Name }} []byte `db:"{{ .Col.ColumnName }}" json:"{{ .Col.ColumnName }}"` // {{ .Col.ColumnName }}
	{{ else if eq (retype .Type) "custom.JSON" -}}{{ .Name }} []byte `db:"{{ .Col.ColumnName }}" json:"{{ .Col.ColumnName }}"` // {{ .Col.ColumnName }}
    {{- else -}}{{ .Name }} {{ retype .Type }} `db:"{{ .Col.ColumnName }}" json:"{{ .Col.ColumnName }}"` // {{ .Col.ColumnName }}
	{{- end -}}
{{- end }}
}

{{ $prefix := .Name }}
{{ $table_name := .Table.TableName }}

// zero{{ $prefix }} zero value of dto
var zero{{ $prefix }} = {{ $prefix }}{}

func  (t {{ $prefix }}) IsEmpty() bool {
    return reflect.DeepEqual(t, zero{{ $prefix }})
}
