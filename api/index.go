package api

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"

	"github.com/jung-kurt/gofpdf"
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
		if len(queryparam) > 0 && queryparam["section"] != nil {
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
        queryparam := r.URL.Query()
        if len(queryparam) > 0 && queryparam["pdf"][0] == "true" {
            data := parseFormData(r)
            pdf, err := generatePDF(data)
            if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
            w.Header().Set("Content-Type", "application/pdf")
            w.Write(pdf.Bytes())
            return
        }else{

		templateStr := fmt.Sprintf(`
			%s
            <form action="/?pdf=true" method="POST">
			<div id="form-container">
				%s
			</div>
                <input type="submit" value="Generate PDF">
            </form>
			%s
		`, HTML_TEMPL_START, formContent, HTML_TEMPL_END)

		renderTemplate(w, templateStr)
        }
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
        <input type='text' hidden name='name' value='%s'>
		<p>Bio: %s</p>
        <input type='text' hidden name='bio' value='%s'>
	`, data.Name, data.Name, data.Bio, data.Bio)

	formContent += "<hr><h2>Education</h2>"
	for _, degree := range data.EducationDegrees {
		formContent += fmt.Sprintf(`
        <input type='text' hidden name='education-degree' value='%s'>
			<p>Education: %s</p>
		`, degree, degree)
	}

	formContent += "<hr><h2>Experience</h2>"
	for i := range data.ExperienceCompanies {
		formContent += fmt.Sprintf(`
            <input type='text' hidden name='experience-company' value='%s'>
			<p>Company: %s</p>
            <input type='text' hidden name='experience-position' value='%s'>
			<p>Position: %s</p>
		`, data.ExperienceCompanies[i], data.ExperienceCompanies[i], data.ExperiencePositions[i],  data.ExperiencePositions[i])
	}

	formContent += "<hr><h2>Skills</h2>"
	for _, skill := range data.Skills {
		formContent += fmt.Sprintf(`
			<p>Skills: %s</p>
            <input type='text' hidden name='skill' value='%s'>
		`, skill, skill)
	}
    formContent += "<hr><h2>Custom Fields</h2>"
    for k, v := range data.CustomFields {
        formContent += fmt.Sprintf(`
            <input type='text' hidden name='custom-key' value='%s'>
            <input type='text' hidden name='custom-field-value' value='%s'>
            <p>%s: %s</p>
        `, k, v, k, v)
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


func generatePDF(data FormData) (*bytes.Buffer, error) {
	var pdfBuffer bytes.Buffer

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)

	pdf.Cell(40, 10, "Resume")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 12)

	pdf.Cell(40, 10, fmt.Sprintf("Name: %s", data.Name))
	pdf.Ln(8)

	pdf.Cell(40, 10, fmt.Sprintf("Bio: %s", data.Bio))
	pdf.Ln(8)

	pdf.Cell(40, 10, "Education:")
	pdf.Ln(8)

	for _, degree := range data.EducationDegrees {
		pdf.Cell(40, 10, fmt.Sprintf("- %s", degree))
		pdf.Ln(8)
	}

	pdf.Cell(40, 10, "Experience:")
	pdf.Ln(8)

	for i := range data.ExperienceCompanies {
		pdf.Cell(40, 10, fmt.Sprintf("- Company: %s", data.ExperienceCompanies[i]))
		pdf.Ln(8)

		pdf.Cell(40, 10, fmt.Sprintf("  Position: %s", data.ExperiencePositions[i]))
		pdf.Ln(8)
	}

	pdf.Cell(40, 10, "Skills:")
	pdf.Ln(8)

	for _, skill := range data.Skills {
		pdf.Cell(40, 10, fmt.Sprintf("- %s", skill))
		pdf.Ln(8)
	}

	pdf.Cell(40, 10, "Custom Fields:")
	pdf.Ln(8)

	for k, v := range data.CustomFields {
		pdf.Cell(40, 10, fmt.Sprintf("- %s: %s", k, v))
		pdf.Ln(8)
	}

	err := pdf.Output(&pdfBuffer)
	if err != nil {
		return nil, err
	}

	return &pdfBuffer, nil
}


func main() {
	http.HandleFunc("/", Handler)
	http.ListenAndServe(":8080", nil)
}
