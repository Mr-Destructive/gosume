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
    <title>GoSume</title>
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <link href="https://cdn.jsdelivr.net/npm/tailwindcss@2.2.19/dist/tailwind.min.css" rel="stylesheet">
    <script src="https://unpkg.com/htmx.org@1.9.6"></script>
</head>
<body class="bg-gray-100 p-8">
  <div class="max-w-lg mx-auto bg-white p-6 rounded shadow-md">
        <h1 class="text-2xl font-semibold mb-4">GoSume Resume</h1>
`

const HTML_TEMPL_END = `
</body>
</html>
`

const ResumeForm = `
        <form action="/" method="POST">
            <div class="mb-4">
                <label for="name" class="block text-gray-700 text-sm font-bold mb-2">Name:</label>
                <input type="text" id="name" name="name" class="w-full border p-2 rounded">
            </div>

            <div class="mb-4">
                <label for="bio" class="block text-gray-700 text-sm font-bold mb-2">Bio:</label>
                <input type="text" id="bio" name="bio" class="w-full border p-2 rounded">
            </div>

            <div class="mb-4" hx-swap="outerHTML">
                <button type="button" hx-get="/?section=education"
                    class="bg-blue-500 text-white px-4 py-2 rounded mr-2 hover:bg-blue-600">Add Education</button>
                <button type="button" hx-get="/?section=experience"
                    class="bg-blue-500 text-white px-4 py-2 rounded mr-2 hover:bg-blue-600">Add Experience</button>
                <button type="button" hx-get="/?section=skills"
                    class="bg-blue-500 text-white px-4 py-2 rounded mr-2 hover:bg-blue-600">Add Skills</button>
                <button type="button" hx-get="/?section=custom-fields"
                    class="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600">Add Custom Fields</button>
            </div>

            <input type="submit" value="Submit"
                class="bg-green-500 text-white px-4 py-2 rounded hover:bg-green-600">
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
