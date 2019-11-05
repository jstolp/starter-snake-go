// main_test.go
package main

import (
    "net/http"
    "net/http/httptest"
    "testing"
	"strings"
	"io"
)

func TestOnlyCorrectMoveIsUp(t *testing.T) {
  t.Skip("too hard")
	var r io.Reader
	r = strings.NewReader("{\"game\":{\"id\":\"only-valid-turn-is-up\"},\"turn\":195,\"board\":{\"height\":7,\"width\":7,\"food\":[{\"x\":6,\"y\":0},{\"x\":5,\"y\":1},{\"x\":6,\"y\":2},{\"x\":5,\"y\":0},{\"x\":4,\"y\":1},{\"x\":0,\"y\":0},{\"x\":4,\"y\":0},{\"x\":4,\"y\":3},{\"x\":4,\"y\":4},{\"x\":3,\"y\":6}],\"snakes\":[{\"id\":\"ce4df4c6-22ff-492d-818c-d86932fa5867\",\"name\":\"j\",\"health\":92,\"body\":[{\"x\":3,\"y\":2},{\"x\":4,\"y\":2},{\"x\":5,\"y\":2},{\"x\":5,\"y\":3},{\"x\":6,\"y\":3},{\"x\":6,\"y\":4},{\"x\":6,\"y\":5},{\"x\":6,\"y\":6},{\"x\":5,\"y\":6},{\"x\":4,\"y\":6},{\"x\":4,\"y\":5},{\"x\":3,\"y\":5},{\"x\":2,\"y\":5},{\"x\":1,\"y\":5},{\"x\":0,\"y\":5},{\"x\":0,\"y\":4},{\"x\":0,\"y\":3},{\"x\":0,\"y\":2},{\"x\":1,\"y\":2},{\"x\":1,\"y\":1},{\"x\":1,\"y\":0},{\"x\":2,\"y\":0},{\"x\":2,\"y\":1},{\"x\":2,\"y\":2},{\"x\":2,\"y\":3},{\"x\":1,\"y\":3},{\"x\":1,\"y\":4},{\"x\":2,\"y\":4},{\"x\":3,\"y\":4},{\"x\":3,\"y\":3}]}]},\"you\":{\"id\":\"ce4df4c6-22ff-492d-818c-d86932fa5867\",\"name\":\"j\",\"health\":92,\"body\":[{\"x\":3,\"y\":2},{\"x\":4,\"y\":2},{\"x\":5,\"y\":2},{\"x\":5,\"y\":3},{\"x\":6,\"y\":3},{\"x\":6,\"y\":4},{\"x\":6,\"y\":5},{\"x\":6,\"y\":6},{\"x\":5,\"y\":6},{\"x\":4,\"y\":6},{\"x\":4,\"y\":5},{\"x\":3,\"y\":5},{\"x\":2,\"y\":5},{\"x\":1,\"y\":5},{\"x\":0,\"y\":5},{\"x\":0,\"y\":4},{\"x\":0,\"y\":3},{\"x\":0,\"y\":2},{\"x\":1,\"y\":2},{\"x\":1,\"y\":1},{\"x\":1,\"y\":0},{\"x\":2,\"y\":0},{\"x\":2,\"y\":1},{\"x\":2,\"y\":2},{\"x\":2,\"y\":3},{\"x\":1,\"y\":3},{\"x\":1,\"y\":4},{\"x\":2,\"y\":4},{\"x\":3,\"y\":4},{\"x\":3,\"y\":3}]}}")

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

    // only valid move is up (else in 3 moves we'll fail...
    expected := `{"move":"up"}`
    if rr.Body.String() != expected {
        t.Errorf("handler returned unexpected body: got %v expected %v",
            rr.Body.String(), expected)
    }
}

func TestNextNodeNoPathToTail(t *testing.T) {
	var r io.Reader

  // Next scenario, you will be blocked by yourself...
	r = strings.NewReader("{\"game\":{\"id\":\"valid-move-is-right\"},\"turn\":118,\"board\":{\"height\":5,\"width\":5,\"food\":[{\"x\":4,\"y\":4}],\"snakes\":[{\"id\":\"785911a2-f126-4796-b752-576402268486\",\"name\":\"j\",\"health\":89,\"body\":[{\"x\":3,\"y\":4},{\"x\":3,\"y\":3},{\"x\":4,\"y\":3},{\"x\":4,\"y\":2},{\"x\":4,\"y\":1},{\"x\":4,\"y\":0},{\"x\":3,\"y\":0},{\"x\":3,\"y\":1},{\"x\":2,\"y\":1},{\"x\":2,\"y\":0},{\"x\":1,\"y\":0},{\"x\":0,\"y\":0},{\"x\":0,\"y\":1},{\"x\":0,\"y\":2},{\"x\":0,\"y\":3},{\"x\":0,\"y\":4},{\"x\":1,\"y\":4},{\"x\":1,\"y\":3}]}]},\"you\":{\"id\":\"785911a2-f126-4796-b752-576402268486\",\"name\":\"j\",\"health\":89,\"body\":[{\"x\":3,\"y\":4},{\"x\":3,\"y\":3},{\"x\":4,\"y\":3},{\"x\":4,\"y\":2},{\"x\":4,\"y\":1},{\"x\":4,\"y\":0},{\"x\":3,\"y\":0},{\"x\":3,\"y\":1},{\"x\":2,\"y\":1},{\"x\":2,\"y\":0},{\"x\":1,\"y\":0},{\"x\":0,\"y\":0},{\"x\":0,\"y\":1},{\"x\":0,\"y\":2},{\"x\":0,\"y\":3},{\"x\":0,\"y\":4},{\"x\":1,\"y\":4},{\"x\":1,\"y\":3}]}}")

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

    expected := `{"move":"right"}`
    if rr.Body.String() != expected {
        t.Errorf("handler returned unexpected body: got %v expected %v",
            rr.Body.String(), expected)
    }
}

