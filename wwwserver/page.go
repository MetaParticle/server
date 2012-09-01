package wwwserver

import (
	"github.com/MetaParticle/metaparticle/entity"
	
	"bytes"
	"html/template"
)

type Page struct {
    Name string
    Title string
    Css []string
    Scripts []string
    Body string
}

func (p *Page) ApplyPageTemplate(player *entity.Player) {
	p.Title, _ = templatifyString(p.Title, player)
	p.Body, _ = templatifyString(p.Body, player)
}

func templatifyString(s string, player *entity.Player) (string, error) {
	t := template.New("The Template")
	t, err := t.Parse(s)
	
	if err == nil {
		buf := bytes.NewBufferString("")
		t.Execute(buf, player)
		return buf.String(), err
	}
	return s, err
}
