package auth

import (
	"bytes"
	"context"
	"embed"
	"html/template"
)

//go:embed templates/*
var templateFS embed.FS

var templates = template.Must(template.ParseFS(templateFS, "templates/*"))

type Service struct {
}

func NewService() (*Service, error) {
	return &Service{}, nil
}

func (s *Service) ImplicitFlow(ctx context.Context, clientID, redirectURL string) ([]byte, error) {
	buf := &bytes.Buffer{}
	err := templates.ExecuteTemplate(
		buf,
		"implicit_flow.html",
		map[string]interface{}{"clientID": clientID, "redirectURL": redirectURL},
	)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
