package main

import (
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	route := mux.NewRouter()

	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	route.HandleFunc("/", home).Methods("GET")
	route.HandleFunc("/contact", contact).Methods("GET")
	route.HandleFunc("/project", project).Methods("GET")
	route.HandleFunc("/add-project", addProject).Methods("POST")
	route.HandleFunc("/detail/{index}", detail).Methods("GET")
	route.HandleFunc("/delete/{index}", delete).Methods("GET")
	route.HandleFunc("/edit-project/{id}", myProjectEdited).Methods("POST")
	route.HandleFunc("/form-edit-project/{index}", myProjectFormEditProject).Methods("GET")
	// r.HandleFunc("/add-contact", addContact).Methods("POST")
	fmt.Println("server on in port 8080")
	http.ListenAndServe("localhost:8080", route)
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("views/index.html")

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	card := map[string]interface{}{
		"Add": data,
	}

	tmpl.Execute(w, card)
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("views/contact.html")
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	tmpl.Execute(w, "")
}

func project(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	tmpl, err := template.ParseFiles("views/addProject.html")
	if err != nil {
		w.Write([]byte(err.Error()))
		
	}

	tmpl.Execute(w, "")
}

type Project struct {
	Name         string
	Start_date   string
	End_date     string
	Duration     string
	Desc         string
	Id          int
	Technologies []string
}

var data = []Project{}

// editdata

func addProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	var name = r.PostForm.Get("inputName")
	var start_date = r.PostForm.Get("startDate")
	var end_date = r.PostForm.Get("endDate")
	var desc = r.PostForm.Get("desc")
	var technologies []string
	technologies = r.Form["technologies"]

	// fmt.Println(technologies)
	// fmt.Println(start_date)

	layout := "2006-01-02"
	dateStart, _ := time.Parse(layout, start_date)
	dateEnd, _ := time.Parse(layout, end_date)

	hours := dateEnd.Sub(dateStart).Hours()
	daysInHours := hours / 24
	monthInDay := daysInHours / 30
	yearInMonth := monthInDay / 12 // Njir prettier ya

	var duration string
	var month, _ float64 = math.Modf(monthInDay)
	var year, _ float64 = math.Modf(yearInMonth)

	if year > 0 {
		duration = strconv.FormatFloat(year, 'f', 0, 64) + " Years"
		// fmt.Println(year, " Years")
	} else if month > 0 {
		duration = strconv.FormatFloat(month, 'f', 0, 64) + " Months"
		// fmt.Println(month, " Months")
	} else if daysInHours > 0 {
		duration = strconv.FormatFloat(daysInHours, 'f', 0, 64) + " Days"
		// fmt.Println(daysInHours, " Days")
	} else if hours > 0 {{{  }}
		duration = strconv.FormatFloat(hours, 'f', 0, 64) + " Hours"
		// fmt.Println(hours, " Hours")
	} else {
		duration = "0 Days"
		// fmt.Println("0 Days")
	}

	// fmt.Println(daysInHours)
	// fmt.Println(month)
	// fmt.Println(year)

	var newData = Project{
		Name:         name,
		Start_date:   start_date,
		End_date:     end_date,
		Duration: duration,
		Desc:         desc,
		Technologies: technologies,
		Id:          len(data),
	}

	data = append(data, newData)
	// fmt.Println(data)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func detail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	tmpl, err := template.ParseFiles("views/detail.html")
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	var Detail = Project{}

	index, _ := strconv.Atoi(mux.Vars(r)["index"])
	// fmt.Println(index)

	for i, data := range data {
		if index == i {
			Detail = Project{
				Name:       data.Name,
				Start_date: data.Start_date,
				End_date:   data.End_date,
				Desc:       data.Desc,
				Duration:  data.Duration,
			}
		}
	}

	data := map[string]interface{}{
		"Details": Detail,
	}
	// fmt.Println(data)
	tmpl.Execute(w, data)
}

func delete(w http.ResponseWriter, r *http.Request) {
	index, _ := strconv.Atoi(mux.Vars(r)["index"])

	data = append(data[:index], data[index+1:]...)

	http.Redirect(w, r, "/", http.StatusFound)
}
func myProjectEdited(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	err := r.ParseForm()

	if err != nil {
		log.Fatal(err)
	}

	Name := r.PostForm.Get("Name")
	Start_date := r.PostForm.Get("Start_date")
	End_date := r.PostForm.Get("End_date")
	Desc := r.PostForm.Get("Desc")

	editDataForm := Project{
		Name: Name,
		Start_date:   Start_date,
		End_date:     End_date,
		Desc: Desc,
		Id:          id,
		// Duration:    time.Now().String(),
	}

	data[id] = editDataForm

	fmt.Println(data)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)

}

func myProjectFormEditProject(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("views/myProjectFormEditProject.html")

	index, _ := strconv.Atoi(mux.Vars(r)["index"])

	ProjectEdit := Project{}

	for i, data := range data {
		if index == i {
			ProjectEdit = Project{
				Name: data.Name,
				Start_date:   data.Start_date,
				End_date:     data.End_date,
				Desc: data.Desc,
				Id:          data.Id,
			}
		}
	}

	response := map[string]interface{}{
		"Project": ProjectEdit,
	}

	if err == nil {
		tmpl.Execute(w, response)
	} else {
		panic(err)
	}
}


// func addContact(w http.ResponseWriter, r *http.Request){
// 	err := r.ParseForm()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	name := r.PostForm.Get("name")
// 	email := r.PostForm.Get("email")
// 	phone := r.PostForm.Get("phone")
// 	subject := r.PostForm.Get("select")
// 	description := r.PostForm.Get("description")

// 	fmt.Println("Name : " + name)
// 	fmt.Println("Email : " + email)
// 	fmt.Println("Phone Number : " + phone)
// 	fmt.Println("Subject : " + subject)
// 	fmt.Println("Description : " + description)


// }
