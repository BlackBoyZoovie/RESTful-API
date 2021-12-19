package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "8080"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}
type Routes []Route

var routes = Routes{
	Route{
		"getEmployees",
		"GET",
		"/employees",
		getEmployees,
	},
	Route{
		"getEmployee",
		"GET",
		"/employee/Id",
		getEmployee,
	},
	Route{
		"addEmployee",
		"POST",
		"/employee/{add}",
		addEmployee,
	},
	Route{
		"updateEmployee",
		"PUT",
		"/employee/update",
		updateEmployee,
	},
	Route{
		"deleteEmployee",
		"DELETE",
		"/employee/delete",
		deleteEmployee,
	},
}

type Employee struct {
	Id        string `json:"Id"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}
type Employees []Employee

var employees []Employee
var employeesV1 []Employee
var employeesV2 []Employee

func init() {
	employees = Employees{
		Employee{Id: "1", FirstName: "Simi", LastName: "ilesanmi"},
		Employee{Id: "2", FirstName: "Gbolahan", LastName: "Fakorede"},
	}
	employeesV1 = Employees{
		Employee{Id: "1", FirstName: "Dami", LastName: "ilesanmi"},
		Employee{Id: "2", FirstName: "Fatoke", LastName: "Ademola"},
	}
	employeesV2 = Employees{
		Employee{Id: "1", FirstName: "Felix", LastName: "Djaho"},
		Employee{Id: "2", FirstName: "Asine", LastName: "Black"},
	}
}

func getEmployees(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/v1") {
		json.NewEncoder(w).Encode(employeesV1)
	} else if strings.HasPrefix(r.URL.Path, "/v2") {
		json.NewEncoder(w).Encode(employeesV2)
	} else {
		json.NewEncoder(w).Encode(employees)
	}
}

func getEmployee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	Id := vars["Id"]
	for _, employee := range employees {
		if employee.Id == Id {
			if err := json.NewEncoder(w).Encode(employee); err != nil {
				log.Print("Error getting requested employee :: ", err)
			}
		}
	}
}

func updateEmployee(w http.ResponseWriter, r *http.Request) {
	employee := Employee{}
	err := json.NewDecoder(r.Body).Decode(&employee)
	if err != nil {
		log.Print("Error occurred while decoding new employee :: ", err)
		return
	}
	var isUpsert = true
	for Idx, emp := range employees {
		if emp.Id == employee.Id {
			isUpsert = false
			log.Printf("Updating employee Id as :: %s with firstname as :: %s and lastname as :: %s", employee.Id, employee.FirstName, employee.LastName)
			employees[Idx].FirstName = employee.FirstName
			employees[Idx].LastName = employee.LastName
			break
		}
	}
	if isUpsert {
		log.Printf("Upserting employee Id :: %s with firstname as :: %s and lastname as :: %s", employee.Id, employee.FirstName, employee.LastName)
		employees = append(employees, Employee{Id: employee.Id, FirstName: employee.FirstName, LastName: employee.LastName})
		json.NewEncoder(w).Encode(employees)
	}
}

func addEmployee(w http.ResponseWriter, r *http.Request) {
	employee := Employee{}
	err := json.NewDecoder(r.Body).Decode(&employee)
	if err != nil {
		log.Print("Error occurred while decoding new employee :: ", err)
		return
	}
	log.Printf("Add employee Id :: %s with firstname as :: %s and lastname as :: %s", employee.Id, employee.FirstName, employee.LastName)
	employees = append(employees, Employee{Id: employee.Id, FirstName: employee.FirstName, LastName: employee.LastName})
	json.NewEncoder(w).Encode(employees)
}

func deleteEmployee(w http.ResponseWriter, r *http.Request) {
	employee := Employee{}
	err := json.NewDecoder(r.Body).Decode(&employee)
	if err != nil {
		log.Print("Error occurred while decoding new employee :: ", err)
		return
	}
	log.Printf("Deleting employee Id :: %s with firstname as :: %s and lastname as :: %s", employee.Id, employee.FirstName, employee.LastName)
	index := GetIndex(employee.Id)
	employees = append(employees[:index], employees[index+1])
	json.NewEncoder(w).Encode(employees)
}

func GetIndex(id string) int {
	for i := 0; i < len(employees); i++ {
		if employees[i].Id == id {
			return i
		}
	}
	return -1
}

func AddRoutes(r *mux.Router) *mux.Router {
	for _, route := range routes {
		r.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			HandlerFunc(route.HandlerFunc)
	}
	return r
}

func main() {
	muxRouter := mux.NewRouter().StrictSlash(true)
	r := AddRoutes(muxRouter)
	AddRoutes(muxRouter.PathPrefix("/v1").Subrouter())
	AddRoutes(muxRouter.PathPrefix("/v2").Subrouter())
	log.Println("Starting server at Port 8080...")
	err := http.ListenAndServe(CONN_HOST+":"+CONN_PORT, r)
	if err != nil {
		log.Fatal("Error starting server :: ", err)
		return
	}
}
