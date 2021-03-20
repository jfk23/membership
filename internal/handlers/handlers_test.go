package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"

	// "net/url"
	"strings"

	// "net/url"

	// "net/url"
	"testing"

	"github.com/jfk23/gobookings/internal/model"
)

type postData struct {
	key string
	value string
}

var theTests = []struct{
	name string
	url string
	method string
	// params []postData
	expectedCode int
} {
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"generalsroom", "/generals-room", "GET", http.StatusOK},
	{"majorsroom", "/majors-room", "GET", http.StatusOK},
	{"makereservation", "/make-reservation", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},
	{"search", "/search", "GET", http.StatusOK},
	{"non-exist", "/blue/sky", "GET", http.StatusNotFound},
	{"dashboard", "/admin/dashboard", "GET", http.StatusOK},
	{"res all", "/admin/reservations-all", "GET", http.StatusOK},
	{"res new", "/admin/reservations-new", "GET", http.StatusOK},
	{"res show", "/admin/reservation/new/1/show", "GET", http.StatusOK},
	{"res cal", "/admin/reservations-calendar", "GET", http.StatusOK},
	
	// {"make-reserve-post", "/make-reservation", "POST", []postData{
	// 	{"start", "01-02-2021"},
	// 	{"end", "01-10-2021"},
	// }, http.StatusOK},

	// {"search-post", "/search", "POST", []postData{
	// 	{"start", "01-02-2021"},
	// 	{"end", "01-10-2021"},
	// }, http.StatusOK},

	// {"search-post-json", "/search-json", "POST", []postData{
	// 		{"first_name" , "Jay"},
	// 		{"last_name" , "Choi"},
	// 		{"email" , "me@here.com"},
	// 		{"phone" , "123-333-4444"},
	// },
	// http.StatusOK},
}
	

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, test := range theTests {
		if test.method == "GET" {
			resp, err := ts.Client().Get(ts.URL+test.url)
			if err != nil {
				t.Error(err)
				t.Fatal(err)
			}
			if resp.StatusCode != test.expectedCode {
				t.Errorf("for %s, expected outcome is %d, but got %d instead", test.name, test.expectedCode, resp.StatusCode)
			} 
		} 
			
	}

}	 

func Test_MakeReseration(t *testing.T) {
	reservation := model.Reservation{
		RoomID : 1,
		Room : model.Room{
			ID : 1,
			RoomName: "General's Quaters",
		},
	}
	
	req, _ := http.NewRequest("GET", "/make-reservation", nil)
	ctx := getCTX(req)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	sessionManager.Put(ctx, "reservation", reservation)


	handle := http.HandlerFunc(Repo.Reservation)

	handle.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("we got wrong code, expected %d, but got %d", http.StatusOK, rr.Code)
	}

	// testing with existing session data (reservation)

	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCTX(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()

	handle = http.HandlerFunc(Repo.Reservation)

	handle.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("we got wrong code, expected %d, but got %d", http.StatusTemporaryRedirect, rr.Code)
	}

	// testing with non-existent room search

	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCTX(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()

	reservation.RoomID = 3

	sessionManager.Put(ctx, "reservation", reservation)

	handle = http.HandlerFunc(Repo.Reservation)

	handle.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("we got wrong code, expected %d, but got %d", http.StatusTemporaryRedirect, rr.Code)
	}

}

