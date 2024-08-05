package forms

import (
	"log"
	"net/url"
	"testing"
)

var form = New(url.Values{})

func TestNew(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	if len(form.Errors) == 0 {
		log.Print("error map is empty")
	} else {
		t.Error("new map should not have errors")
	}
}

func TestValid(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	if !form.Valid() {
		t.Error("form is not valid")
	}
}

func TestRequired(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	form.Required("a")
	if form.Valid() {
		t.Error("form false valid, required some values")
	}

	postedData = url.Values{}
	postedData.Add("a", "some value")
	form = New(postedData)
	form.Required("a")

	if !form.Valid() {
		t.Error("form false invalid")
	}
}

func TestForm_Has(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	has := form.Has("whatever")
	if has {
		t.Error("form show has field when empty")
	}

	postedData = url.Values{}
	postedData.Add("a", "some value")

	form = New(postedData)

	has = form.Has("a")
	if !has {
		t.Error("false !has")
	}
}

func TestForm_MinLength(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	form.MinLength("x", 10)
	if form.Valid() {
		t.Error("false valid min length")
	}

	isError := form.Errors.Get("x")
	if isError == "" {
		t.Error("error is not added into stack")
	}

	postedData = url.Values{}
	postedData.Add("a", "12345")
	form = New(postedData)
	form.MinLength("a", 6)

	if form.Valid() {
		t.Error("false valid for insufficient length")
	}

	postedData = url.Values{}
	postedData.Add("a", "12345")
	form = New(postedData)
	form.MinLength("a", 5)
	if !form.Valid() {
		t.Error("false valid for insufficient length")
	}

	isError = form.Errors.Get("a")
	if isError != "" {
		t.Error("should not have error, but got one")
	}
}

func TestForm_IsEmail(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)
	form.IsEmail("a")
	if form.Valid() {
		t.Error("false valid [non-existent field]")
	}

	postedData = url.Values{}
	postedData.Add("email", "asdn@gmail.com")
	form = New(postedData)
	form.IsEmail("email")
	if !form.Valid() {
		t.Error("false invalid [email is valid]")
	}

	postedData = url.Values{}
	postedData.Add("email", "asdngmail.com")
	form = New(postedData)
	form.IsEmail("email")
	if form.Valid() {
		t.Error("false valid [email is invalid]")
	}

}
