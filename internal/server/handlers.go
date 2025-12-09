package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"sensorsProject/internal/rawData"
)

const TEMPLATES = "web/templates"

func (s *Server) IndexHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/users", http.StatusSeeOther)
}

func (s *Server) UsersHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("web/templates/index.html", "web/templates/users_page.html"))
	//tmplPath := filepath.Join(TEMPLATES, "users_page.html")
	//tmpl, err := template.ParseFiles(tmplPath)
	//if err != nil {
	//	http.Error(w, err.Error(), 500)
	//	return
	//}

	//data := map[string]interface{}{
	//	"Message": "Это HTML, сгенерированный сервером на Go!",
	//}

	err := tmpl.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		return
	}
}

func (s *Server) DatasetHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("web/templates/index.html", "web/templates/dataset_page.html"))

	err := tmpl.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		return
	}
}

func (s *Server) SendDataHandler(w http.ResponseWriter, r *http.Request) {
	var body rawData.RawData
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "bad json", 400)
		return
	}

	fmt.Println(body.Data)
	if err := rawData.AddRaw(body); err != nil {
		http.Error(w, "db error", 500)
		return
	}

	w.WriteHeader(200)
	_, err := w.Write([]byte("ok"))
	if err != nil {
		return
	}
}
