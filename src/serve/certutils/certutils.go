/*
   Certutils manages certificates.

   To Create a Certificate start with:

    var ca CA

   From there you will want to edit:

   - Location

   - Name

   - X509Cert

   - Signer (If it is suppose to be signed by another CA, you must add it before initializing)

   Then initialize it:

    ca.Initialize

   Create the files:

    ca.WriteToFile()


*/
package certutils

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"log"
	"math/big"
	"os"
	"path"
)

// I want to give Shane Utt credit for being the example for how this module creates and signs certs

/*
 * CA holds information regarding a key pair and certificate
 */
type CA struct {
	Location   string // Location is the location of the keys
	Name       string // This is the name of the certificate, it will append ".key" and ".crt" when writing them:w
	TLSConf    *tls.Config
	CRTPEM     []byte            // This will hold the certificate pem byte file until it is writen to files
	PrivKeyPEM []byte            // This will hold the private key pem byte file until it is writen to files
	X509Cert   *x509.Certificate // Make sure to change this to match what you want your certificate to say
	Signer     *CA               // Signer is the address of who you want to sign this key
	keyPair    *rsa.PrivateKey
}

// If a location is given it will check for one at the location first.
// If one does not exist it will create one using the provided certificate.
// If one does exist it will check it against the certificate and decided if a new one needs to be made.
// If a Certificate is not provided then the default template will be used.
// The name it looks for is "ca": ca.key, ca.crt if no name is provided
func (ca CA) Initialize() (CA, error) {
	var err error

	// Check for name
	if ca.Name == "" {
		ca.Name = "ca"
	}

	// Check if CA Exists
	var caBuf CA
	caBuf, err = fillFromFile(ca)
	if err == nil {
		ca = caBuf
	} else {
		log.Print(err)

		// Check for certificate
		if ca.X509Cert == nil {
			return ca, errors.New("No x509 certificate provided")
		}

		// If CA does not exist then create one
		ca, err = ca.Build()
		if err != nil {
			return ca, err
		}
	}

	return ca, err

}

func fillFromFile(ca CA) (CA, error) {
	var err error

	file := path.Join(ca.Location + ca.Name)
	_, err = os.Stat(file + ".crt")
	if err != nil {
		return ca, err
	}
	ca.CRTPEM, err = os.ReadFile(file + ".crt")
	if err != nil {
		return ca, err
	}

	_, err = os.Stat(file + ".key")
	if err != nil {
		return ca, err
	}
	ca.PrivKeyPEM, err = os.ReadFile(file + ".key")
	if err != nil {
		return ca, err
	}

	log.Print("Reading keys from file")

	//  Decode crt pem
	var p *pem.Block

	p, _ = pem.Decode(ca.CRTPEM)
	if p == nil {
		return ca, errors.New("Failed to parse: " + file + ".crt")
	}

	certBuf, err := x509.ParseCertificate(p.Bytes)
	if err != nil {
		return ca, err
	}
	// This is where it needs to check if they match, for now it will just assign it
	ca.X509Cert = certBuf

	// Decode key pem
	p, _ = pem.Decode(ca.PrivKeyPEM)
	if p == nil {
		return ca, errors.New("Failed to parse: " + file + ".key")
	}
	ca.keyPair, err = x509.ParsePKCS1PrivateKey(p.Bytes)
	if err != nil {
		return ca, err
	}

	return ca, err
}

// WriteToFile will create a key and crt file using the Location and Name on the CA
func (ca CA) WriteToFile() (err error) {
	var file string

	file = path.Join(ca.Location, ca.Name+".crt")
	log.Print("Creating: " + file)
	err = os.WriteFile(file, ca.CRTPEM, 0660)
	if err != nil {
		return
	}

	file = path.Join(ca.Location, ca.Name+".key")
	log.Print("Creating: " + file)
	err = os.WriteFile(file, ca.PrivKeyPEM, 0660)

	return
}

// This is run when initializing unless the files already exist. You can force the building process here afterwards if you want to use the current one's x509 certificate but regenerate everything else.
func (ca CA) Build() (CA, error) {
	log.Print("Generating key" + ca.Name)
	var err error

	// Changing the default serial number
	// For some reason it uses 2019 as the default, I need to watch this in the event it changes
	if ca.X509Cert.SerialNumber.Cmp(big.NewInt(2019)) == 0 {
		// generate a random serial number
		// Later should modify this to use issuerDN+serial
		serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
		serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
		if err != nil {
			return ca, err
		}

		ca.X509Cert.SerialNumber = serialNumber
	}

	// create our private and public key
	caKeyPair, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return ca, err
	}
	ca.keyPair = caKeyPair

	if ca.Signer == nil {
		ca.Signer = &ca
	}

	// create the CA
	caBytes, err := x509.CreateCertificate(rand.Reader, ca.X509Cert, ca.Signer.X509Cert, &caKeyPair.PublicKey, ca.Signer.keyPair)
	if err != nil {
		return ca, err
	}

	// pem encode
	crtPEM := new(bytes.Buffer)
	pem.Encode(crtPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})

	privKeyPEM := new(bytes.Buffer)
	pem.Encode(privKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(caKeyPair),
	})

	ca.CRTPEM = crtPEM.Bytes()
	ca.PrivKeyPEM = privKeyPEM.Bytes()

	return ca, err
}

func (ca CA) TlsConfigCreate() (*tls.Config, error) {
	certpool := x509.NewCertPool()
	if certpool.AppendCertsFromPEM(ca.CRTPEM) == false {
		return nil, errors.New("Failed to AppendCertsFromPEM" + ca.Name)
	}
	TLSConf := &tls.Config{
		RootCAs: certpool,
	}

	return TLSConf, nil

}
