package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

// use name of TestForm -> because function Valid has a receiver of (f *Form)
func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Error("got invalid when should have been valid")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form shows valid when required fields missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "value")
	postedData.Add("b", "value")
	postedData.Add("c", "value")

	r, _ = http.NewRequest("POST", "/whatever", nil)
	r.PostForm = postedData

	form = New(r.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("shows doesn't have required fields when it does")
	}
}

func TestForm_Has(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	if form.Has("a") {
		t.Error("form shows valid when field missing")
	}

	postedData = url.Values{}
	postedData.Add("a", "value")
	form = New(postedData)

	if !form.Has("a") {
		t.Error("shows doesn't have the fields when it does")
	}
}

func TestForm_MinLength(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	form.MinLength("x", 10)
	if form.Valid() {
		t.Error("form shows min length for non-existent field")
	}

	//unit test for Get function in errors.go
	isError := form.Errors.Get("x")
	if isError == "" {
		t.Error("should have an error, but didn't get one")
	}

	postedData = url.Values{}
	postedData.Add("a", "123")
	form = New(postedData)
	if form.MinLength("a", 4) {
		t.Error("shows minlength of 4 met when data is shorter")
	}

	postedData = url.Values{}
	postedData.Add("b", "123")
	form = New(postedData)
	if !form.MinLength("b", 3) {
		t.Error("shows minlength of 3 is not met when it is")
	}
	//unit test for Get function in errors.go
	isError = form.Errors.Get("b")
	if isError != "" {
		t.Error("shouldn't have an error, but get one")
	}
}

func TestForm_IsEmail(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	form.IsEmail("x")
	if form.Valid() {
		t.Error("form shows valid email for non-existent field")
	}

	postedData = url.Values{}
	postedData.Add("wrongEmail", "123")
	form = New(postedData)
	form.IsEmail("wrongEmail")
	if form.Valid() {
		t.Error("form shows valid when field email is invalid")
	}

	postedData = url.Values{}
	postedData.Add("correctEmail", "email@mail.com")
	form = New(postedData)
	form.IsEmail("correctEmail")
	if !form.Valid() {
		t.Error("form shows invalid when field email is valid")
	}
}
