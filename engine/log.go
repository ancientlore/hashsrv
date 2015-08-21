package engine

import (
	"fmt"
	"html/template"
	"log"
)

const (
	htmlHeader = `<html>
<head>
<title>hashsrv</title>
<style>
body {
	font-family: "Lucida Sans Unicode", "Lucida Grande", Sans-Serif;
	font-size: 14px;
}
h2 {
	font-family: "Lucida Sans Unicode", "Lucida Grande", Sans-Serif;
	font-size: 16px;
}
table
{
	font-family: "Lucida Sans Unicode", "Lucida Grande", Sans-Serif;
	font-size: 12px;
	margin: 45px;
	width: 800px;
	text-align: left;
	border-collapse: collapse;
	border: 1px solid #69c;
}
th
{
	padding: 12px 17px 12px 17px;
	font-weight: normal;
	font-size: 14px;
	color: #039;
	border-bottom: 1px dashed #69c;
}
td
{
	padding: 7px 17px 7px 17px;
	color: #669;
}
td.fixed
{
	width: 260px;
	padding: 7px 17px 7px 17px;
	color: #669;
	word-break: break-all;
}
tbody tr:hover td
{
	color: #339;
	background: #d0dafd;
}
</style>
</head>
<body>`

	htmlFooter = `
</body>
</html>`

	templates = `
{{define "Text"}}<h2>{{.}}</h2>{{end}}
{{define "Vars"}}Variables:<table><thead><tr><th>Name</th><th>Length</th><th>Text</th><th>Bytes</th></tr></thead><tbody>
{{range $n, $v := .}}<tr><td>{{$n}}</td><td>{{len $v}}</td><td class="fixed">{{$v | printf "%+q"}}</td><td>{{ $v | printf "% x"}}</td></tr>{{end}}
</tbody></table>
{{end}}
{{define "Stack"}}<table><thead><tr><th>Position</th><th>Length</th><th>Text</th><th>Bytes</th></tr></thead><tbody>
{{range $i, $v := .}}<tr><td>{{$i}}</td><td>{{len $v}}</td><td class="fixed">{{$v | printf "%+q"}}</td><td>{{ $v | printf "% x"}}</td></tr>{{end}}
</tbody></table>
{{end}}
{{define "Funcs"}}<table><thead><tr><th>Stack In</th><th>Function</th><th>Stack Out</th><th>Description</th></tr></thead><tbody>
{{range $k, $v := .}}<tr><td>{{$v.In}}</td><td><b>{{$k}}</b></td><td>{{$v.Out}}</td><td>{{$v.Desc}}</td></tr>{{end}}
</tbody></table>
{{end}}
`
)

var (
	tpl = template.Must(template.New("templates").Parse(templates))
)

// Log writes information to the debug log
func (e *Engine) Log(parms ...interface{}) {
	log.Print(parms...)
	if e.DebugMode {
		// ignoring errors
		tpl.ExecuteTemplate(e.logBuf, "Text", fmt.Sprint(parms...))
	}
}

// Logf writes formatted information to the debug log
func (e *Engine) Logf(format string, parms ...interface{}) {
	log.Printf(format, parms...)
	if e.DebugMode {
		// ignoring errors
		tpl.ExecuteTemplate(e.logBuf, "Text", fmt.Sprintf(format, parms...))
	}
}

func (e *Engine) LogValues() {
	for k, v := range e.values {
		log.Printf("%s = %+q [% x]", k, string(v), v)
	}
	if e.DebugMode {
		tpl.ExecuteTemplate(e.logBuf, "Vars", e.values)
	}
}

func (e *Engine) LogStack() {
	log.Print("Stack length: ", e.stack.Len())
	if e.DebugMode {
		a := e.stack.ToArray()
		tpl.ExecuteTemplate(e.logBuf, "Stack", a)
	}
}
