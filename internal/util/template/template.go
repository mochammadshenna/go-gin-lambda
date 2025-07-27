package template

import (
	"ai-service/internal/util/logger"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var templates *template.Template

// Init initializes the HTML templates
func Init() {
	// Get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		logger.Error(context.Background(), "failed to get working directory", err)
		return
	}

	// Define template directory path
	templateDir := filepath.Join(wd, "internal", "util", "template")

	// Parse all HTML templates with timestamp to avoid caching
	templates = template.New(fmt.Sprintf("templates_%d", time.Now().Unix())).Funcs(template.FuncMap{
		"percent": func(a, b int64) int {
			if b == 0 {
				return 0
			}
			return int((float64(a) / float64(b)) * 100)
		},
		"add": func(a, b int) int {
			return a + b
		},
		"sub": func(a, b int) int {
			return a - b
		},
		"mul": func(a, b int) int {
			return a * b
		},
		"div": func(a, b int) int {
			if b == 0 {
				return 0
			}
			return a / b
		},
		"safe": func(s string) template.HTML {
			return template.HTML(s)
		},
		"toJSON": func(v interface{}) template.HTML {
			data, err := json.Marshal(v)
			if err != nil {
				return template.HTML("{}")
			}
			return template.HTML(string(data))
		},
	})

	// Parse all HTML files in the template directory
	_, err = templates.ParseGlob(filepath.Join(templateDir, "*.html"))
	if err != nil {
		logger.Error(context.Background(), "failed to parse templates", err)
		return
	}

	logger.Info(context.Background(), "templates initialized successfully")
}

// GetTemplates returns the parsed templates
func GetTemplates() *template.Template {
	return templates
}

// ExecuteTemplate executes a specific template
func ExecuteTemplate(templateName string, data interface{}) (string, error) {
	if templates == nil {
		return "", fmt.Errorf("templates not initialized")
	}

	var buf strings.Builder
	err := templates.ExecuteTemplate(&buf, templateName, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