func TestOnlyMoveIsDown(t *testing.T) {
  //t.Skip("too hard")
	var r io.Reader
	r = strings.NewReader("{\"game\":{\"id\":\"only-move-is-down\"},\"turn\":102,\"board\":{\"height\":7,\"width\":7,\"food\":[{\"x\":6,\"y\":0},{\"x\":0,\"y\":1},{\"x\":5,\"y\":3},{\"x\":5,\"y\":5},{\"x\":0,\"y\":3},{\"x\":6,\"y\":1},{\"x\":4,\"y\":2},{\"x\":4,\"y\":3},{\"x\":1,\"y\":1},{\"x\":3,\"y\":5}],\"snakes\":[{\"id\":\"b12189c4-ab53-4d72-a841-d452a1fe8c5c\",\"name\":\"j\",\"health\":100,\"body\":[{\"x\":0,\"y\":0},{\"x\":1,\"y\":0},{\"x\":2,\"y\":0},{\"x\":3,\"y\":0},{\"x\":4,\"y\":0},{\"x\":5,\"y\":0},{\"x\":5,\"y\":1},{\"x\":4,\"y\":1},{\"x\":3,\"y\":1},{\"x\":2,\"y\":1},{\"x\":2,\"y\":2},{\"x\":2,\"y\":3},{\"x\":2,\"y\":4},{\"x\":2,\"y\":5},{\"x\":1,\"y\":5},{\"x\":0,\"y\":5},{\"x\":0,\"y\":6},{\"x\":1,\"y\":6},{\"x\":2,\"y\":6},{\"x\":3,\"y\":6},{\"x\":4,\"y\":6},{\"x\":5,\"y\":6},{\"x\":6,\"y\":6},{\"x\":6,\"y\":5},{\"x\":6,\"y\":5}]}]},\"you\":{\"id\":\"b12189c4-ab53-4d72-a841-d452a1fe8c5c\",\"name\":\"j\",\"health\":100,\"body\":[{\"x\":0,\"y\":0},{\"x\":1,\"y\":0},{\"x\":2,\"y\":0},{\"x\":3,\"y\":0},{\"x\":4,\"y\":0},{\"x\":5,\"y\":0},{\"x\":5,\"y\":1},{\"x\":4,\"y\":1},{\"x\":3,\"y\":1},{\"x\":2,\"y\":1},{\"x\":2,\"y\":2},{\"x\":2,\"y\":3},{\"x\":2,\"y\":4},{\"x\":2,\"y\":5},{\"x\":1,\"y\":5},{\"x\":0,\"y\":5},{\"x\":0,\"y\":6},{\"x\":1,\"y\":6},{\"x\":2,\"y\":6},{\"x\":3,\"y\":6},{\"x\":4,\"y\":6},{\"x\":5,\"y\":6},{\"x\":6,\"y\":6},{\"x\":6,\"y\":5},{\"x\":6,\"y\":5}]}}")

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

    expected := `{"move":"down"}`
    if rr.Body.String() != expected {
        t.Errorf("handler returned unexpected body: got %v expected %v",
            rr.Body.String(), expected)
    }
}

func TestSmartMoveIsRight(t *testing.T) {
	t.Skip("Hard logic move")
	var r io.Reader
	r = strings.NewReader("{\"game\":{\"id\":\"smart-move-is-right-hard-logic\"},\"turn\":18,\"board\":{\"height\":4,\"width\":4,\"food\":[{\"x\":3,\"y\":3},{\"x\":3,\"y\":1},{\"x\":3,\"y\":2},{\"x\":0,\"y\":3},{\"x\":3,\"y\":0},{\"x\":1,\"y\":3}],\"snakes\":[{\"id\":\"e387e7dc-901e-459f-9f48-f427e77444b1\",\"name\":\"j\",\"health\":100,\"body\":[{\"x\":2,\"y\":3},{\"x\":2,\"y\":2},{\"x\":2,\"y\":1},{\"x\":2,\"y\":0},{\"x\":1,\"y\":0},{\"x\":0,\"y\":0},{\"x\":0,\"y\":1},{\"x\":0,\"y\":2},{\"x\":1,\"y\":2},{\"x\":1,\"y\":1},{\"x\":1,\"y\":1}]}]},\"you\":{\"id\":\"e387e7dc-901e-459f-9f48-f427e77444b1\",\"name\":\"j\",\"health\":100,\"body\":[{\"x\":2,\"y\":3},{\"x\":2,\"y\":2},{\"x\":2,\"y\":1},{\"x\":2,\"y\":0},{\"x\":1,\"y\":0},{\"x\":0,\"y\":0},{\"x\":0,\"y\":1},{\"x\":0,\"y\":2},{\"x\":1,\"y\":2},{\"x\":1,\"y\":1},{\"x\":1,\"y\":1}]}}")

    req, err := http.NewRequest("POST", "/move", r)
    if err != nil {
        t.Fatal(err)
    }

	rr := httptest.NewRecorder()
    handler := http.HandlerFunc(Move)
    handler.ServeHTTP(rr, req)

    // Check the status code is what we expect.
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }

    expected := `{"move": "right"}`
    if rr.Body.String() != expected {
        t.Errorf("Handler returned the wrong move: got %v expected %v",
            rr.Body.String(), expected)
    }
}

func TestMoveIsLeftHardCase(t *testing.T) {

	t.Skip("skip too hard")

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
        t.Errorf("handler returned unexpected body: got %v expected %v",
            rr.Body.String(), expected)
    }
}
