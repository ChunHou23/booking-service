package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/ChunHou23/booking-service/internal/models"
)

var theTests = []struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},
	{"js", "/junior-suite", "GET", http.StatusOK},
	{"gs", "/general-suite", "GET", http.StatusOK},
	{"search-a", "/search-availability", "GET", http.StatusOK},
	// {"make-re", "/make-reservation", "GET", http.StatusOK},
	// {"re-summary", "/reservation-summary", "GET", http.StatusOK},

	// {"post-search-a", "/search-availability", "POST", []postData{
	// 	{key: "start", value: "2020-01-01"},
	// 	{key: "end", value: "2020-01-02"},
	// }, http.StatusOK},
	// {"post-search-a-json", "/search-availability-json", "POST", []postData{
	// 	{key: "start", value: "2020-01-01"},
	// 	{key: "end", value: "2020-01-02"},
	// }, http.StatusOK},
	{"post-make-reservation", "/make-reservation", "POST", http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {
		if e.method == "GET" {
			resp, err := ts.Client().Get(ts.URL + e.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("For %s, expected code: %d, but code: %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}

		}
	}
}

func TestRepository_Reservation(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Suite",
		},
	}

	req, _ := http.NewRequest("GET", "/make-reservation", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.MakeReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("reservation handler return wrong response : %d, wanted: %d", rr.Code, http.StatusOK)
	}

	// test when no reservation in session
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("reservation handler return wrong response : %d, wanted: %d", rr.Code, http.StatusOK)
	}

	// test with non existent room
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()
	reservation.RoomID = 100
	reservation.Room.ID = 100
	session.Put(ctx, "reservation", reservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("reservation handler return wrong response : %d, wanted: %d", rr.Code, http.StatusOK)
	}
}

func TestRepository_PostReservation(t *testing.T) {
	reqBody := "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=test@gmail.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone_number=123")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", models.Reservation{})

	handler := http.HandlerFunc(Repo.PostMakeReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("post reservation handler return wrong response : %d, wanted: %d", rr.Code, http.StatusSeeOther)
	}
}

func TestRepository_PostReservation_MissingPostBody(t *testing.T) {
	req, _ := http.NewRequest("POST", "/make-reservation", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", models.Reservation{})

	handler := http.HandlerFunc(Repo.PostMakeReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("post reservation missing form handler return wrong response : %d, wanted: %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_PostReservation_InvalidRoomId(t *testing.T) {
	// test invalid room id
	reqBody := "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=James")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=test@gmail.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone_number=123")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=invalid")
	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))

	ctx := getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.PostMakeReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("post reservation missing form handler return wrong response : %d, wanted: %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_PostReservation_InvalidForm(t *testing.T) {
	// test invalid room id
	reqBody := "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=J")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=test@gmail.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone_number=123")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=invalid")
	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))

	ctx := getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", models.Reservation{})

	handler := http.HandlerFunc(Repo.PostMakeReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("post reservation missing form handler return wrong response : %d, wanted: %d", rr.Code, http.StatusSeeOther)
	}
}

func TestRepository_PostReservation_InvalidInsertReservation(t *testing.T) {
	reqBody := "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=James")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=test@gmail.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone_number=123")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=2")
	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))

	ctx := getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", models.Reservation{
		RoomID: 2,
	})

	handler := http.HandlerFunc(Repo.PostMakeReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("post reservation missing form handler return wrong response : %d, wanted: %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_PostReservation_InvalidInsertRestriction(t *testing.T) {
	reqBody := "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=James")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=test@gmail.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone_number=123")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1000")
	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))

	ctx := getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", models.Reservation{
		RoomID: 1000,
	})

	handler := http.HandlerFunc(Repo.PostMakeReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("post reservation missing form handler return wrong response : %d, wanted: %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_PostAvailabilityHaveRoom(t *testing.T) {
	reqBody := "start=2100-01-01&end=2100-01-03"
	req, _ := http.NewRequest("POST", "/search-availability", strings.NewReader(reqBody))
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.PostAvailability)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("post reservation missing form handler return wrong response : %d, wanted: %d", rr.Code, http.StatusOK)
	}
}

