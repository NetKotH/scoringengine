// +build ignore

package main

import (
	"bytes"
	"io/ioutil"

	"github.com/brimstone/traefik-cert/client"
	"github.com/dave/jennifer/jen"
)

func main() {
	f := jen.NewFile("main")

	f.HeaderComment("//go:generate go run certs_gen.go")

	jwt, err := ioutil.ReadFile("jwt.key")
	if err != nil {
		panic(err)
	}
	jwt = bytes.TrimSpace(jwt)

	cert, key, err := client.GetCert("certs.sprinkle.cloud", "scoreboard.netkoth.org", string(jwt))
	if err != nil {
		panic(err)
	}
	f.Var().Id("certBytes").Op("=").Lit(string(cert))
	f.Var().Id("keyBytes").Op("=").Lit(string(key))

	err = f.Save("certs.go")
	if err != nil {
		panic(err)
	}
}
