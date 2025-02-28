package templates

import (
	"html/template"
)

var MedicineOverview = template.Must(template.New("MedicineOverview").Funcs(template.FuncMap{}).Parse(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>{{.MedicineName}}</title>
	<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/purecss@3.0.0/build/pure-min.css" integrity="sha384-X38yfunGUhNzHpBaEBsWLO+A0HDYOQi8ufWDkZ0k9e0eXz/tH3II7uKZ9msv++Ls" crossorigin="anonymous">
</head>
<body>
	<h1>{{.MedicineName}}</h1>
	<div class="pure-g">
		{{ range .People }}
			<div class="pure-u-1-2">
				<figure>
					<a href="/{{$.MedicineName}}/{{.Name}}">
						<img class="pure-img" src="{{.PhotoUrl}}" alt="{{.Name}}">
					</a>
				</figure>
			</div>
		{{ end }}
	</div>
</body>
</html>
`))
