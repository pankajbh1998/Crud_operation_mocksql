package main

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"log"
	"reflect"
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

func TestUpdateFunc(t *testing.T){
	fdb,mock,err:=sqlmock.New()
	if err != nil {
		t.Fatalf("Can't connect to a mock database")
	}
	dbh:=DatabaseHandler{fdb}
	db:=dbh.Db
	defer db.Close()
	testCase:=[]struct{
		input string
		output Employee
		expectedRetrunvalue  int64
		err error
	}{
		{
			"1",Employee{"1", "Pankaj Sharma", 22, "M", "1"},1,nil,
		},
		{
			"2",Employee{"2", "Rudra Bhardwaj", 21, "M", "2"},0,errors.New("User enterted the existing same data"),
		},
		{
			"3",Employee{"3", "Vivek Sharma", 26, "M", "3"},1,nil,
		},
	}
	//str:=[]string{"Id","Name","Age","Gender","Role"}
	for i,tc:= range testCase {
		mock.ExpectExec("Update Employee*").WithArgs(tc.output.Name,tc.output.Age,tc.output.Gender,tc.output.Role,tc.input).WillReturnResult(sqlmock.NewResult(0,tc.expectedRetrunvalue))
		result,err:=dbh.UpdateData(tc.input,tc.output)
		//if err!=nil && tc.err!=nil && err.Error() != tc.err.Error() {
		if !reflect.DeepEqual(err,tc.err){
			t.Fatalf("Failed at %v\nExpected Error : %v\nActual Error : %v\n", i+1, tc.err.Error(), err.Error())
			//t.Logf("%v\n%v\n%v\n%T\n%T",i+1,err,tc.err,err,tc.err)
		}
		if result != tc.output {
			t.Fatalf("Failed at %v\nExpected Output : %v\nActual Output : %v\n", i+1, tc.output, result)
		}
		t.Logf("Passed at %v\n", i+1)
	}

}