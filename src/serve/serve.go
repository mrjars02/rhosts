// Provides the web server for rhosts to relay altered content
package serve

import (
	"net/http"
)

func Start(certLoc string) {
	go httpServer()
	//go httpsServer(&certLoc)

}

func httpServer() (err error) {
	err = http.ListenAndServe("127.0.0.1:80", http.HandlerFunc(httpHandler))
	return
}
func httpsServer(certLoc *string) (err error) {
	err = http.ListenAndServeTLS("127.0.0.1:80", *certLoc+"ca.crt", *certLoc+"ca.key", http.HandlerFunc(httpHandler))
	return
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Test", 200)
}
