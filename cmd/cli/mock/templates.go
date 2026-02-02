package mock

import (
	"fmt"

	"github.com/nraghuveer/vibecast/lib/db"
	"github.com/nraghuveer/vibecast/lib/models"
)

type Template = models.Template

var defaultTemplates = []Template{
	{
		ID:      "tech-visionary",
		Name:    "Tech Visionary",
		Topic:   "The future of artificial intelligence",
		Persona: "A Silicon Valley tech founder with deep knowledge of AI and machine learning",
	},
	{
		ID:      "startup-founder",
		Name:    "Startup Journey",
		Topic:   "Building a successful startup from scratch",
		Persona: "An experienced entrepreneur who has built and sold multiple companies",
	},
	{
		ID:      "wellness-guru",
		Name:    "Wellness Expert",
		Topic:   "Holistic health and modern wellness practices",
		Persona: "A wellness coach with expertise in nutrition, mindfulness, and fitness",
	},
	{
		ID:      "creative-artist",
		Name:    "Creative Mind",
		Topic:   "The creative process and finding inspiration",
		Persona: "A multi-disciplinary artist who works across music, visual arts, and writing",
	},
	{
		ID:      "food-explorer",
		Name:    "Culinary Journey",
		Topic:   "World cuisines and stories behind food",
		Persona: "A chef and food writer who has traveled the world exploring culinary traditions",
	},
}

func InitializeDefaultTemplates() {
	for _, t := range defaultTemplates {
		exists, err := db.TemplateExists(t.ID)
		if err == nil && !exists {
			if err := db.CreateTemplate(t); err != nil {
				fmt.Printf("Warning: failed to initialize default template %s: %v\n", t.ID, err)
			}
		}
	}
}

func GetTemplates() []Template {
	dbTemplates, err := db.GetAllTemplates()
	if err != nil {
		return defaultTemplates
	}

	templates := make([]Template, 0, len(dbTemplates))
	for _, dt := range dbTemplates {
		templates = append(templates, Template{
			ID:      dt.ID,
			Name:    dt.Name,
			Topic:   dt.Topic,
			Persona: dt.Persona,
		})
	}

	return templates
}

func GetCustomTemplates() []Template {
	all := GetTemplates()
	custom := make([]Template, 0)

	for _, t := range all {
		isDefault := false
		for _, dt := range defaultTemplates {
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

func AddTemplate(t Template) error {
	return db.CreateTemplate(t)
}

func UpdateTemplate(t Template) error {
	return db.UpdateTemplate(t)
}

func DeleteTemplate(id string) error {
	return db.DeleteTemplate(id)
}
