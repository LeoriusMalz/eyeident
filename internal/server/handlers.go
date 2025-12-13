package server

import (
	"encoding/json"
	"eyeident/internal/rawData"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Body struct {
	Count   string            `json:"count"`
	Dataset []rawData.Dataset `json:"dataset"`
}

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

func (s *Server) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := rawData.GetUsers()
	if err != nil {
		http.Error(w, "Error getting users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		http.Error(w, "Error users encoding", http.StatusInternalServerError)
		return
	}
}

func (s *Server) GetDatasetHandler(w http.ResponseWriter, r *http.Request) {
	startDate, _ := strconv.ParseInt(r.URL.Query().Get("startDate"), 10, 64)
	endDate, _ := strconv.ParseInt(r.URL.Query().Get("endDate"), 10, 64)
	ids := strings.Split(r.URL.Query().Get("id"), ",")
	types := strings.Split(r.URL.Query().Get("type"), ",")

	log.Println(startDate, endDate, ids, types)
	cnt, err := rawData.GetDataset(ids, types, startDate, endDate, "/data/dataset.csv")
	if err != nil {
		http.Error(w, "Error getting dataset:"+err.Error(), http.StatusInternalServerError)
		return
	}

	dataset, err := rawData.ReadCSVPreview("/data/dataset.csv", 10)
	if err != nil {
		log.Println("Error reading preview:" + err.Error())
		http.Error(w, "Error reading preview"+err.Error(), 500)
		return
	}

	var body = Body{cnt, dataset}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(body)
	if err != nil {
		log.Println("Error sending dataset:" + err.Error())
		http.Error(w, "Error dataset encoding", http.StatusInternalServerError)
		return
	}
}

func (s *Server) ConnectUserHandler(w http.ResponseWriter, r *http.Request) {
	var body rawData.UserData
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "bad json", 400)
		return
	}

	log.Println("User", body.Id, "waiting to connection!")
	if err := rawData.AddUser(body); err != nil {
		http.Error(w, "db error", 500)
		return
	}

	w.WriteHeader(200)
	_, err := w.Write([]byte("ok"))
	if err != nil {
		return
	}
}

func (s *Server) DisconnectUserHandler(w http.ResponseWriter, r *http.Request) {
	var body rawData.UserData
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "bad json", 400)
		return
	}

	log.Println("User", body.Id, "wants to disconnect!")
	if err := rawData.RemoveUser(body); err != nil {
		http.Error(w, "db error", 500)
		return
	}

	w.WriteHeader(200)
	_, err := w.Write([]byte("ok"))
	if err != nil {
		return
	}
}

func (s *Server) SendDataHandler(w http.ResponseWriter, r *http.Request) {
	var body rawData.SensorPacket
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "bad json", 400)
		return
	}

	log.Println("Received", len(body.Samples), "samples from user", body.Id)
	if err := rawData.Add2RawSet(body); err != nil {
		http.Error(w, "db error", 500)
		return
	}

	w.WriteHeader(200)
	_, err := w.Write([]byte("ok"))
	if err != nil {
		return
	}
}

func (s *Server) GetDatasetParamsHandler(w http.ResponseWriter, r *http.Request) {
	params, err := rawData.GetParams()
	if err != nil {
		http.Error(w, "Error getting params", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(params)
	if err != nil {
		http.Error(w, "Error params encoding", http.StatusInternalServerError)
		return
	}
}
