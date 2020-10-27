package main

import (
	"encoding/json"
	"encoding/xml"
	_ "fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

type dataStructure struct {
	Id        int    `xml:"id"`
	Guid      string `xml:"guid"`
	IsActive  bool   `xml:"isActive"`
	Balance   string `xml:"balance"`
	Picture   string `xml:"picture"`
	Age       int    `xml:"age"`
	EyeColor  string `xml:"eyeColor"`
	FirstName string `xml:"first_name"`
	LastName  string `xml:"last_name"`
	Gender    string `xml:"gender"`
	Company   string `xml:"company"`
	Email     string `xml:"email"`
	Phone     string `xml:"phone"`
	Address   string `xml:"address"`
	About     string `xml:"about"`
}

type xmlStructure struct {
	Version string          `xml:"version"`
	Row     []dataStructure `xml:"row"`
}

const limitSize = 25

func SearchServer(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadFile("dataset.xml")
	if err != nil {
		panic(err)
	}

	usersInXml := &xmlStructure{}
	xml.Unmarshal(data, &usersInXml)

	var users []User

	for _, user := range usersInXml.Row {
		users = append(users, User{
			Id:     user.Id,
			Name:   user.FirstName,
			Age:    user.Age,
			About:  user.About,
			Gender: user.Gender,
		})
	}

	offset, _ := strconv.Atoi(r.FormValue("offset"))
	limit, _ := strconv.Atoi(r.FormValue("limit"))

	var initialRow int
	if offset > 0 {
		initialRow = offset * limitSize
	}

	finalRow := initialRow + limit
	users = users[initialRow:finalRow]

	resultJson, err := json.Marshal(users)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resultJson)
}

func SSLimitMalfunction(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadFile("dataset.xml")
	if err != nil {
		panic(err)
	}

	usersInXml := &xmlStructure{}
	xml.Unmarshal(data, &usersInXml)

	var users []User

	for _, user := range usersInXml.Row {
		users = append(users, User{
			Id:     user.Id,
			Name:   user.FirstName,
			Age:    user.Age,
			About:  user.About,
			Gender: user.Gender,
		})
	}

	resultJson, err := json.Marshal(users)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resultJson)
}

func SSJsonError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, `"err": "bad json"}`)
}

func SSTimeoutError(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Second * 3)
	w.WriteHeader(http.StatusOK)
}

func SSUnknownError(w http.ResponseWriter, r *http.Request) {}

func SSStatusUnauthorizedError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusUnauthorized)
}

func SSInternalServerError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
}

func SSBadRequestError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
}

func SSFieldError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	jsonResponse, _ := json.Marshal(SearchErrorResponse{Error: "ErrorBadOrderField"})
	w.Write(jsonResponse)
}

func SSBadError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	jsonResponse, _ := json.Marshal(SearchErrorResponse{Error: "Unknown error"})
	w.Write(jsonResponse)
}

func TestErrorResponse(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(SearchServer))

	searchClient := &SearchClient{
		URL: ts.URL,
	}

	searchRequest := SearchRequest{
		Limit:  5,
		Offset: 0,
	}

	_, err := searchClient.FindUsers(searchRequest)

	if err != nil {
		t.Error("Dosn't work success request")
	}

	searchRequest.Limit = -1

	_, err = searchClient.FindUsers(searchRequest)
	if err.Error() != "limit must be > 0" {
		t.Error("limit must be > 0")
	}

	searchRequest.Limit = 1
	searchRequest.Offset = -1
	_, err = searchClient.FindUsers(searchRequest)
	if err.Error() != "offset must be > 0" {
		t.Error("offset must be > 0")
	}

	ts.Close()
}

func TestLimitFailed(t *testing.T) {
	limit := 11
	ts := httptest.NewServer(http.HandlerFunc(SSLimitMalfunction))

	searchClient := &SearchClient{
		URL: ts.URL,
	}

	response, _ := searchClient.FindUsers(SearchRequest{Limit: limit})

	if limit == len(response.Users) {
		t.Error("Limit not true")
	}
	ts.Close()
}

func TestBadJson(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SSJsonError))
	searchClient := &SearchClient{
		URL: ts.URL,
	}
	_, err := searchClient.FindUsers(SearchRequest{})

	if err.Error() != `cant unpack result json: invalid character ':' after top-level value` {
		t.Error("Bad json test ")
	}
	ts.Close()
}

func TestOverLimit(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	searchClient := &SearchClient{
		URL: ts.URL,
	}

	response, _ := searchClient.FindUsers(SearchRequest{Limit: 126})

	if 25 != len(response.Users) {
		t.Error("Overlimit has been detected")
	}
	ts.Close()
}

func TestTimeoutError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SSTimeoutError))
	searchClient := &SearchClient{
		URL: ts.URL,
	}

	_, err := searchClient.FindUsers(SearchRequest{})

	if err == nil {
		t.Error("Timeout error detected")
	}

	ts.Close()
}

func TestUnknownError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SSUnknownError))
	searchClient := &SearchClient{
		URL: "bad_link",
	}

	_, err := searchClient.FindUsers(SearchRequest{})

	if err == nil {
		t.Error("Unknown Error has been detected")
	}

	ts.Close()
}

func TestStatusUnauthorized(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SSStatusUnauthorizedError))
	searchClient := &SearchClient{URL: ts.URL}
	_, err := searchClient.FindUsers(SearchRequest{})

	if err.Error() != "Bad AccessToken" {
		t.Error("Access Token error detected")
	}

	ts.Close()
}

func TestStatusInternalServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SSInternalServerError))
	searchClient := &SearchClient{URL: ts.URL}
	_, err := searchClient.FindUsers(SearchRequest{})

	if err.Error() != "SearchServer fatal error" {
		t.Error("Internal Server error detected")
	}

	ts.Close()
}

func TestBadRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SSBadRequestError))
	searchClient := &SearchClient{URL: ts.URL}
	_, err := searchClient.FindUsers(SearchRequest{})

	if err.Error() != "cant unpack error json: unexpected end of JSON input" {
		t.Error("Bad Request error detected")
	}

	ts.Close()
}

func TestBadField(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SSFieldError))
	searchClient := &SearchClient{URL: ts.URL}
	_, err := searchClient.FindUsers(SearchRequest{})
	if err.Error() != "OrderFeld  invalid" {
		t.Error("Field error detected")
	}

	ts.Close()
}

func TestBadRequestError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SSBadError))
	searchClient := &SearchClient{URL: ts.URL}
	_, err := searchClient.FindUsers(SearchRequest{})
	if err == nil {
		t.Error("Bad request error detected")
	}

	ts.Close()
}
