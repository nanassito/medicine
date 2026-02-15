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
	<style>
		.medicine-cards { padding: 0 0 1rem; }
		.medicine-card {
			display: block;
			padding: 1.25rem 1.5rem;
			margin-bottom: 0.75rem;
			background: #f0f0f0;
			border: 1px solid #ccc;
			border-radius: 8px;
			text-decoration: none;
			color: #333;
			font-size: 1.1rem;
			transition: background 0.15s, border-color 0.15s, box-shadow 0.15s;
		}
		.medicine-card:hover {
			background: #e0e0e0;
			border-color: #999;
			box-shadow: 0 2px 8px rgba(0,0,0,0.1);
		}
	</style>
</head>
<body>
	<h1>All Medicines</h1>
	<div class="medicine-cards">
		{{ range .Medicines }}
		<a class="medicine-card" href="./{{.}}">{{.}}</a>
		{{ end }}
	</div>
</body>
</html>
`))
