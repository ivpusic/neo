{{define "testindex"}}
<html>
<head>
    <title>Test</title>
</head>
<body>
    {{template "testheader"}}
    Body {{.Name}}
    {{template "testfooter"}}
</body>
</html>
{{end}}