package main

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"log"
	"testing"
)


func mockDbConnection()*sql.DB{
	db, mock, err := sqlmock.New() //create mock database
	if err != nil {
		log.Fatalf("Error comes when connecting a Stub Database : %v", err)
	}

	// Here we are creating rows in our mocked database.
	rows := sqlmock.NewRows([]string{"Id", "Name", "Age", "Gender", "Role"}).
		AddRow("1", "Pankaj Sharma", 22, "M", "1").
		AddRow("2", "Rudra Bhardwaj", 21, "M", "2").
		AddRow("3", "Vivek Sharma", 26, "M", "3")
	mock.ExpectQuery("^select (.+) from Employe*").WillReturnRows(rows)
	mock.ExpectQuery("^select (.+) from Employe*").WithArgs("id").WillReturnRows(rows)
	return db
}
func TestMock(t *testing.T) {
	db:=mockDbConnection()
	defer db.Close()
	//ctx:=context.TODO()
	//ans:=ReadData(ctx,db)
	ans,err:=ReadData(db)
	if err !=nil {
		t.Fatal(err.Error())
	}
	expectedOutput:=[]Employee{
		{"1", "Pankaj Sharma", 22, "M", "1"},
		{"2", "Rudra Bhardwaj", 21, "M", "2"},
		{"3", "Vivek Sharma", 26, "M", "3"},
	}
	//t.Log(ans)
	//t.Log(expectedOutput)
	for i,val:=range expectedOutput {
		if val!=ans[i]{
			t.Fatalf("Failed at %v\nExpected Output : %v\nActual Output : %v\n",i+1,val,ans[i])
		}
		t.Logf("Passed at %v\n",i+1)
	}
}

func Testmockid(t *testing.T){
	db:=mockDbConnection()
	defer db.Close()
	testCases:=[]struct{
		input string
		output Employee
	}{
		{"1",{"1", "Pankaj Sharma", 22, "M", "1"}},
		{"2",{"2", "Rudra Bhardwaj", 21, "M", "2"}},
		{"3",{"3", "Vivek Sharma", 26, "M", "3"}},
	}
	for i,tc:=range testCases {
		res,err:=ReadDataid(db, tc.input)
		if err != nil {
			t.Fatal(err.Error())
		}
		if res != tc.output {
			t.Fatalf("Failed at %v\nExpected Output : %v\nActual Output : %v\n",i+1,tc.output,res)
		}
		t.Logf("Passed at %v\n",i+1)
	}

}