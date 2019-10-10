package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/hyperledger/burrow/execution/wasm/slurp"
	"github.com/hyperledger/burrow/util"
	"github.com/tmthrgd/go-hex"
)

type WASMDotGo struct {
	Package string
	Name    string
	WASMHex string
}

var tmpl = template.Must(template.New("wasm.go").Parse(`package {{ .Package }}

import (
	"github.com/hyperledger/burrow/execution/wasm/slurp"
	hex "github.com/tmthrgd/go-hex"
)

var Bytecode_{{ .Name }} = slurp.MustGunzip(hex.MustDecodeString("{{ .WASMHex }}"))
`))

func main() {
	flag.Parse()
	infile := flag.Arg(0)
	bs, err := ioutil.ReadFile(infile)
	if err != nil {
		util.Fatalf("could not read %s: %v", infile, err)
	}

	base := strings.ReplaceAll(path.Base(infile), "-", "_")
	name := strings.TrimSuffix(base, filepath.Ext(base))
	dir := path.Dir(infile)
	pkg := path.Base(dir)
	outfile := path.Join(dir, fmt.Sprintf("%s.wasm.go", name))

	gzbs, err := slurp.Gzip(bs)
	if err != nil {
		util.Fatalf("could not gzip %s: %v", infile, err)
	}

	w, err := os.Create(outfile)
	if err != nil {
		util.Fatalf("could not create file %s: %v", outfile, err)
	}
	err = tmpl.Execute(w, WASMDotGo{
		Package: pkg,
		Name:    name,
		WASMHex: hex.EncodeUpperToString(gzbs),
	})
	if err != nil {
		util.Fatalf("could not execute template for %s: %v", infile, err)
	}
}
