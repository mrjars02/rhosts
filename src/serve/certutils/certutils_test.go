package certutils

import (
	"testing"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"time"
	"net"
)

// Parent template for all further certificates
var x509CATmpl = &x509.Certificate{
	SerialNumber: big.NewInt(2019),
	Subject: pkix.Name{
		Organization:  []string{"Master CA of Localhost"},
		Country:       []string{"XX"},
		Province:      []string{""},
		Locality:      []string{"Your computer"},
		StreetAddress: []string{""},
		PostalCode:    []string{"5432"},
	},
	NotBefore:             time.Now(),
	NotAfter:              time.Now().AddDate(10, 0, 0),
	IsCA:                  true,
	ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
	KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	BasicConstraintsValid: true,
}

var x509ServerTmpl = &x509.Certificate{
	SerialNumber: big.NewInt(2019),
	Subject: pkix.Name{
		Organization:  []string{"The Website"},
		Country:       []string{"US"},
		Province:      []string{""},
		Locality:      []string{"webserver"},
		StreetAddress: []string{"bus 1"},
		PostalCode:    []string{"8"},
	},
	IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
	NotBefore:    time.Now(),
	NotAfter:     time.Now().AddDate(10, 0, 0),
	SubjectKeyId: []byte{1, 2, 3, 4, 6},
	ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
	KeyUsage:     x509.KeyUsageDigitalSignature,
}


func TestWriteToFile (t *testing.T){
	var ca, server CA
	var err error

	ca.Name = "TestCertificateAuthority"
	ca.X509Cert = x509CATmpl
	ca, err = ca.Initialize()
	if (err != nil){
		t.Error("Failed to genrate the TestCertificateAuthority certificate" + err.Error())
	}


	server.Name = "server"
	server.X509Cert = x509ServerTmpl
	server.Signer = &ca
	server, err = server.Initialize()
	if (err != nil){
		t.Error("Failed to generate the server certificate" + err.Error())
	}

	err = ca.WriteToFile()
	if (err != nil){
		t.Error("Failed to write the TestCertificateAuthority files" + err.Error())
	}
	err = server.WriteToFile()
	if (err != nil){
		t.Error("Failed to write the sever portion" + err.Error())
	}

}
