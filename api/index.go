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
        <button type="button" hx-get="/?section=custom-fields">Add Custom Fields</button>
    </div>

    <input type="submit" value="Submit">
</form>
`

const EducationForm = `
<label for="education-degree">Degree:</label>
<input type="text" id="education-degree" name="education-degree">
<button type="button" hx-get="/?section=education">+</button>
`

const ExperienceForm = `
<label for="experience-company">Company:</label>
<input type="text" id="experience-company" name="experience-company">
<label for="experience-position">Position:</label>
<input type="text" id="experience-position" name="experience-position">
<button type="button" hx-get="/?section=experience">+</button>
`

const SkillsForm = `
<label for="skill">Skill:</label>
<input type="text" id="skill" name="skill">
<button type="button" hx-get="/?section=skills">+</button>
`

const CustomFieldsForm = `
<label for="custom-field">Custom Field Key:</label>
<input type="text" id="custom-key" name="custom-key">
<label for="custom-field-value">Custom Field Value:</label>
<input type="text" id="custom-field-value" name="custom-field-value">
<button type="button" hx-get="/?section=custom-fields">+</button>
`

type FormData struct {
	Name                string
	Bio                 string
	EducationDegrees    []string
	ExperienceCompanies []string
	ExperiencePositions []string
	Skills              []string
    CustomFields        map[string]string
}

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
            case "custom-fields":
                formContent = CustomFieldsForm
			default:
				http.Error(w, "Invalid section", http.StatusBadRequest)
				return
			}
			templateStr := fmt.Sprintf(`
				<div id="form-container">
					%s
				</div>
			`, formContent)

			renderTemplate(w, templateStr)
			return
		}

		templateStr := fmt.Sprintf(`
			%s
			<div id="form-container">
				%s
			</div>
			%s
		`, HTML_TEMPL_START, ResumeForm, HTML_TEMPL_END)

		renderTemplate(w, templateStr)
	} else {
		formData := parseFormData(r)

		formContent := renderFormContent(formData)

		templateStr := fmt.Sprintf(`
			%s
			<div id="form-container">
				%s
			</div>
			%s
		`, HTML_TEMPL_START, formContent, HTML_TEMPL_END)

		renderTemplate(w, templateStr)
	}
}

func parseFormData(r *http.Request) FormData {
	err := r.ParseForm()
	if err != nil {
		return FormData{}
	}
    customFields := make(map[string]string)
    for i, k := range r.Form["custom-key"] {
        customFields[k] = r.Form["custom-field-value"][i]
    }

	return FormData{
		Name:                r.Form["name"][0],
		Bio:                 r.Form["bio"][0],
		EducationDegrees:    r.Form["education-degree"],
		ExperienceCompanies: r.Form["experience-company"],
		ExperiencePositions: r.Form["experience-position"],
		Skills:              r.Form["skill"],
        CustomFields:        customFields,
	}
}

func renderFormContent(data FormData) string {
	formContent := fmt.Sprintf(`
		<h1>Resume</h1>
		<hr>
		<h2>Name: %s</h2>
		<p>Bio: %s</p>
	`, data.Name, data.Bio)

	formContent += "<hr><h2>Education</h2>"
	for _, degree := range data.EducationDegrees {
		formContent += fmt.Sprintf(`
			<p>Education: %s</p>
		`, degree)
	}

	formContent += "<hr><h2>Experience</h2>"
	for i := range data.ExperienceCompanies {
		formContent += fmt.Sprintf(`
			<p>Company: %s</p>
			<p>Position: %s</p>
		`, data.ExperienceCompanies[i], data.ExperiencePositions[i])
	}

	formContent += "<hr><h2>Skills</h2>"
	for _, skill := range data.Skills {
		formContent += fmt.Sprintf(`
			<p>Skills: %s</p>
		`, skill)
	}
    formContent += "<hr><h2>Custom Fields</h2>"
    for k, v := range data.CustomFields {
        formContent += fmt.Sprintf(`
            <p>%s: %s</p>
        `, k, v)
    }

	return formContent
}

func renderTemplate(w http.ResponseWriter, templateStr string) {
	temp := template.New("index")
	t, err := temp.Parse(templateStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, nil)
}

func main() {
	http.HandleFunc("/", Handler)
	http.ListenAndServe(":8080", nil)
}
