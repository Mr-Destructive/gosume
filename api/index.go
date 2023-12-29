package api

import (
    "fmt"
	"html/template"
	"net/http"
)

const HTML_TEMPL_START = `
<!DOCTYPE html>
<html>
    <head>
        <title></title>
        <script src="https://unpkg.com/htmx.org@1.9.6"></script>
    </head>
    <body>
        <h1>GoSume</h1>
`

const HTML_TEMPL_END = `
    </body>
</html>
`

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		templateStr :=  fmt.Sprintf(`
        %s
        <form action="/" method="post">
                    <label for="username">Name:</label>
                    <input type="text" id="name" name="name">
                    <label for="bio" >Bio:</label>
                    <input type="text" id="bio" name="bio">
                    <input type="submit" value="Submit">
                </form>
        %s
        `, HTML_TEMPL_START, HTML_TEMPL_END)

		temp := template.New("index")
		t, err := temp.Parse(templateStr)
		t.Execute(w, nil)
		if err != nil {
			panic(err)
		}
	} else {
        name := r.FormValue("name")
        bio := r.FormValue("bio")
        // render pdf for resume
        templateStr := fmt.Sprintf(`
        %s
                <p>Name: %s</p>
                <p>Bio: %s</p>
        %s
        `, HTML_TEMPL_START, name, bio, HTML_TEMPL_END)
        temp := template.New("index")
        t, err := temp.Parse(templateStr)
        t.Execute(w, nil)
        if err != nil {
            panic(err)
        }

	}
}


