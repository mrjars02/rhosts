// Provides the web server for rhosts to relay altered content
package serve

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"net/http"
	"jbreich/rhosts/serve/certutils"
	"jbreich/rhosts/cfg"
	"log"
	"os"
	"path"
	"math/big"
	"time"
	"net"
)

func Start(exit chan bool) {
	config := cfg.Create()
	if config.WebServer.Enabled == false {
		log.Print("Webserver was disabled in the config file")
		exit <- true
		return
	}




	go httpServer()
	go httpsServer(path.Join(config.System.Var, "certs"))
}

func httpServer(){
	err := http.ListenAndServe("127.0.0.1:80", http.HandlerFunc(httpHandler))
	if (err != nil) {log.Fatal("Failed to start httls server")}
}
func httpsServer(certPath string) {
	var err error

	// Create certificates if they do not exist
	// CA
	err = os.MkdirAll(certPath,0755)
	if (err != nil){log.Fatal("Could not create cert path: " + err.Error())}

	var ca certutils.CA
	ca.Location = certPath
	ca.Name = "ca"
	ca.X509Cert = &x509.Certificate{
		SerialNumber: big.NewInt(2019),
		Subject: pkix.Name{
			Organization:  []string{"Rhost"},
			Country:       []string{""},
			Province:      []string{""},
			Locality:      []string{"Your computer"},
			StreetAddress: []string{""},
			PostalCode:    []string{""},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	ca,err = ca.Initialize()
	if(err!=nil){log.Fatal("Failed to build the ca: " + err.Error())}
	err = ca.WriteToFile()
	if(err!=nil){log.Fatal("Failed to write the ca: " + err.Error())}

	// Server Certificate
	var serverCa certutils.CA
	serverCa.Location = certPath
	serverCa.Name = "server"
	serverCa.Signer = &ca

	serverCa.X509Cert = &x509.Certificate{
		SerialNumber: big.NewInt(2019),
		Subject: pkix.Name{
			Organization:  []string{"The Website"},
			Country:       []string{"US"},
			Province:      []string{""},
			Locality:      []string{"webserver"},
			StreetAddress: []string{""},
			PostalCode:    []string{""},
		},
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	serverCa,err = serverCa.Initialize()
	if(err!=nil){log.Fatal("Failed to build the server certificate: " + err.Error())}
	serverCa.WriteToFile()
	if(err!=nil){log.Fatal("Failed to write the server certificate: " + err.Error())}


	// Starting the server
	err = http.ListenAndServeTLS("127.0.0.1:443", path.Join(certPath, "server.crt"), path.Join(certPath, "server.key"), http.HandlerFunc(httpHandler))
	if (err != nil) {log.Fatal("Failed to start httls server: " + err.Error())}
	return
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Test", 200)
}
