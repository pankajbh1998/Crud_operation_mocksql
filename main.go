package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	_ "github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type Employee struct {
	Id string			`json:="id"`
	Name string 		`json:="name"`
	Age int				`json:="age"`
	Gender string		`json:="gender"`
	Role string			`json:="role"`
}
type DatabaseHandler struct{
	Db *sql.DB
}

func (dbh* DatabaseHandler)ReadDataAll()([]Employee,error){
	db:=dbh.Db
	var ans []Employee
	res,err:=db.Query("select Id,Name,Age,Gender,Role from Employee")
	defer res.Close()
	if err !=nil{
		return ans,err
	}
	for res.Next(){
		var emp Employee
		res.Scan(&emp.Id,&emp.Name,&emp.Age,&emp.Gender,&emp.Role)
		ans=append(ans,emp)
	}
	return ans,nil
}
func (dbh* DatabaseHandler)ReadDataId(id string)(Employee,error){
	db:=dbh.Db
	var emp Employee
	res:=db.QueryRow("select ID,Name,Age,Gender,Role from Employee where Id=?",id)
	if res.Err() !=nil {
		return emp,errors.New("Data is not found")
	}
	res.Scan(&emp.Id,&emp.Name,&emp.Age,&emp.Gender,&emp.Role)
	return emp,nil
}
func (dbh *DatabaseHandler)dbConnection() error{
	user:="Pankaj"
	password:="Pankaj@123"
	ip:="127.0.0.1"
	dbName:="Company"
	//db,err:=sql.Open("mysql","Pankaj:Pankaj@123@tcp(127.0.0.1)/Company")
	db,err:=sql.Open("mysql",fmt.Sprintf("%v:%v@tcp(%v)/%v",user,password,ip,dbName))
	if err !=nil {
		return(errors.New("Cannot connectd with database"))
	}
	dbh.Db=db
	return nil
}
func (dbh *DatabaseHandler)CreateData(emp Employee)Employee{
	db:=dbh.Db
	insert, _ := db.Exec("insert into Employee(Name,Age,Gender,Role) values(?,?,?,?)", emp.Name, emp.Age, emp.Gender, emp.Role)
	num,_:=insert.LastInsertId()
	emp.Id=strconv.Itoa(int(num))
	return emp
}
func CreateDataHandler(w http.ResponseWriter,r* http.Request){
	var dbh DatabaseHandler
	err:=dbh.dbConnection()
	defer dbh.Db.Close()
	if err != nil {
		http.Error(w,err.Error(),http.StatusInternalServerError)
	}
	var emp Employee
	err = json.NewDecoder(r.Body).Decode(&emp)
	log.Print(emp)
	if err != nil {
		http.Error(w,"Input is not in correct format",http.StatusBadRequest)
	}
	log.Print(emp)
	emp=dbh.CreateData(emp)
	post, _ := json.Marshal(emp)
	w.Write(post)
}

func ReadDataAllHandler(w http.ResponseWriter,r* http.Request){
	w.Header().Set("content-type","application/json")
	var dbh DatabaseHandler
	err:=dbh.dbConnection()
	defer dbh.Db.Close()
	if err != nil {
		http.Error(w,err.Error(),http.StatusInternalServerError)
	}
	result,err:=dbh.ReadDataAll()
	if err !=nil {
		http.Error(w,err.Error(),http.StatusBadRequest)
	}
	post, err := json.Marshal(result)
	if err !=nil {
		http.Error(w,err.Error(),http.StatusInternalServerError)
	} else {
		w.Write(post)
	}
}


func ReadDataId(w http.ResponseWriter,r *http.Request) {
	w.Header().Set("content-type","application/json")
	var dbh DatabaseHandler
	err:=dbh.dbConnection()
	defer dbh.Db.Close()
	if err != nil {
		http.Error(w,err.Error(),http.StatusInternalServerError)
	}
	vars:=mux.Vars(r)
	id:=vars["id"]
	result,err:=dbh.ReadDataId(id)
	if err != nil {
		http.Error(w,err.Error(),http.StatusBadRequest)
	}
	post, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.Write(post)
	}
}

