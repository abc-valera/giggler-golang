package greetings

import "net/http"

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from the golang template!"))
}