func Test_PostReservation(t *testing.T) {
	reqBody := "start_date=2022-01-05"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2022-01-07")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=Johnyl")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smithon")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=smith@mem.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=2221112222")

	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx := getCTX(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type","application/x-www-form-urlencoded")
	
	rr := httptest.NewRecorder()

	handle := http.HandlerFunc(Repo.PostReservation)

	handle.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("we got wrong code, expected %d, but got %d", http.StatusSeeOther, rr.Code)
	}




	//testing invalid form data - non existant id
	reqBody = "start_date=2050-12-10"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-12-11")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=Chris")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=2")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=smith@mem.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=44444")

	//print(reqBody)
	
	// postBody := url.Values{}
	// postBody.Add("start_date", "2050-01-25")
	// postBody.Add("end_date", "2050-01-29")
	// postBody.Add("first_name", "Jonh")
	// postBody.Add("last_name", "Smith")
	// postBody.Add("email", "Jonh@smith.com")
	// postBody.Add("phone", "1112223333") 
	// postBody.Add("room_id", "1000")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCTX(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type","application/x-www-form-urlencoded")
	
	rr = httptest.NewRecorder()

	handle = http.HandlerFunc(Repo.PostReservation)

	handle.ServeHTTP(rr, req)
	//print(rr.Code)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("we got wrong code for failing to insert reservation data(wrong room id), expected %d, but got %d", http.StatusTemporaryRedirect, rr.Code)
	}

	// testing no form body
	req, _ = http.NewRequest("POST", "/make-reservation", nil)
	ctx = getCTX(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type","application/x-www-form-urlencoded")
	
	rr = httptest.NewRecorder()

	handle = http.HandlerFunc(Repo.PostReservation)

	handle.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("we got wrong code for no form body, expected %d, but got %d", http.StatusTemporaryRedirect, rr.Code)
	}

	// testing cannot convert start time
	reqBody = "start_date=invalid"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-10")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=Johny")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=smith@mem.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=2221112222")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCTX(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type","application/x-www-form-urlencoded")
	
	rr = httptest.NewRecorder()

	handle = http.HandlerFunc(Repo.PostReservation)

	handle.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("we got wrong code for invalid time, expected %d, but got %d", http.StatusTemporaryRedirect, rr.Code)
	}

	// testing cannot convert end time
	reqBody = "start_date=2050-01-10"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=invalid")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=Johny")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=smith@mem.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=2221112222")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCTX(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type","application/x-www-form-urlencoded")
	
	rr = httptest.NewRecorder()

	handle = http.HandlerFunc(Repo.PostReservation)

	handle.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("we got wrong code for invalid time, expected %d, but got %d", http.StatusTemporaryRedirect, rr.Code)
	}

	// testing invalid room_id
	reqBody = "start_date=2050-01-10"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-11")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=Johny")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=invalid")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=smith@mem.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=2221112222")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCTX(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type","application/x-www-form-urlencoded")
	
	rr = httptest.NewRecorder()

	handle = http.HandlerFunc(Repo.PostReservation)

	handle.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("we got wrong code for invalid room id, expected %d, but got %d", http.StatusTemporaryRedirect, rr.Code)
	}

	// testing invalid form data - too short first_name
	reqBody = "start_date=2050-01-10"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-11")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=Jo")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=smith@mem.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=2221112222")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCTX(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type","application/x-www-form-urlencoded")
	
	rr = httptest.NewRecorder()

	handle = http.HandlerFunc(Repo.PostReservation)

	handle.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("we got wrong code for invalid form data(first name), expected %d, but got %d", http.StatusSeeOther, rr.Code)
	}

	// reqBody := "start_date=2050-01-01"
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-10")
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=Johny")
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "email=smith@mem.com")
	// reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=2221112222")

	// testing invalid form data - no first name
	reqBody = "start_date=2050-11-10"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-11-11")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=smith@mem.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=2221112222")

	

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCTX(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type","application/x-www-form-urlencoded")
	
	rr = httptest.NewRecorder()

	handle = http.HandlerFunc(Repo.PostReservation)

	handle.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("we got wrong code for failing to insert reservation data(no first_name), expected %d, but got %d", http.StatusTemporaryRedirect, rr.Code)
	}

	
}
	
	

func Test_ReservationSummary(t *testing.T) {
	
	req, _ := http.NewRequest("GET", "/reservation-summary", nil)
	ctx := getCTX(req)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()


	handle := http.HandlerFunc(Repo.ReservationSummary)

	handle.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("we got wrong code, expected %d, but got %d", http.StatusTemporaryRedirect, rr.Code)
	}


}

func Test_PostSearchJson(t *testing.T) {

	reqBody := "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-10")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, _ := http.NewRequest("POST", "/search-json", strings.NewReader(reqBody))
	ctx := getCTX(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type","application/x-www-form-urlencoded")
	
	rr := httptest.NewRecorder()

	handle := http.HandlerFunc(Repo.PostSearchJson)

	handle.ServeHTTP(rr, req)

	var r responseJson

	err := json.Unmarshal([]byte(rr.Body.Bytes()), &r)

	if err != nil {
		t.Error("there is something wrong with correct input data")
	}

}

var loginTests = []struct {
	name string
	email string
	password string
	expectedStatusCode int
	expectedHTML string
	expectedLocation string
}{
	 {
		name : "valid credentials",
		email : "admin@a.a",
		password : "admin1",
		expectedStatusCode: http.StatusSeeOther,
		expectedHTML: "",
		expectedLocation: "/",
	 },

	 {
		name : "invalid credentials",
		email : "jack@a.a",
		password : "admin1",
		expectedStatusCode: http.StatusSeeOther,
		expectedHTML: "",
		expectedLocation: "/user/login",
	 },
	 {
		name : "invalid form",
		email : "admin",
		password : "admin1",
		expectedStatusCode: http.StatusOK,
		expectedHTML: `action="/user/login"`,
		expectedLocation: "",
	 },
	 
}

func Test_ShowLogin(t *testing.T) {
	for _, e := range loginTests {
		postData := url.Values{}
		postData.Add("email", e.email)
		postData.Add("password", e.password)

		req, _ := http.NewRequest("POST", "/user/login", strings.NewReader(postData.Encode()))
		ctx := getCTX(req)
		req = req.WithContext(ctx)
		req.Header.Set("Content-Type","application/x-www-form-urlencoded")
		
		rr := httptest.NewRecorder()

		handle := http.HandlerFunc(Repo.PostShowLogin)

		handle.ServeHTTP(rr, req)

		//fmt.Println(rr.Result().Location())

		if rr.Code != e.expectedStatusCode {
			t.Errorf("we got wrong code, expected %d, but got %d", e.expectedStatusCode, rr.Code)
		}

		if e.expectedLocation != "" {
			resultLoc, _ := rr.Result().Location()
			//log.Println(resultLoc.String())

			if resultLoc.String() != e.expectedLocation {
				t.Errorf("for the test of %s, we expected to land on %s, but instead %s", e.name, e.expectedLocation, resultLoc.String())
			}
		}

		if e.expectedHTML != "" {
			html := rr.Body.String()

			if !strings.Contains(html, e.expectedHTML) {
				t.Errorf("for the test of %s, we expected to have this html %s, but didn't have", e.name, e.expectedHTML)
			}

		}

	}
}

func getCTX(req *http.Request) context.Context{
	ctx, _ := sessionManager.Load(req.Context(), req.Header.Get("X-Session"))
	return ctx

}