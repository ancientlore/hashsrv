package engine

import (
	"bytes"
)

// Help returns some basic help about the commands
func (e *Engine) Help() []byte {
	var b bytes.Buffer

	b.Write([]byte(htmlHeader))

	b.Write([]byte(`<h1>hashsrv</h1>`))
	b.Write([]byte(`hashsrv is a web service that performs hashing, encryption, encoding, and compression. Available functions include:`))

	tpl.ExecuteTemplate(&b, "Funcs", e.funcMap)

	tpl.ExecuteTemplate(&b, "Vars", e.values)

	b.Write([]byte(htmlFooter))

	return b.Bytes()
}
