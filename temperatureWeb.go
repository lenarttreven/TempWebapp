package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"log"

	_ "github.com/lib/pq"
)

const (  
  host     = "localhost"
  port     = 5432
  user     = "postgres"
  password = "Try to guess it"
  dbname   = "lenart_demo"
)

var (
	globalDB *sql.DB
)

type Page struct {
    Title string
    Body  []byte
}

type Temp struct {
	ID int `json:"ID,string"` 
	Temperature int `json:"Temperature,string"` 
	Location string `json:"Location"`
}

type WebapiResponse struct {
	Status string //ok or error
	Description string
	Data interface{} `json:"data"`
}

func (p *Page) save() error {
    filename := p.Title + ".txt"
    return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
    filename := title + ".txt"
    body, err := ioutil.ReadFile(filename)
    if err != nil {
    	return nil, err
    }
    return &Page{Title: title, Body: body}, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
    t, err := template.ParseFiles(tmpl + ".html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
    err = t.Execute(w, p)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func handler(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Path[len("/"):]
    p, _ := loadPage(title)
    renderTemplate(w, "temperatures", p)
}


func webapiHandler(w http.ResponseWriter, r *http.Request) {
	method := r.FormValue("method")

	fmt.Println(method)
	var executionError error
	var data interface{}
	var description string

	switch method {
	case "insertTemp":
		data, executionError = insertTempHanlder(r)
		description = "Insertion successfully performed."
	case "updateTemp":
		data, executionError = updateTempHandler(r)
		description = "Temperature was successfully updated."
	case "deleteID":
		data, executionError = deleteIDHandler(r)
		description = "ID was successfully deleted."
	case "deleteTemp":
		data, executionError = deleteTempHandler(r)
		description = "Temperature was successfully deleted."
	case "deleteLoc":
		data, executionError = deleteLocHandler(r)
		description = "Location was successfully deleted."
	case "getData":
	 	data, executionError = getDataHandler(r)
	 	description = "Data successfully updated."
	}

	
	// Prepare response
	response := WebapiResponse{
		Status: "OK",
		Data: data,
		Description: description,
	}
	if executionError != nil{
		response.Status = "Error"
		response.Description = executionError.Error()
	}
	//write log
	writeLog(response.Status, response.Description)

	b, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

func insertTempHanlder(r *http.Request) (interface{}, error) {
	body, readErr := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if readErr != nil || body == nil{
		return nil, readErr
	}
	fmt.Println(string(body))
	
	var temp Temp
	err := json.Unmarshal(body, &temp)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	fmt.Println(temp.Temperature)
	fmt.Println(temp.ID)
	fmt.Println(temp.Location)
	err = insertTemp(temp.Temperature, temp.Location, globalDB)
	if err != nil {
		fmt.Println(err.Error())
	}
	return nil, err
}

func updateTempHandler(r *http.Request) (interface{}, error) {
	body, readErr := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if readErr != nil || body == nil{
		return  nil, readErr
	}

	var temp Temp
	err := json.Unmarshal(body, &temp)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	fmt.Println(temp.Temperature)
	fmt.Println(temp.ID)
	fmt.Println(temp.Location)
	err = updateTemp(temp.ID, temp.Temperature, globalDB)
	return nil, err
}

func deleteIDHandler(r *http.Request) (interface{}, error) {
	body, readErr := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if readErr != nil || body == nil{
		return  nil, readErr
	}

	var temp Temp
	err := json.Unmarshal(body, &temp)
	if err != nil {
		return nil, err
	}
	fmt.Println(temp.Temperature)
	fmt.Println(temp.ID)
	fmt.Println(temp.Location)
	err = deleteID(temp.ID, globalDB)
	return nil, err
}

func deleteTempHandler(r *http.Request) (interface{}, error) {
	body, readErr := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if readErr != nil || body == nil{
		return  nil, readErr
	}

	var temp Temp
	err := json.Unmarshal(body, &temp)
	if err != nil {
		return nil, err
	}
	fmt.Println(temp.Temperature)
	fmt.Println(temp.ID)
	fmt.Println(temp.Location)
	err = deleteTemp(temp.Temperature, globalDB)
	return nil, err
}

func deleteLocHandler(r *http.Request) (interface{}, error) {
	body, readErr := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if readErr != nil || body == nil{
		return nil, readErr
	}
	var temp Temp
	err := json.Unmarshal(body, &temp)
	if err != nil {
		return nil, err
	}
	fmt.Println(temp.Temperature)
	fmt.Println(temp.ID)
	fmt.Println(temp.Location)
	err = deleteLoc(temp.Location, globalDB)
	return nil, err

}

func getDataHandler(r *http.Request) (interface{}, error) {
	data, err := printAll(globalDB)
	if err != nil{
		return nil, err
	}
	return data, err 
}


func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
    "password=%s dbname=%s sslmode=disable",
    host, port, user, password, dbname)
  	db, err := sql.Open("postgres", psqlInfo)
  	if err != nil {
    	fmt.Println(err.Error())
  	}
  	globalDB = db
  	defer db.Close()
  	log.SetFlags(0)
  	
    http.HandleFunc("/", handler)
    http.HandleFunc("/webapi", webapiHandler)
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
    
    http.ListenAndServe(":8080", nil)



}