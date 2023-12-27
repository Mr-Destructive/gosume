package api

import (
	"html/template"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		templateStr := `<!DOCTYPE html>
        <html>
            <head>
                <title></title>
                <script src="https://unpkg.com/htmx.org@1.9.6"></script>
            </head>
            <body>
                <h1>GoSume</h1>
            </body>
        </html>`

		temp := template.New("index")
		t, err := temp.Parse(templateStr)
		t.Execute(w, nil)
		if err != nil {
			panic(err)
		}
	} else {
        
	}
}


