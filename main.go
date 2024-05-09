package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func conexionDB() (conexion *sql.DB) {
	Driver := "mysql"
	User := "root"
	Password := ""
	name := "system"

	conexion, err := sql.Open(Driver, User+":"+Password+"@tcp(127.0.0.1:3306)/"+name)
	if err != nil {
		panic(err.Error())

	}
	return conexion
}

var templates = template.Must(template.ParseGlob("templates/*"))

func main() {
	http.HandleFunc("/", Init)
	http.HandleFunc("/create", Create)
	http.HandleFunc("/insert", Insert)

	http.HandleFunc("/delete", Delete)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/update", Update)

	log.Println("Server started .......")

	http.ListenAndServe(":8080", nil)
}

type Employee struct {
	Id    int
	Name  string
	Email string
}

func Init(w http.ResponseWriter, r *http.Request) {
	establishedConnection := conexionDB()
	defer establishedConnection.Close()

	register, err := establishedConnection.Prepare("SELECT * FROM employee")

	if err != nil {
		panic(err.Error())
	}
	defer register.Close()

	employee := Employee{}
	arrayEmployees := []Employee{}
	rows, err := register.Query()

	if err != nil {
		fmt.Fprintf(w, "Error ejecutando consulta: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name, email string
		err := rows.Scan(&id, &name, &email)
		if err != nil {
			panic(err.Error())
		}
		employee.Id = id
		employee.Name = name
		employee.Email = email
		arrayEmployees = append(arrayEmployees, employee)

	}

	//fmt.Fprintf(w, "Hello, World!")

	templates.ExecuteTemplate(w, "init", arrayEmployees)
}

func Create(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "create", nil)
}

func Insert(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		name := r.FormValue("name")
		email := r.FormValue("email")

		establishedConnection := conexionDB()
		insertRegister, err := establishedConnection.Prepare("INSERT INTO employee(name, email) VALUES(?,?)")
		if err != nil {
			panic(err.Error())
		}
		insertRegister.Exec(name, email)
		http.Redirect(w, r, "/", 301)
	}

}

func Delete(w http.ResponseWriter, r *http.Request) {
	idEmployee := r.URL.Query().Get("id")

	establishedConnection := conexionDB()
	deleteRegister, err := establishedConnection.Prepare("DELETE FROM employee WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}

	deleteRegister.Exec(idEmployee)
	fmt.Println("Empleado eliminado: " + idEmployee)
	http.Redirect(w, r, "/", 301)
}

func Edit(w http.ResponseWriter, r *http.Request) {
	idEmployee := r.URL.Query().Get("id")
	fmt.Println(idEmployee)

	establishedConnection := conexionDB()
	defer establishedConnection.Close()

	stmt, err := establishedConnection.Prepare("SELECT * FROM employee WHERE id = ?")
	if err != nil {
		fmt.Fprintf(w, "Error en la consulta: %v", err)
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(idEmployee)
	if err != nil {
		fmt.Fprintf(w, "Error en la consulta: %v", err)
		return
	}
	defer rows.Close()

	var employee Employee
	found := false

	for rows.Next() {
		var id int
		var name, email string
		err = rows.Scan(&id, &name, &email)
		if err != nil {
			fmt.Fprintf(w, "Error al leer la consulta: %v", err)
			return
		}
		employee.Id = id
		employee.Name = name
		employee.Email = email
		found = true
	}

	if !found {
		fmt.Fprintf(w, "Error: Emplyee whit ID %s not found", idEmployee)
		return
	}
	templates.ExecuteTemplate(w, "edit", employee)

}

func Update(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := r.FormValue("id")
		name := r.FormValue("name")
		email := r.FormValue("email")

		establishedConnection := conexionDB()
		modifierRegister, err := establishedConnection.Prepare("UPDATE employee SET name = ?, email = ? WHERE id = ?")
		if err != nil {
			panic(err.Error())
		}
		modifierRegister.Exec(name, email, id)
		http.Redirect(w, r, "/", 301)
	}

}
