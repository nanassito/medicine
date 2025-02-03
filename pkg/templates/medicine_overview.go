package templates

import "html/template"

var MedicineOverview = template.Must(template.New("MedicineOverview").Parse(`
<!DOCTYPE html>
<html>
<head>
	<title>{{.Name}}</title>
</head>
<body>
	<h1>{{.Name}}</h1>
</body>
</html>
`))
