package data

import (
	"fmt"

	"github.com/nraghuveer/vibecast/cmd/cli/mock"
	"github.com/nraghuveer/vibecast/lib/db"
	"github.com/nraghuveer/vibecast/lib/models"
)

func InitializeDefaultTemplates(database *db.DB) {
	for _, t := range mock.DefaultTemplates {
		exists, err := database.TemplateExists(t.ID)
		if err == nil && !exists {
			if err := database.CreateTemplate(t); err != nil {
				fmt.Printf("Warning: failed to initialize default template %s: %v\n", t.ID, err)
			}
		}
	}
}

func GetTemplates(database *db.DB) []models.Template {
	dbTemplates, err := database.GetAllTemplates()
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

func GetCustomTemplates(database *db.DB) []models.Template {
	all := GetTemplates(database)
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

func AddTemplate(database *db.DB, t models.Template) error {
	return database.CreateTemplate(t)
}

func UpdateTemplate(database *db.DB, t models.Template) error {
	return database.UpdateTemplate(t)
}

func DeleteTemplate(database *db.DB, id string) error {
	return database.DeleteTemplate(id)
}