func TestRepository_PostAvailabilityNoRoom(t *testing.T) {
	reqBody := "start=2050-01-01&end=2050-01-03"
	req, _ := http.NewRequest("POST", "/search-availability", strings.NewReader(reqBody))
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.PostAvailability)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("post reservation missing form handler return wrong response : %d, wanted: %d", rr.Code, http.StatusSeeOther)
	}
}

func TestRepository_PostAvailabilityInvalidForm(t *testing.T) {
	req, _ := http.NewRequest("POST", "/search-availability", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.PostAvailability)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("post reservation missing form handler return wrong response : %d, wanted: %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_PostAvailabilityInvalidStartDate(t *testing.T) {
	reqBody := "start=2050-01-031&end=2050-01-03"
	req, _ := http.NewRequest("POST", "/search-availability", strings.NewReader(reqBody))
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.PostAvailability)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("post reservation missing form handler return wrong response : %d, wanted: %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_PostAvailabilityInvalidEndDate(t *testing.T) {
	reqBody := "start=2050-01-31&end=2050-01-63"
	req, _ := http.NewRequest("POST", "/search-availability", strings.NewReader(reqBody))
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.PostAvailability)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("post reservation missing form handler return wrong response : %d, wanted: %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_PostAvailabilityInvalidSearchDB(t *testing.T) {
	reqBody := "start=2090-01-01&end=2090-01-03"
	req, _ := http.NewRequest("POST", "/search-availability", strings.NewReader(reqBody))
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.PostAvailability)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("post reservation missing form handler return wrong response : %d, wanted: %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_PostAvailabilityJSON(t *testing.T) {
	reqBody := "start_date=2050-01-01&end_date=2050-01-03&room_id=1"
	req, _ := http.NewRequest("POST", "/search-availability-json", strings.NewReader(reqBody))
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.AvailabilityJSON)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("post reservation missing form handler return wrong response : %d, wanted: %d", rr.Code, http.StatusOK)
	}
}

func TestRepository_PostAvailabilityJSONInvalidForm(t *testing.T) {
	req, _ := http.NewRequest("POST", "/search-availability-json", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.AvailabilityJSON)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("post reservation missing form handler return wrong response : %d, wanted: %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_PostAvailabilityJSONInvalidStartDate(t *testing.T) {
	reqBody := "start_date=2050-01-41&end_date=2050-01-03"
	req, _ := http.NewRequest("POST", "/search-availability-json", strings.NewReader(reqBody))
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.AvailabilityJSON)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("post reservation missing form handler return wrong response : %d, wanted: %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_PostAvailabilityJSONInvalidEndDate(t *testing.T) {
	reqBody := "start_date=2050-01-01&end_date=2050-01-53"
	req, _ := http.NewRequest("POST", "/search-availability-json", strings.NewReader(reqBody))
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.AvailabilityJSON)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("post reservation missing form handler return wrong response : %d, wanted: %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_PostAvailabilityJSONInvalidRoomId(t *testing.T) {
	reqBody := "start_date=2050-01-01&end_date=2050-01-03&room_id=invalid"
	req, _ := http.NewRequest("POST", "/search-availability-json", strings.NewReader(reqBody))
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.AvailabilityJSON)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("post reservation missing form handler return wrong response : %d, wanted: %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_PostAvailabilityJSONInvalidSearchRoomById(t *testing.T) {
	reqBody := "start_date=2050-01-01&end_date=2050-01-03&room_id=1000"
	req, _ := http.NewRequest("POST", "/search-availability-json", strings.NewReader(reqBody))
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.AvailabilityJSON)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("post reservation missing form handler return wrong response : %d, wanted: %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_ReservationSummary(t *testing.T) {
	sd, _ := time.Parse(time.DateOnly, "2050-01-01")
	ed, _ := time.Parse(time.DateOnly, "2050-01-03")
	reservation := models.Reservation{
		StartDate: sd,
		EndDate:   ed,
	}
	req, _ := http.NewRequest("Get", "/reservation-summary", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	session.Put(ctx, "reservation", reservation)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.ReservationSummary)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("post reservation missing form handler return wrong response : %d, wanted: %d", rr.Code, http.StatusOK)
	}
}

