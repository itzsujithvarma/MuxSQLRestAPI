package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

const (
	db_user = "root"
	db_pwd  = "A1sujith."
	db_addr = "localhost"
	db_db   = "muxtest"
)

var s = fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", db_user, db_pwd, db_addr, db_db)
var db *sql.DB
var err error
var id int
var fn string
var ln string
var city string
var mobile string
var emp Employee

type Employee struct {
	Id         int    `json:"id"`
	First_Name string `json:"first_Name"`
	Last_Name  string `json:"last_Name"`
	City       string `json:"city"`
	Mobile     string `json:"mobile"`
}

func GetEmps(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	l := "select * from employee"
	res, err1 := db.Query(l)
	if err1 == nil {
		defer res.Close()
		var employees []Employee
		for res.Next() {

			err2 := res.Scan(&id, &fn, &ln, &city, &mobile)

			if err2 == nil {
				emp = Employee{Id: id, First_Name: fn, Last_Name: ln, City: city, Mobile: mobile}
				employees = append(employees, emp)
			}
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(employees)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func GetEmp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	l := "select * from employee where id =" + params["id"]
	res := db.QueryRow(l)

	err2 := res.Scan(&id, &fn, &ln, &city, &mobile)
	if err2 == nil {
		emp = Employee{Id: id, First_Name: fn, Last_Name: ln, City: city, Mobile: mobile}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(emp)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func UpdateEmp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	_ = json.NewDecoder(r.Body).Decode(&emp)
	l := "update employee set First_Name='" + emp.First_Name + "',Last_Name='" + emp.Last_Name + "',City='" + emp.City + "',Mobile='" + emp.Mobile + "' where id =" + params["id"]
	res, err2 := db.Exec(l)
	if err2 == nil {
		count, err3 := res.RowsAffected()
		if err3 == nil && count == 1 {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(emp)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func DelEmp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	l := "delete from employee where id =" + params["id"]
	res, err2 := db.Exec(l)
	if err2 == nil {
		count, err3 := res.RowsAffected()
		if err3 == nil && count == 1 {
			w.WriteHeader(http.StatusOK)

		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func CreateEmp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	_ = json.NewDecoder(r.Body).Decode(&emp)
	l := "insert into Employee(First_Name,Last_Name,City,Mobile) values('" + emp.First_Name + "','" + emp.Last_Name + "','" + emp.City + "','" + emp.Mobile + "')"
	res, err2 := db.Exec(l)
	if err2 == nil {
		count, err3 := res.RowsAffected()
		if err3 == nil && count == 1 {
			w.WriteHeader(http.StatusCreated)
			//json.NewEncoder(w).Encode(emp)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(l))
	}
}

func main() {
	db, err = sql.Open("mysql", s)
	if err != nil {
		fmt.Println("MYSQL ERRORRRRRRR")
	}
	defer db.Close()
	r := mux.NewRouter()
	r.HandleFunc("/api/emp", GetEmps).Methods("GET")
	r.HandleFunc("/api/emp/{id}", GetEmp).Methods("GET")
	r.HandleFunc("/api/emp/{id}", DelEmp).Methods("DELETE")
	r.HandleFunc("/api/emp/{id}", UpdateEmp).Methods("PUT")
	r.HandleFunc("/api/emp", CreateEmp).Methods("POST")
	http.ListenAndServe(":8081", r)
}
