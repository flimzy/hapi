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
    Name        string
    Accept      string
    Handlers    map[string]Handle   
    Result      string
    HandlerID   string
}

func NewTest(name,accept,result,id string) *typeTest {
    return &typeTest{
        name,
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
        t.Fatalf("%s: TypeNegotiator returned '%s', expected '%s'\n", testData.Name, negType, testData.Result)
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
        t.Fatalf("%s: Handler identified itself as '%s', expected '%s'\n", testData.Name, id, testData.HandlerID)
    }
}

// Tests 'text/html' handler
func TestTypeNegotiator(t *testing.T) {
    // Accept: matches exactly
    testData := NewTest("test/html 1","text/html","text/html","1")
    testData.AddHandler("text/html","1")
    testData.DoTests(t)

    // Accept: Doesn't match
    testData = NewTest("test/html 2","text/plain","","")
    testData.AddHandler("text/html","1")
    testData.DoTests(t)

    // Accept: doesn't exist
    testData = NewTest("test/html 3","","text/html","1")
    testData.AddHandler("text/html","1")
    testData.DoTests(t)

    // Accept: */*
    testData = NewTest("test/html 4","*/*","text/html","1")
    testData.AddHandler("text/html","1")
    testData.DoTests(t)

    // Accept: text/*
    testData = NewTest("test/html 5","text/*","text/html","1")
    testData.AddHandler("text/html","1")
    testData.DoTests(t)

    // Accept: foo/*
    testData = NewTest("test/html 6","text/plain","","")
    testData.AddHandler("text/html","1")
    testData.DoTests(t)
}

// Tests for '*/*' handler
func TestTypeNegotiator2(t *testing.T) {
    // Accept: text/html
    testData := NewTest("*/* 1","text/html","text/html","1")
    testData.AddHandler("*/*","1")
    testData.DoTests(t)

    testData = NewTest("*/* 2","*/*","*/*","1")
    testData.AddHandler("*/*","1")
    testData.DoTests(t)

    testData = NewTest("*/* 3","","*/*","1")
    testData.AddHandler("*/*","1")
    testData.DoTests(t)
}

// Tests for two routers 1: "text/html", 2: "*/*"
func TestTypeNegotiator3(t *testing.T) {
    testData := NewTest("text/html + */* 1","text/html","text/html","1")
    testData.AddHandler("text/html","1")
    testData.AddHandler("*/*","2")
    testData.DoTests(t)

    testData = NewTest("text/html + */* 2","text/html","text/html","2")
    testData.AddHandler("text/plain","1")
    testData.AddHandler("*/*","2")
    testData.DoTests(t)

    testData = NewTest("text/html + */* 3","text/*","text/*","2")
    testData.AddHandler("text/html","1")
    testData.AddHandler("*/*","2")
    testData.DoTests(t)
}

// Two routers: 1: text/html, 2: text/plain
func TestTypeNegotiator4(t *testing.T) {
    testData := NewTest("text/html + text/plain 1","text/html","text/html","1")
    testData.AddHandler("text/html","1")
    testData.AddHandler("text/plain","2")
    testData.DoTests(t)

    testData = NewTest("text/html + text/plain 2","text/plain","text/plain","2")
    testData.AddHandler("text/html","1")
    testData.AddHandler("text/plain","2")
    testData.DoTests(t)

    testData = NewTest("text/html + text/plain 3","text/*","text/*","1")
    testData.AddHandler("text/*","1")
    testData.AddHandler("text/plain","2")
    testData.DoTests(t)
}

// Three routers: 1: */*, 2: text/html, 3: text/plain
func TestTypeNegotiator5(t *testing.T) {
    testData := NewTest("*/* + text/html + text/plain 1", "text/html","text/html","2")
    testData.AddHandler("*/*","1")
    testData.AddHandler("text/html","2")
    testData.AddHandler("text/plain","3")

    testData = NewTest("*/* + text/html + text/plain 2", "text/html","text/foo","1")
    testData.AddHandler("*/*","1")
    testData.AddHandler("text/html","2")
    testData.AddHandler("text/plain","3")

    testData = NewTest("*/* + text/html + text/plain 3", "text/*","text/*","2")
    testData.AddHandler("*/*","1")
    testData.AddHandler("text/html","2")
    testData.AddHandler("text/plain","3")

}