func TestRepository_ReservationSummaryNoSession(t *testing.T) {
	req, _ := http.NewRequest("Get", "/reservation-summary", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.ReservationSummary)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("post reservation missing form handler return wrong response : %d, wanted: %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_ChooseRoom(t *testing.T) {
	/*****************************************
	// first case -- reservation in session
	*****************************************/
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}

	req, _ := http.NewRequest("GET", "/choose-room/1", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	// set the RequestURI on the request so that we can grab the ID
	// from the URL
	req.RequestURI = "/choose-room/1"

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.ChooseRoom)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("ChooseRoom handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	///*****************************************
	//// second case -- reservation not in session
	//*****************************************/
	req, _ = http.NewRequest("GET", "/choose-room/1", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.RequestURI = "/choose-room/1"

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.ChooseRoom)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("ChooseRoom handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	///*****************************************
	//// third case -- missing url parameter, or malformed parameter
	//*****************************************/
	req, _ = http.NewRequest("GET", "/choose-room/fish", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.RequestURI = "/choose-room/fish"

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.ChooseRoom)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("ChooseRoom handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_BookRoom(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}

	req, _ := http.NewRequest("GET", "/book-room/1", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	// set the RequestURI on the request so that we can grab the ID
	// from the URL
	req.RequestURI = "/book-room/1"

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.BookRoom)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("BookRoom handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}
}

func TestRepository_BookRoomInvalidParam(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}

	req, _ := http.NewRequest("GET", "/book-room/test", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	// set the RequestURI on the request so that we can grab the ID
	// from the URL
	req.RequestURI = "/book-room/test"

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.BookRoom)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("BookRoom handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_BookRoomInvalidRoomId(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 3,
		Room: models.Room{
			ID:       3,
			RoomName: "General's Quarters",
		},
	}

	req, _ := http.NewRequest("GET", "/book-room/3", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	// set the RequestURI on the request so that we can grab the ID
	// from the URL
	req.RequestURI = "/book-room/3"

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.BookRoom)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("BookRoom handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

// loginTests is the data for the Login handler tests
var loginTests = []struct {
	name               string
	email              string
	expectedStatusCode int
	expectedHTML       string
	expectedLocation   string
}{
	{
		"valid-credentials",
		"me@here.ca",
		http.StatusSeeOther,
		"",
		"/",
	},
	{
		"invalid-credentials",
		"jack@nimble.com",
		http.StatusSeeOther,
		"",
		"/user/login",
	},
	{
		"invalid-data",
		"j",
		http.StatusOK,
		`action="/user/login"`,
		"",
	},
}

// func TestLogin(t *testing.T) {
// 	// range through all tests
// 	for _, e := range loginTests {
// 		postedData := url.Values{}
// 		postedData.Add("email", e.email)
// 		postedData.Add("password", "password")

// 		// create request
// 		req, _ := http.NewRequest("POST", "/user/login", strings.NewReader(postedData.Encode()))
// 		ctx := getCtx(req)
// 		req = req.WithContext(ctx)

// 		// set the header
// 		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
// 		rr := httptest.NewRecorder()

// 		// call the handler
// 		handler := http.HandlerFunc(Repo.PostShowLogin)
// 		handler.ServeHTTP(rr, req)

// 		if rr.Code != e.expectedStatusCode {
// 			t.Errorf("failed %s: expected code %d, but got %d", e.name, e.expectedStatusCode, rr.Code)
// 		}

// 		if e.expectedLocation != "" {
// 			// get the URL from test
// 			actualLoc, _ := rr.Result().Location()
// 			if actualLoc.String() != e.expectedLocation {
// 				t.Errorf("failed %s: expected location %s, but got location %s", e.name, e.expectedLocation, actualLoc.String())
// 			}
// 		}

// 		// checking for expected values in HTML
// 		if e.expectedHTML != "" {
// 			// read the response body into a string
// 			html := rr.Body.String()
// 			if !strings.Contains(html, e.expectedHTML) {
// 				t.Errorf("failed %s: expected to find %s but did not", e.name, e.expectedHTML)
// 			}
// 		}
// 	}
// }

var adminPostShowReservationTests = []struct {
	name                 string
	url                  string
	postedData           url.Values
	expectedResponseCode int
	expectedLocation     string
	expectedHTML         string
}{
	{
		name: "valid-data-from-new",
		url:  "/admin/reservations/new/1/show",
		postedData: url.Values{
			"first_name": {"John"},
			"last_name":  {"Smith"},
			"email":      {"john@smith.com"},
			"phone":      {"555-555-5555"},
		},
		expectedResponseCode: http.StatusSeeOther,
		expectedLocation:     "/admin/reservations-new",
		expectedHTML:         "",
	},
	{
		name: "valid-data-from-all",
		url:  "/admin/reservations/all/1/show",
		postedData: url.Values{
			"first_name": {"John"},
			"last_name":  {"Smith"},
			"email":      {"john@smith.com"},
			"phone":      {"555-555-5555"},
		},
		expectedResponseCode: http.StatusSeeOther,
		expectedLocation:     "/admin/reservations-all",
		expectedHTML:         "",
	},
	{
		name: "valid-data-from-cal",
		url:  "/admin/reservations/cal/1/show",
		postedData: url.Values{
			"first_name": {"John"},
			"last_name":  {"Smith"},
			"email":      {"john@smith.com"},
			"phone":      {"555-555-5555"},
			"year":       {"2022"},
			"month":      {"01"},
		},
		expectedResponseCode: http.StatusSeeOther,
		expectedLocation:     "/admin/reservations-calendar?y=2022&m=01",
		expectedHTML:         "",
	},
}

// TestAdminPostShowReservation tests the AdminPostReservation handler
func TestAdminPostShowReservation(t *testing.T) {
	for _, e := range adminPostShowReservationTests {
		var req *http.Request
		if e.postedData != nil {
			req, _ = http.NewRequest("POST", "/user/login", strings.NewReader(e.postedData.Encode()))
		} else {
			req, _ = http.NewRequest("POST", "/user/login", nil)
		}
		ctx := getCtx(req)
		req = req.WithContext(ctx)
		req.RequestURI = e.url

		// set the header
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		// call the handler
		handler := http.HandlerFunc(Repo.AdminPostShowReservations)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedResponseCode {
			t.Errorf("failed %s: expected code %d, but got %d", e.name, e.expectedResponseCode, rr.Code)
		}

		if e.expectedLocation != "" {
			// get the URL from test
			actualLoc, _ := rr.Result().Location()
			if actualLoc.String() != e.expectedLocation {
				t.Errorf("failed %s: expected location %s, but got location %s", e.name, e.expectedLocation, actualLoc.String())
			}
		}

		// checking for expected values in HTML
		if e.expectedHTML != "" {
			// read the response body into a string
			html := rr.Body.String()
			if !strings.Contains(html, e.expectedHTML) {
				t.Errorf("failed %s: expected to find %s but did not", e.name, e.expectedHTML)
			}
		}
	}
}

var adminPostReservationCalendarTests = []struct {
	name                 string
	postedData           url.Values
	expectedResponseCode int
	expectedLocation     string
	expectedHTML         string
	blocks               int
	reservations         int
}{
	{
		name: "cal",
		postedData: url.Values{
			"year":  {time.Now().Format("2006")},
			"month": {time.Now().Format("01")},
			fmt.Sprintf("add_block_1_%s", time.Now().AddDate(0, 0, 2).Format("2006-01-2")): {"1"},
		},
		expectedResponseCode: http.StatusSeeOther,
	},
	{
		name:                 "cal-blocks",
		postedData:           url.Values{},
		expectedResponseCode: http.StatusSeeOther,
		blocks:               1,
	},
	{
		name:                 "cal-res",
		postedData:           url.Values{},
		expectedResponseCode: http.StatusSeeOther,
		reservations:         1,
	},
}

func TestPostReservationCalendar(t *testing.T) {
	for _, e := range adminPostReservationCalendarTests {
		var req *http.Request
		if e.postedData != nil {
			req, _ = http.NewRequest("POST", "/admin/reservations-calendar", strings.NewReader(e.postedData.Encode()))
		} else {
			req, _ = http.NewRequest("POST", "/admin/reservations-calendar", nil)
		}
		ctx := getCtx(req)
		req = req.WithContext(ctx)

		now := time.Now()
		bm := make(map[string]int)
		rm := make(map[string]int)

		currentYear, currentMonth, _ := now.Date()
		currentLocation := now.Location()

		firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
		lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

		for d := firstOfMonth; d.After(lastOfMonth) == false; d = d.AddDate(0, 0, 1) {
			rm[d.Format("2006-01-2")] = 0
			bm[d.Format("2006-01-2")] = 0
		}

		if e.blocks > 0 {
			bm[firstOfMonth.Format("2006-01-2")] = e.blocks
		}

		if e.reservations > 0 {
			rm[lastOfMonth.Format("2006-01-2")] = e.reservations
		}

		session.Put(ctx, "block_map_1", bm)
		session.Put(ctx, "reservation_map_1", rm)

		// set the header
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		// call the handler
		handler := http.HandlerFunc(Repo.AdminPostReservationsCalendar)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedResponseCode {
			t.Errorf("failed %s: expected code %d, but got %d", e.name, e.expectedResponseCode, rr.Code)
		}

	}
}

var adminProcessReservationTests = []struct {
	name                 string
	queryParams          string
	expectedResponseCode int
	expectedLocation     string
}{
	{
		name:                 "process-reservation",
		queryParams:          "",
		expectedResponseCode: http.StatusSeeOther,
		expectedLocation:     "",
	},
	{
		name:                 "process-reservation-back-to-cal",
		queryParams:          "?y=2021&m=12",
		expectedResponseCode: http.StatusSeeOther,
		expectedLocation:     "",
	},
}

func TestAdminProcessReservation(t *testing.T) {
	for _, e := range adminProcessReservationTests {
		req, _ := http.NewRequest("GET", fmt.Sprintf("/admin/process-reservation/cal/1/do%s", e.queryParams), nil)
		ctx := getCtx(req)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(Repo.AdminProcessReservation)
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusSeeOther {
			t.Errorf("failed %s: expected code %d, but got %d", e.name, e.expectedResponseCode, rr.Code)
		}
	}
}

var adminDeleteReservationTests = []struct {
	name                 string
	queryParams          string
	expectedResponseCode int
	expectedLocation     string
}{
	{
		name:                 "delete-reservation",
		queryParams:          "",
		expectedResponseCode: http.StatusSeeOther,
		expectedLocation:     "",
	},
	{
		name:                 "delete-reservation-back-to-cal",
		queryParams:          "?y=2021&m=12",
		expectedResponseCode: http.StatusSeeOther,
		expectedLocation:     "",
	},
}

func TestAdminDeleteReservation(t *testing.T) {
	for _, e := range adminDeleteReservationTests {
		req, _ := http.NewRequest("GET", fmt.Sprintf("/admin/process-reservation/cal/1/do%s", e.queryParams), nil)
		ctx := getCtx(req)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(Repo.AdminDeleteReservation)
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusSeeOther {
			t.Errorf("failed %s: expected code %d, but got %d", e.name, e.expectedResponseCode, rr.Code)
		}
	}
}

// gets the context
func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}
