package testhelper

import (
	"io/ioutil"
	"log"
)

//go:generate go get github.com/loggregator/go-bindata/...
//go:generate ../../scripts/generate-test-certs
//go:generate go-bindata -nocompress -pkg testhelper -prefix test-certs/ test-certs/
//go:generate rm -rf test-certs

func Cert(filename string) string {
	contents := MustAsset(filename)

	tmpfile, err := ioutil.TempFile("", "")

	if err != nil {
		log.Fatal(err)
	}

	if _, err := tmpfile.Write(contents); err != nil {
		log.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}

	return tmpfile.Name()
}
