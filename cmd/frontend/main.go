package main

import (
	"ChiliOverFlow/pkg/config"
	"ChiliOverFlow/pkg/db"
	"fmt"
	"github.com/Masterminds/sprig"
	"github.com/go-resty/resty/v2"
	"github.com/julienschmidt/httprouter"
	"github.com/spf13/viper"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"time"
)

var (
	conf *viper.Viper
	version = "undefined"
	tmpl *template.Template
	client *resty.Client

)
func getAllTemplates() []string {
	templateFiles := []string{}
	files, err := ioutil.ReadDir("templates")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		templateFiles = append(templateFiles, path.Join("templates", f.Name()))
	}
	fmt.Println(templateFiles)
	return templateFiles
}


func main() {

	var err error
	conf, err = config.New()
	if err != nil {
		log.Fatal("Could not get configuration: ", err)
	}

	client = resty.New()
	dir, _ := os.Getwd()
	webRoot := dir

	tmpl = template.Must(template.New("dummy").Funcs(sprig.FuncMap()).ParseFiles(getAllTemplates()...))

	router := httprouter.New()


	router.GET("/", frontPageHandler)

	router.ServeFiles("/static/*filepath", http.Dir(path.Join(webRoot, "static")))
	imagePath := path.Join(conf.GetString("configPath"), "images")
	log.Println("looking in", imagePath)
	router.ServeFiles("/images/*filepath", http.Dir(imagePath))

	log.Fatal(http.ListenAndServe(":8080", router))
}

func frontPageHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	tmpl = template.Must(template.New("dummy").Funcs(sprig.FuncMap()).ParseFiles(getAllTemplates()...))

	chilis := []db.StockedChili{}
	_, _ = client.R().
		EnableTrace().
		SetResult(&chilis).
		Get("http://localhost:8081/inventory/available")
	//TODO error handling
	execute := tmpl.ExecuteTemplate(w, "main",chilis)
	if execute != nil {
		log.Println(execute)
	}
}

var myClient = &http.Client{Timeout: 10 * time.Second}