package main

import (
	"github.com/DATA-DOG/go-sqlmock"
	"log"
	"testing"
)

func TestReadAll(t *testing.T) {
	fdb, mock, err := sqlmock.New() //create mock database
	if err != nil {
		log.Fatalf("Error comes when connecting a mock Database : %v", err)
	}
	dbh := DatabaseHandler{fdb}
	db := dbh.Db
	defer db.Close()
	// Here we are creating rows in our mocked database.
	rows := sqlmock.NewRows([]string{"Id", "Name", "Age", "Gender", "Role"})
	rows.AddRow("1", "Pankaj Sharma", 22, "M", "1")
	rows.AddRow("2", "Rudra Bhardwaj", 21, "M", "2")
	rows.AddRow("3", "Vivek Sharma", 26, "M", "3")

	mock.ExpectQuery("^select (.+) from Employee*").WillReturnRows(rows)

	//ctx:=context.TODO()
	//ans:=ReadData(ctx,db)
	ans, err := dbh.ReadDataAll()
	if err != nil {
		t.Fatal(err.Error())
	}
	expectedOutput := []Employee{
		{"1", "Pankaj Sharma", 22, "M", "1"},
		{"2", "Rudra Bhardwaj", 21, "M", "2"},
		{"3", "Vivek Sharma", 26, "M", "3"},
	}
	//t.Log(ans)
	//t.Log(expectedOutput)
	for i, val := range expectedOutput {
		if val != ans[i] {
			t.Fatalf("Failed at %v\nExpected Output : %v\nActual Output : %v\n", i+1, val, ans[i])
		}
		t.Logf("Passed at %v\n", i+1)
	}
}

func TestReadById(t *testing.T) {
	fdb, mock, err := sqlmock.New() //create mock database
	if err != nil {
		log.Fatalf("Error comes when connecting a Stub Database : %v", err)
	}
	dbh := DatabaseHandler{fdb}
	db := dbh.Db
	defer db.Close()
	ExpectedOutput := []Employee {
		{"1", "Pankaj Sharma", 22, "M", "1"},
		{"2", "Rudra Bhardwaj", 21, "M", "2"},
		{"3", "Vivek Sharma", 26, "M", "3"},
	}
	str:=[]string{"Id", "Name", "Age", "Gender", "Role"}
	for i, tc := range ExpectedOutput {
		row:=mock.NewRows(str).AddRow(tc.Id,tc.Name,tc.Age,tc.Gender,tc.Role)
		mock.ExpectQuery("^select (.+) from Employee*").WithArgs(tc.Id).WillReturnRows(row)
		res, err := dbh.ReadDataId(tc.Id)
		if err != nil {
			t.Fatal(err.Error())
		}
		if res != tc {
			t.Fatalf("Failed at %v\nExpected Output : %v\nActual Output : %v\n", i+1, tc, res)
		}
		t.Logf("Passed at %v\n", i+1)
	}
}

//func TestUpdateFunc(t *testing.T){
//	fdb,mock,err:=sqlmock.New()
//	if err != nil {
//		t.Fatalf("Can't connect to a mock database")
//	}
//	dbh:=DatabaseHandler{fdb}
//	db:=dbh.Db
//
//
//}