package hapi

import (
    "testing"
    "net/http"
//    "net/http/httptest"
)

type mockResponseWriter struct{}

func (m *mockResponseWriter) Header() (h http.Header) {
    return http.Header{}
}

func (m *mockResponseWriter) Write(p []byte) (n int, err error) {
    return len(p), nil
}

func (m *mockResponseWriter) WriteString(s string) (n int, err error) {
    return len(s), nil
}

func (m *mockResponseWriter) WriteHeader(int) {}


func TestNew(t *testing.T) {
    h := New()
    if h.Router == nil {
        t.Fatal("h.Router not initilized")
    }
}

type typeTest struct{
    Accept      string
    Handlers    map[string]Handle   
    Result      string
    HandlerID   string
}

func NewTest(accept,result,id string) *typeTest {
    return &typeTest{
        accept,
        make(map[string]Handle),
        result,
        id,
    }
}

func (testData *typeTest) AddHandler(ctype,id string) {
    testData.Handlers[ctype] = func(c *Context) { c.Stash["id"] = id }
}

func (testData *typeTest) DoTests(t *testing.T) {
    negType,negHandler := TypeNegotiator(testData.Accept,testData.Handlers)
    if negType != testData.Result {
        t.Fatalf("TypeNegotiator returned '%s', expected '%s'\n", negType, testData.Result)
    }
    context := &Context{}
    context.Stash = make(map[string]interface{})
    if negHandler != nil {
        negHandler( context )
    }
    var id string
    if context.Stash["id"] != nil {
        id = context.Stash["id"].(string)
    }
    if id != testData.HandlerID {
        t.Fatalf("Handler identified itself as '%s', expected '%x'\n", id, testData.HandlerID)
    }
}

// Test a handler for 'text/html' only
func TestTypeNegotiator(t *testing.T) {
    // Accept: matches exactly
    testData := NewTest("text/html","text/html","1")
    testData.AddHandler("text/html","1")
    testData.DoTests(t)

    // Accept: Doesn't match
    testData = NewTest("text/plain","","")
    testData.AddHandler("text/html","1")
    testData.DoTests(t)

    // Accept: doesn't exist
    testData = NewTest("","text/html","1")
    testData.AddHandler("text/html","1")
    testData.DoTests(t)
    
    // Accept: */*
    testData = NewTest("*/*","text/html","1")
    testData.AddHandler("text/html","1")
    testData.DoTests(t)

    // Accept: text/*
    testData = NewTest("text/*","text/html","1")
    testData.AddHandler("text/html","1")
    testData.DoTests(t)

    // Accept: foo/*
    testData = NewTest("text/plain","","")
    testData.AddHandler("text/html","1")
    testData.DoTests(t)
}