func ReadDataQueryHandler(w http.ResponseWriter,r* http.Request){
	var dbh DatabaseHandler
	err:=dbh.dbConnection()
	defer dbh.Db.Close()
	if err != nil {
		http.Error(w,err.Error(),http.StatusInternalServerError)
	}
	//vars:=mux.Vars(r)
	//id:=vars["query"]
	result,err:=dbh.ReadDataAll()
	if err != nil {
		http.Error(w,err.Error(),http.StatusBadRequest)
	}
	post,err :=json.Marshal(result)
	if err != nil {
		http.Error(w,err.Error(),http.StatusInternalServerError)
	} else {
		w.Write(post)
	}
}
func (dbh *DatabaseHandler)checkuser(id string)bool {
	db:=dbh.Db
	res:=db.QueryRow(fmt.Sprintf("select Id from Employee where Id= %v",id))
	var temp string
	res.Scan(&temp)
	if temp != "" {
		return true
	}
	return false
}
func (dbh *DatabaseHandler)UpdateData(id string,emp Employee)(Employee,error){
	db:= dbh.Db
	var emp1 Employee
	result,err:=db.Exec("Update Employee set Name=?,Age=?,Gender=?,Role=? where id =?",emp.Name,emp.Age,emp.Gender,emp.Role,id)
	if err != nil {
		return emp1,errors.New("Input is not in appropriate format")
	}
	emp.Id=id
	if num,_:=result.RowsAffected();int(num) ==0 {
		return emp,errors.New("User enterted the existing same data")
	}
	return emp,nil
}
func UpdateDataHandler(w http.ResponseWriter,r* http.Request) {
	w.Header().Set("content-text","application/json")
	var dbh DatabaseHandler
	err:=dbh.dbConnection()
	defer dbh.Db.Close()
	if err != nil {
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	vars := mux.Vars(r)
	id := vars["id"]
	if dbh.checkuser(id) == false {
		http.Error(w, "ID does not exist", http.StatusNotFound)
		return
	}
	var emp Employee
	err = json.NewDecoder(r.Body).Decode(&emp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Print(emp)
	result, err := dbh.UpdateData(id, emp)
	post, _ := json.Marshal(result)
	w.Write(post)
}
func (dbh *DatabaseHandler)DeleteData(id string)error{
	db:=dbh.Db
	result,err:=db.Exec("delete from Employee where id = ?",id)
	if num,_:=result.RowsAffected();err != nil || int(num)==0 {
		return errors.New("Id does not exist")
	}
	return nil
}
func DeleteDataHandler(w http.ResponseWriter,r* http.Request){
	var dbh DatabaseHandler
	err:=dbh.dbConnection()
	if err !=nil {
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	defer dbh.Db.Close()
	vars:=mux.Vars(r)
	id:=vars["id"]
	err=dbh.DeleteData(id)
	if err != nil	{
		http.Error(w,err.Error(),http.StatusBadRequest)
		return
	}
	w.Write([]byte("Deleted Successfully"))
	w.WriteHeader(http.StatusOK)
}


func main(){
	r:=mux.NewRouter()
	r.HandleFunc("/employee",CreateDataHandler).Methods("POST")
	r.HandleFunc("/employee",ReadDataAllHandler).Methods("GET")
	//r.HandleFunc("/employee/{id}",ReadDataIdHandler).Methods("GET")
	//r.HandleFunc("/employee/{query}",ReadDataQuery).Methods("GET")
	r.HandleFunc("/employee/{id}",UpdateDataHandler).Methods("PUT")
	r.HandleFunc("/employee/{id}",DeleteDataHandler).Methods("DELETE")
	http.ListenAndServe(":8080",r)
}

