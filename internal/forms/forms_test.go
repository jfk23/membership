package forms

import (
	// "fmt"
	"net/http"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	// r, _ := http.NewRequest("POST", "/someform", nil)
	// form := New(r.PostForm)
	var form Form

	isValid := form.Valid()

	if !isValid {
		t.Error("There is error in the form")
	}

}

func TestForm_New(t *testing.T) {
	r, _ := http.NewRequest("POST", "/someform", nil)

	form := New(r.PostForm)

	if form == nil {
		t.Error("got error when creating new Form")
	}

}

func TestForm_Required(t *testing.T) {
	r, _ := http.NewRequest("POST", "/someform", nil)
	form := New(r.PostForm)

	form.Required("a", "b", "c") 
	v := form.Valid()
	if v {
		t.Error("supposed to be invalid because required fields are missing")
	}

	r, _ = http.NewRequest("POST", "/someform", nil)
	

	postData := url.Values{"a": {"a"}, "b":{"b"}, "c":{"c"}}
	r.PostForm = postData
	form = New(r.PostForm)


	form.Required("a", "b", "c")
	v = form.Valid()
	if !v {
		t.Error("got error even when required fields are provided")
	}
}

func TestForm_MinLength (t *testing.T) {
	r, _ := http.NewRequest("POST", "/someform", nil)
	form := New(r.PostForm)

	form.MinLength("b", 1)
	if form.Valid() {
		t.Error("should have an error for non-existant field, but passed")
	}

	e := form.Error.Get("b")
	if e == "" {
		t.Error("should have an error item, but nothing??")
	}
	
	postData := url.Values{"a": {"a"},}
	r.PostForm = postData
	form = New(r.PostForm)
	// print(form.Get("a"))
	form.MinLength("a", 3)

	v := form.Valid()
	// print(v)
	if v {
		t.Error("should be error because required field length is less than minimum")
	}

	postData = url.Values{"a": {"abcde"},}
	r.PostForm = postData

	form = New(r.PostForm)
	// print(form.Get("a"))

	form.MinLength("a", 3)

	v = form.Valid()
	// print(v)
	if !v {
		t.Error("supposed to be true because required field length is over minimum")
	}

	e = form.Error.Get("a")
	if e != "" {
		t.Error("should have no error item, but one??")
	}

}

func TestForm_EmailValidate (t *testing.T) {
	

	postData := url.Values{"a": {"a"},}
	
	form := New(postData)

	form.EmailValidate("a")
	v := form.Valid()
	if v {
		t.Error("supposed to be error for mal-formatted email entry")
	}

	postData = url.Values{"a": {"a@b.com"},}
	
	form = New(postData)

	form.EmailValidate("a")
	v = form.Valid()
	if !v {
		t.Error("supposed to be error-free for well-formatted email entry")
	}
	
}

func TestForm_Has (t *testing.T) {
	
	postData := url.Values{"a": {""},}
	
	form := New(postData)

	v := form.Has("a")
	
	if v {
		t.Error("supposed to be error for empty required field")
	}

	
	postData = url.Values{"a": {"a"},}
	
	form = New(postData)
	
	
	v1 := form.Has("a")

	if !v1 {
		t.Error("supposed to be true, error-free for filled required field")
	}

	
}



