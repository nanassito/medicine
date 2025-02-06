package templates

import (
	"html/template"
)

var List = template.Must(template.New("List").Funcs(template.FuncMap{}).Parse(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>all Medicines</title>
	<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/purecss@3.0.0/build/pure-min.css" integrity="sha384-X38yfunGUhNzHpBaEBsWLO+A0HDYOQi8ufWDkZ0k9e0eXz/tH3II7uKZ9msv++Ls" crossorigin="anonymous">
</head>
<body>
	<h1>All Medicines</h1>
	<div class="pure-g">
		<ul>
		{{ range .Medicines }}
			<li><a href="./{{.}}">{{.}}</a></li>
		{{ end }}
		</ul>
	</div>
</body>
</html>
`))
