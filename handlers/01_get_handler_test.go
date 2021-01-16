package handlers

import (
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/token"
	"golang-fifa-world-cup-web-service/data"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"path"
	"reflect"
	"runtime"
	"testing"
)

func TestRootHandlerReturnsNoContentStatus(t *testing.T) {
	handler := http.HandlerFunc(RootHandler)
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusNoContent {
		t.Error("Did not return status 204 - No Content")
	}
}

func TestListWinnersSetsContentType(t *testing.T) {
	handler := http.HandlerFunc(ListWinners)
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/winners", nil)
	handler.ServeHTTP(rr, req)
	if ctype := rr.Header().Get("Content-Type"); ctype != "application/json" {
		t.Error("Did not set Content-Type response header to application/json")
	}
}

func TestListWinnersReturnsAllWinners(t *testing.T) {
	setup()

	handler := http.HandlerFunc(ListWinners)
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/winners", nil)
	handler.ServeHTTP(rr, req)

	body := rr.Body.String()
	var winners data.Winners
	json.Unmarshal([]byte(body), &winners)
	if len(winners.Winners) != 21 {
		t.Error("Did not return all winners from /winners")
	}
}

func TestListWinnersReadsYearQueryString(t *testing.T) {
	failingTestMessage := "Did not assign query string \"year\" to variable year"
	defer func() {
		if recover() != nil {
			t.Error(failingTestMessage)
		}
	}()
	hasAssignedToYear := false
	isReadingStringArg := false

	listWinnersFunc := getFuncDecl("ListWinners", "handlers.go")
	funcBody := listWinnersFunc.Body
	for _, funcStatement := range funcBody.List {
		if reflect.TypeOf(funcStatement).String() == "*ast.AssignStmt" {
			// Checks if `year` variable is being assigned to
			assignment := funcStatement.(*ast.AssignStmt)
			left := assignment.Lhs[0].(*ast.Ident)
			hasAssignedToYear = (left.Name == "year")

			// Checks if req.URL.Query().Get("year")
			// is being called
			right := assignment.Rhs[0].(*ast.CallExpr)
			rightArg := right.Args[0].(*ast.BasicLit)
			isReadingStringArg = (rightArg.Value == "\"year\"")

			if hasAssignedToYear && isReadingStringArg {
				break
			}
		}
	}
	if !hasAssignedToYear || !isReadingStringArg {
		t.Error(failingTestMessage)
	}
}

func TestListWinnersReturnsAllWinnersByYear(t *testing.T) {
	handler := http.HandlerFunc(ListWinners)
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/winners", nil)
	q := req.URL.Query()
	q.Add("year", "2018")
	req.URL.RawQuery = q.Encode()

	handler.ServeHTTP(rr, req)

	body := rr.Body.String()
	var winners data.Winners
	json.Unmarshal([]byte(body), &winners)
	if len(winners.Winners) != 1 {
		t.Error("Did not return winners filtered by year")
	}
}

func TestListWinnersReturnsBadRequestWhenInvalidYear(t *testing.T) {
	handler := http.HandlerFunc(ListWinners)
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/winners", nil)
	q := req.URL.Query()
	q.Add("year", "banana")
	req.URL.RawQuery = q.Encode()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Error("Did not return status 400 - Bad Request")
	}
}

func getFuncDecl(name string, filename string) *ast.FuncDecl {
	var found *ast.FuncDecl

	_, currentFile, _, _ := runtime.Caller(1)

	src, err := ioutil.ReadFile(path.Join(path.Dir(currentFile), filename))
	if err != nil {
		log.Fatal(err)
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		panic(err)
	}

	for _, decl := range f.Decls {
		if reflect.TypeOf(decl).String() == "*ast.FuncDecl" {
			funDecl := decl.(*ast.FuncDecl)
			if funDecl.Name.String() == name {
				found = funDecl
				break
			}
		}
	}
	//ast.Print(fset, found)
	return found
}
