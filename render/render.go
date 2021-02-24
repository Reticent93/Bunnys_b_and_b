package render

import (
	"bytes"
	"github.com/Reticent93/Bunnys_b_and_b/config"
	"github.com/Reticent93/Bunnys_b_and_b/models"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)



var app *config.AppConfig

var functions = template.FuncMap{

}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

//NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}


// RenderTemplate renders templates using html/template
func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	
	var tc map[string]*template.Template
	
	if app.UseCache {
		//get the template cached from the app config
		
		tc = app.TemplateCache
	}else {
		tc, _ = CreateTemplateCache()
		
	}
	
	

	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("Could not get template from template cache")
	}
	
	buf := new(bytes.Buffer)
	
	td = AddDefaultData(td)
	
	_ = t.Execute(buf, td)
	
	buf.WriteTo(w)

}

//CreateTemplateCache creates a template cache as a map
func CreateTemplateCache()  (map[string]*template.Template, error){
	myCache := map[string]*template.Template{}
	
	pages, err := filepath.Glob("./templates/*.page.tmpl.html")
	
	if err != nil {
		return myCache, err
	}
	
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache ,err
		}
		
		matches, err := filepath.Glob("./templates/*.layout.html")
		if err != nil {
			return myCache,err
		}
		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.html")
			if err != nil {
				return myCache,err
			}
		}
		
		myCache[name] = ts
	}
	return myCache, nil
}