package api

import (
    "fmt"
    "net/http"
    "text/template"
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

const ResumeForm = `
    <form action="/" method="POST" >
        <label for="name">Name:</label>
        <input type="text" id="name" name="name">

        <label for="bio">Bio:</label>
        <input type="text" id="bio" name="bio">

        <div  hx-swap="outerHTML">
            <button type="button" hx-get="/?section=education">Add Education</button>
            <button type="button" hx-get="/?section=experience">Add Experience</button>
            <button type="button" hx-get="/?section=skills">Add Skills</button>
        </div>

        <input type="submit" value="Submit">
    </form>
`

const EducationForm = `
    <h2>Add Education</h2>
    <label for="education-degree">Degree:</label>
    <input type="text" id="education-degree" name="education-degree">
    <button type="button" hx-post="/">Back to Main Form</button>
`

const ExperienceForm = `
    <h2>Add Experience</h2>
    <label for="company">Company:</label>
    <input type="text" id="company" name="company">
    <label for="position">Position:</label>
    <input type="text" id="position" name="position">
    <button type="button" hx-post="/">Back to Main Form</button>
`

const SkillsForm = `
    <h2>Add Skills</h2>
    <label for="skill">Skill:</label>
    <input type="text" id="skill" name="skill">
    <button type="button" hx-post="/">Back to Main Form</button>
`

func Handler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        queryparam := r.URL.Query()
        formContent := ""
        if len(queryparam) > 0 {
            section := queryparam["section"][0]
            switch section {
            case "education":
                formContent = EducationForm
            case "experience":
                formContent = ExperienceForm
            case "skills":
                formContent = SkillsForm
            default:
                http.Error(w, "Invalid section", http.StatusBadRequest)
            }
        templateStr := fmt.Sprintf(`
            <div id="form-container">
                %s
            </div>
        `, formContent)

        temp := template.New("index")
        t, err := temp.Parse(templateStr)
        t.Execute(w, nil)
        if err != nil {
            panic(err)
        }
            
        return
        }
        templateStr := fmt.Sprintf(`
            %s
            <div id="form-container">
                %s
            </div>
            %s
        `, HTML_TEMPL_START, ResumeForm, HTML_TEMPL_END)

        temp := template.New("index")
        t, err := temp.Parse(templateStr)
        t.Execute(w, nil)
        if err != nil {
            panic(err)
        }
    } else {
        formData := map[string]string{}
        err := r.ParseForm()
        if err != nil {
            panic(err)
        }
        for k, v := range r.Form {
            formData[k] = v[0]
        }

        formContent := ""
        name := formData["name"]
        bio := formData["bio"]
        education := formData["education-degree"]
        experience := formData["experience"]
        skills := formData["skills"]
        formContent = fmt.Sprintf(`
            <p>Name: %s</p>
            <p>Bio: %s</p>
            <p>Education: %s</p>
            <p>Experience: %s</p>
            <p>Skills: %s</p>
        `, name, bio, education, experience, skills)

        templateStr := fmt.Sprintf(`
            %s
            <div id="form-container">
                %s
            </div>
            %s
        `, HTML_TEMPL_START, formContent, HTML_TEMPL_END)

        temp := template.New("index")
        t, err := temp.Parse(templateStr)
        t.Execute(w, nil)
        if err != nil {
            panic(err)
        }
    }
}

func main() {
    http.HandleFunc("/", Handler)

    http.ListenAndServe(":8080", nil)
}
