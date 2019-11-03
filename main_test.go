// main_test.go
package main

import (
    "net/http"
    "net/http/httptest"
    "testing"
	"strings"
	"io"
)

func TestMoveIsLeft(t *testing.T) {
   
	var r io.Reader
	r = strings.NewReader("{\"game\":{\"id\":\"answer-is-left\"},\"turn\":86,\"board\":{\"height\":7,\"width\":7,\"food\":[{\"x\":0,\"y\":0},{\"x\":0,\"y\":1},{\"x\":4,\"y\":6},{\"x\":1,\"y\":0},{\"x\":5,\"y\":2},{\"x\":5,\"y\":6},{\"x\":5,\"y\":0},{\"x\":3,\"y\":0},{\"x\":3,\"y\":2},{\"x\":0,\"y\":2}],\"snakes\":[{\"id\":\"629dd101-a2b1-4a45-a2f1-a9bd68d7a800\",\"name\":\"j\",\"health\":100,\"body\":[{\"x\":6,\"y\":0},{\"x\":6,\"y\":1},{\"x\":6,\"y\":2},{\"x\":6,\"y\":3},{\"x\":6,\"y\":4},{\"x\":5,\"y\":4},{\"x\":5,\"y\":3},{\"x\":4,\"y\":3},{\"x\":3,\"y\":3},{\"x\":2,\"y\":3},{\"x\":2,\"y\":4},{\"x\":1,\"y\":4},{\"x\":0,\"y\":4},{\"x\":0,\"y\":5},{\"x\":1,\"y\":5},{\"x\":1,\"y\":6},{\"x\":2,\"y\":6},{\"x\":2,\"y\":6}]}]},\"you\":{\"id\":\"629dd101-a2b1-4a45-a2f1-a9bd68d7a800\",\"name\":\"j\",\"health\":100,\"body\":[{\"x\":6,\"y\":0},{\"x\":6,\"y\":1},{\"x\":6,\"y\":2},{\"x\":6,\"y\":3},{\"x\":6,\"y\":4},{\"x\":5,\"y\":4},{\"x\":5,\"y\":3},{\"x\":4,\"y\":3},{\"x\":3,\"y\":3},{\"x\":2,\"y\":3},{\"x\":2,\"y\":4},{\"x\":1,\"y\":4},{\"x\":0,\"y\":4},{\"x\":0,\"y\":5},{\"x\":1,\"y\":5},{\"x\":1,\"y\":6},{\"x\":2,\"y\":6},{\"x\":2,\"y\":6}]}}")   
	
    req, err := http.NewRequest("POST", "/move", r)
    if err != nil {
        t.Fatal(err)
    }

    // We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(Move)

    // Our handlers satisfy http.Handler, so we can call their ServeHTTP method 
    // directly and pass in our Request and ResponseRecorder.
    handler.ServeHTTP(rr, req)

    // Check the status code is what we expect.
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }

    // only valid move is left
    expected := `{"move": "left"}`
    if rr.Body.String() != expected {
        t.Errorf("handler returned unexpected body: got %v want %v",
            rr.Body.String(), expected)
    }
}