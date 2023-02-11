package chain

import (
	"bytes"
	"text/template"
)

// prompt template
type PromptTemplate struct {
	Template string
	Inputs   []string
}

func (p *PromptTemplate) Format(args map[string]string) (string, error) {
	tmpl := template.Must(template.New("template").Parse(p.Template))
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, args); err != nil {
		return "", err
	}
	return buf.String(), nil
}
