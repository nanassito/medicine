package templates

import (
	"html/template"
)

var MedicineOverview = template.Must(template.New("MedicineOverview").Funcs(template.FuncMap{}).Parse(`
<!DOCTYPE html>
<html>
<head>
	<title>{{.MedicineName}}</title>
</head>
<body>
	<h1>{{.MedicineName}}</h1>
	{{ range .CanTake }}
		<p>{{ .Who }} can {{if .CanTake }} {{else}}NOT{{end}} take this medicine because {{ .Reason }}</p>
	{{ end }}
</body>
</html>
`))
