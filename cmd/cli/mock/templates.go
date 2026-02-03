package mock

import "github.com/nraghuveer/vibecast/lib/models"

var DefaultTemplates = []models.Template{
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
