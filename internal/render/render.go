package render

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"text/template"
	"time"

	"github.com/cindysurjawann/gobookings/internal/config"
	"github.com/cindysurjawann/gobookings/internal/models"
	"github.com/justinas/nosurf"
)

// functions that available to my golang template
var functions = template.FuncMap{
	"humanDate": HumanDate,
}
var app *config.AppConfig
var pathToTemplates = "./templates"

// sets the config for the template package
func NewRenderer(a *config.AppConfig) {
	app = a
}

// returns time in YYY-MM-DD format
func HumanDate(t time.Time) string {
	return t.Format("2006-01-02")
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.CSRFToken = nosurf.Token(r)
	if app.Session.Exists(r.Context(), "user_id") {
		td.IsAuthenticated = 1
	}
	return td
}

// to render HTML files
func Template(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {
	var tc map[string]*template.Template

	if app.UseCache {
		//get the template cache from the app config
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	//get requested template from cache
	t, ok := tc[tmpl]
	if !ok {
		//if template not exist
		return errors.New("can't get template from cache")
	}
	buf := new(bytes.Buffer)

	td = AddDefaultData(td, r)

	err := t.Execute(buf, td)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	//render the template
	_, err = buf.WriteTo(w)
	if err != nil {
		fmt.Println("Error writing template to browser", err)
		return err
	}
	return nil
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	//get all of the files named *.page.html from ./templates
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.html", pathToTemplates))
	if err != nil {
		return myCache, err
	}

	//range through all files ending with *.page.html
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.page.html", pathToTemplates))
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.page.html", pathToTemplates))
		}
		myCache[name] = ts
	}
	return myCache, nil
}

//SIMPLE TEMPLATE CACHE
// var cache = make(map[string]*template.Template)

// func RenderTemplate(w http.ResponseWriter, t string) {
// 	var tmpl *template.Template
// 	var err error

// 	//check to see if we already have the template in our cache
// 	_, inMap := cache[t]
// 	if !inMap {
// 		//need to create the template
// 		log.Println("creating template and adding to cache")
// 		err = createTemplateCache(t)
// 		if err != nil {
// 			log.Println(err)
// 		}
// 	} else {
// 		//we have the template in the cache
// 		log.Println("using cached template")
// 	}

// 	tmpl = cache[t]
// 	err = tmpl.Execute(w, nil)
// 	if err != nil {
// 		log.Println(err)
// 	}
// }

// func createTemplateCache(t string) error {
// 	templates := []string{
// 		fmt.Sprintf("./templates/%s", t),
// 		"./templates/base.layout.html",
// 	}

// 	//parse the template
// 	tmpl, err := template.ParseFiles(templates...)
// 	if err != nil {
// 		return err
// 	}

// 	//add template to cache
// 	cache[t] = tmpl
// 	return nil
// }
