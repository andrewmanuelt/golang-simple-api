package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Todo struct {
	gorm.Model
	Subject  string `json:"subject"`
	Complete string `json:"complete"`
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", show)
	r.HandleFunc("/show", show)
	r.HandleFunc("/create", create)
	r.HandleFunc("/edit", edit)
	r.HandleFunc("/delete", delete)

	fmt.Println("Server running on 127.0.0.1:9000")

	log.Fatal(http.ListenAndServe(":9000", r))
}

func database() *gorm.DB {
	dsn := "root:@tcp(127.0.0.1:3306)/todo?parseTime=true"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&Todo{})

	return db
}

func show(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json")

	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)

		data := map[string]interface{}{
			"status": "Method not allowed",
		}

		json.NewEncoder(w).Encode(data)

		return
	}

	var count int64

	db := database()

	db.Model(&Todo{}).Count(&count)

	if count == 0 {
		data := map[string]interface{}{
			"status": "OK",
			"data":   "No data",
		}

		json.NewEncoder(w).Encode(data)

		return
	}

	var datas []Todo

	db.Find(&datas)

	data := map[string]interface{}{
		"status": "OK",
		"data":   datas,
	}

	json.NewEncoder(w).Encode(data)
}

func create(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json")

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)

		data := map[string]interface{}{
			"status": "Method not allowed",
		}

		json.NewEncoder(w).Encode(data)

		return
	}

	var todos Todo
	var datas Todo

	json.NewDecoder(r.Body).Decode(&todos)

	db := database()

	db.Create(&Todo{
		Subject:  todos.Subject,
		Complete: todos.Complete,
	})

	db.Last(&datas)

	data := map[string]interface{}{
		"status": "",
		"data":   &datas,
	}

	json.NewEncoder(w).Encode(data)
}

func edit(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "application/json")

	if r.Method != "PATCH" {
		data := map[string]interface{}{
			"status": "Method not allowed",
		}

		json.NewEncoder(w).Encode(data)

		return
	}

	id := r.URL.Query().Get("id")

	if id == "" {
		data := map[string]interface{}{
			"status": "Id parameter are required",
		}

		json.NewEncoder(w).Encode(data)

		return
	}

	db := database()

	var update_data []Todo
	var update_form Todo

	json.NewDecoder(r.Body).Decode(&update_form)

	db.Model(&Todo{}).Where("id", id).Updates(&Todo{
		Subject:  update_form.Subject,
		Complete: update_form.Complete,
	})

	db.Model(&Todo{}).Where("id", id).Find(&update_data)

	data := map[string]interface{}{
		"status": "OK",
		"data":   &update_data,
	}

	json.NewEncoder(w).Encode(data)
}

func delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "DELETE" {
		w.WriteHeader(http.StatusMethodNotAllowed)

		data := map[string]interface{}{
			"status": "Method not allowed",
		}

		json.NewEncoder(w).Encode(data)

		return
	}

	id := r.URL.Query().Get("id")

	if id == "" {
		w.WriteHeader(http.StatusMethodNotAllowed)

		data := map[string]interface{}{
			"status": "Id parameter are required",
		}

		json.NewEncoder(w).Encode(data)

		return
	}

	db := database()

	var delete []Todo

	db.Where("id = ?", id).Delete(&delete)

	var todos []Todo

	db.Model(&Todo{}).Find(&todos)

	data := map[string]interface{}{
		"status": "OK",
		"data":   &todos,
	}

	json.NewEncoder(w).Encode(data)
}
