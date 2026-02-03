package data

import (
	"fmt"

	"github.com/nraghuveer/vibecast/cmd/cli/mock"
	"github.com/nraghuveer/vibecast/lib/db"
	"github.com/nraghuveer/vibecast/lib/models"
)

func InitializeDefaultTemplates() {
	for _, t := range mock.DefaultTemplates {
		exists, err := db.TemplateExists(t.ID)
		if err == nil && !exists {
			if err := db.CreateTemplate(t); err != nil {
				fmt.Printf("Warning: failed to initialize default template %s: %v\n", t.ID, err)
			}
		}
	}
}

func GetTemplates() []models.Template {
	dbTemplates, err := db.GetAllTemplates()
	if err != nil {
		return mock.DefaultTemplates
	}

	templates := make([]models.Template, 0, len(dbTemplates))
	for _, dt := range dbTemplates {
		templates = append(templates, models.Template{
			ID:      dt.ID,
			Name:    dt.Name,
			Topic:   dt.Topic,
			Persona: dt.Persona,
		})
	}

	return templates
}

func GetCustomTemplates() []models.Template {
	all := GetTemplates()
	custom := make([]models.Template, 0)

	for _, t := range all {
		isDefault := false
		for _, dt := range mock.DefaultTemplates {
			if t.ID == dt.ID {
				isDefault = true
				break
			}
		}
		if !isDefault {
			custom = append(custom, t)
		}
	}

	return custom
}

func AddTemplate(t models.Template) error {
	return db.CreateTemplate(t)
}

func UpdateTemplate(t models.Template) error {
	return db.UpdateTemplate(t)
}

func DeleteTemplate(id string) error {
	return db.DeleteTemplate(id)
}
