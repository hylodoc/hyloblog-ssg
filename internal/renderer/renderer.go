package renderer

import "fmt"

type Renderer struct {
}

func NewRenderer(title, theme string) (*Renderer, error) {
	return nil, fmt.Errorf("NewRenderer NOT IMPLEMENTED")
}

func (r *Renderer) Render(title, body string) (string, error) {
	return "", fmt.Errorf("Render NOT IMPLEMENTED")
}
