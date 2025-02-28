package templates

import (
	"html/template"
)

var MedicineFor = template.Must(template.New("MedicineFor").Funcs(template.FuncMap{}).Parse(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>{{.MedicineName}}</title>
	<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/purecss@3.0.0/build/pure-min.css" integrity="sha384-X38yfunGUhNzHpBaEBsWLO+A0HDYOQi8ufWDkZ0k9e0eXz/tH3II7uKZ9msv++Ls" crossorigin="anonymous">
</head>
<body>
	<h1>{{.MedicineName}} - {{.Who.Name}}</h1>
	<img class="pure-img" src="{{.Who.PhotoUrl}}" alt="{{.Who.Name}}">
	<div style="text-align:center; padding-top:10px; padding-bottom:10px; background-color:{{if .CanTake}}#60A561{{else if lt .WaitForPct 0.1}}#FFB400{{else}}#F4442E{{end}};">
		<p>{{.Reason}}</p>
		{{if .CanTake}}{{else}}<p>Do NOT take for another {{.WaitFor}}</p>{{end}}
	</div>
	<h3>Posology</h3>
	<ul>
		<li>Dose: {{.Posology.MaxDoses}} every {{.Posology.DoseInterval}}</li>
		<li>No more than {{.Posology.MaxDoses}} doses every {{.Posology.MaxDosesInterval}}</li>
	</ul>
	<a style="width: 100%" class="pure-button pure-button-primary" href="/{{.MedicineName}}/{{.Who.Name}}/take"><h2>Take</h2></a>
</body>
</html>
`))
