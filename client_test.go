package scgi

import (
	"testing"
)

func TestCreateClient(t *testing.T) {
	cl := NewClient("network", "address")

	if cl.network != "network" || cl.address != "address" {
		t.Fatalf("Well that was unexpected")
	}

	/*
	  if len(cl.headers) != 1 || cl.headers["SCGI"] != "1" {
	    t.Fatalf("Missing/invalid default header")
	  }
	*/
}

func TestSetHeader(t *testing.T) {
	/*
	  cl := NewClient("network", "address")

	  cl.SetHeader("bob", "jim")
	  cl.SetHeader("12345", "-1")
	  cl.SetHeader("!@##@!$#@$\n\r", "...")

	  if (len(cl.headers) != 4) {
	    t.Fatalf("Expected 4 headers, got %d", len(cl.headers))
	  }

	  if cl.headers["SCGI"] != "1" ||
	    cl.headers["bob"] != "jim" ||
	    cl.headers["12345"] != "-1" ||
	    cl.headers["!@##@!$#@$\n\r"] != "..." {
	    t.Fatalf("Expected 4 headers set, one of more invalid")
	  }
	*/
}

func TestBuildRequest(t *testing.T) {
	expected := "70:CONTENT_LENGTH\x0027\x00SCGI\x001\x00REQUEST_METHOD\x00POST\x00REQUEST_URI\x00/deepthought\x00,What is the answer to life?"

	cl := NewClient("network", "address")
	cl.SetHeader("REQUEST_METHOD", "POST")
	cl.SetHeader("REQUEST_URI", "/deepthought")

	body := []byte("What is the answer to life?")
	header := cl.makeHeaders(len(body))
	actual := cl.netstring(header, body)

	if string(actual) != string(expected) {
		t.Fatalf("Expected %s, got %s\n", expected, actual)
	}
}

func TestReuseRequest(t *testing.T) {
	var header, body, actual []byte

	expected1 := "70:CONTENT_LENGTH\x0027\x00SCGI\x001\x00REQUEST_METHOD\x00POST\x00REQUEST_URI\x00/deepthought\x00,What is the answer to life?"
	expected2 := "72:CONTENT_LENGTH\x0051\x00SCGI\x001\x00REQUEST_METHOD\x00POST\x00REQUEST_URI\x00/morethoughts/\x00,The answer is obviously 42...to another question..."

	cl := NewClient("network", "address")

	//first request
	cl.SetHeader("REQUEST_METHOD", "POST")
	cl.SetHeader("REQUEST_URI", "/deepthought")

	body = []byte("What is the answer to life?")
	header = cl.makeHeaders(len(body))
	actual = cl.netstring(header, body)

	if string(actual) != string(expected1) {
		t.Fatalf("Expected %s, got %s\n", expected1, actual)
	}

	//second request
	cl.SetHeader("REQUEST_URI", "/morethoughts/")

	body = []byte("The answer is obviously 42...to another question...")
	header = cl.makeHeaders(len(body))
	actual = cl.netstring(header, body)

	if string(actual) != string(expected2) {
		t.Fatalf("Expected %s, got %s\n", expected2, actual)
	}
}
