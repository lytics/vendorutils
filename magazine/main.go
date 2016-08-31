package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/lytics/vendorutils"
)

var (
	dirpath        string
	writeGlockfile bool
)

func main() {
	log.Infof("magazine: a program to load glock with vendor checkouts from govendor")

	flag.StringVar(&dirpath, "dirpath", "", "path to Go project above vendor/  eg: $GOPATH/src/github.com/lytics/gowrapmx4j")
	flag.BoolVar(&writeGlockfile, "gfile", true, "write the GLOCKFILE to the Go directory specified")
	flag.Parse()

	gp := os.Getenv("GOPATH")
	if gp == "" {
		log.Errorf("error reading GOPATH envvar")
		os.Exit(1)
	}

	vf, err := vendorutils.ReadVendorFile(dirpath+"/vendor", dirpath+"/vendor/vendor.json")
	if err != nil {
		log.Errorf("error reading vendorfile: %v\n", err)
		os.Exit(1)
	}

	revs := make(map[string]string)
	vfpkgs := make(map[string]string)

	for _, p := range vf.Package {
		if r, ok := revs[p.Revision]; ok {
			if len(r) < len(p.Path) {
				revs[p.Revision] = p.Path
			}
		} else {
			revs[p.Revision] = p.Path
			fmt.Printf("%#v\n", p)
		}
		vfpkgs[p.Path] = p.Revision
	}

	pkgs := make([]string, 0)
	for _, v := range revs {
		pkgs = append(pkgs, v)
	}
	sort.Strings(pkgs)

	var gfb bytes.Buffer
	for _, p := range pkgs {
		gps := path.Join(gp, "src")
		v, err := vendorutils.FindVCSRoot(path.Join(gps, p))
		if err != nil {
			log.Errorf("error finding VCS[%s]: %v", v, err)
			os.Exit(1)
		}
		vt := strings.TrimPrefix(v, gps+"/")
		fmt.Fprintf(&gfb, "%s %s\n", vt, vfpkgs[p])
	}

	if writeGlockfile {
		ioutil.WriteFile(dirpath+"/GLOCKFILE", gfb.Bytes(), 0644)
	} else {
		log.Infof("GLOCKFILE deps: \n%s", gfb.Bytes())
	}
}
