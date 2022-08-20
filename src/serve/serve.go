// Provides the web server for rhosts to relay altered content
package serve

import (
	"net/http"
	"jbreich/rhosts/cfg"
	"log"
)

func Start(exit chan bool) {
	config := cfg.Create()
	if config.WebServer.Enabled == false {
		log.Print("Webserver was disabled in the config file")
		exit <- true
		return
	}
	go httpServer()
	go httpsServer(config.System.Var + "/certs/")

}

func httpServer() (err error) {
	err = http.ListenAndServe("127.0.0.1:80", http.HandlerFunc(httpHandler))
	return
}
func httpsServer(certLoc string) (err error) {
	err = http.ListenAndServeTLS("127.0.0.1:80", certLoc+"ca.crt", certLoc+"ca.key", http.HandlerFunc(httpHandler))
	return
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Test", 200)
}
