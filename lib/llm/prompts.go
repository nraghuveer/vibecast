package llm

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

type PromptLoader struct {
	baseDir string
}

func NewPromptLoader(baseDir string) *PromptLoader {
	return &PromptLoader{baseDir: baseDir}
}

func (p *PromptLoader) RenderFile(relPath string, data any) (string, error) {
	content, err := p.readFile(relPath)
	if err != nil {
		return "", err
	}

	tmpl, err := template.New(filepath.Base(relPath)).Option("missingkey=error").Parse(content)
	if err != nil {
		return "", fmt.Errorf("parse prompt template %s: %w", relPath, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("render prompt template %s: %w", relPath, err)
	}

	return buf.String(), nil
}

func (p *PromptLoader) readFile(relPath string) (string, error) {
	// 1) Try baseDir relative to current working directory.
	path1 := filepath.Join(p.baseDir, relPath)
	if b, err := os.ReadFile(path1); err == nil {
		return string(b), nil
	}

	// 2) Try baseDir relative to the executable directory.
	exe, err := os.Executable()
	if err == nil {
		path2 := filepath.Join(filepath.Dir(exe), p.baseDir, relPath)
		if b, err := os.ReadFile(path2); err == nil {
			return string(b), nil
		}
	}

	return "", errors.New("prompt file not found: " + filepath.Join(p.baseDir, relPath))
}
