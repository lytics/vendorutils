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
	dirPath        string
	writeGlockfile bool
	loglvl         string
)

func main() {

	flag.StringVar(&dirPath, "dirPath", "", "path to Go project above vendor/  eg: $GOPATH/src/github.com/lytics/gowrapmx4j")
	flag.BoolVar(&writeGlockfile, "gfile", true, "write the GLOCKFILE to the Go directory specified")
	flag.StringVar(&loglvl, "loglvl", "info", "logrus log level to use")
	flag.Parse()

	lvl, err := log.ParseLevel(loglvl)
	if err != nil {
		log.Errorf("error parsing log level: %v", err)
		os.Exit(1)
	}
	log.SetLevel(lvl)

	log.Debug("magazine: a program to load glock with vendor checkouts from govendor")
	gp := os.Getenv("GOPATH")
	if gp == "" {
		log.Errorf("error reading GOPATH envvar")
		os.Exit(1)
	}

	vf, err := vendorutils.ReadVendorFile(dirPath+"/vendor", dirPath+"/vendor/vendor.json")
	if err != nil {
		log.Errorf("error reading vendorfile: %v\n", err)
		os.Exit(1)
	}

	revs := make(map[string]string)
	vfpkgs := make(map[string]string)

	// iterate over packages from vendorfile and use VCS revisions to remove
	// duplicate package paths. Perfer the shortest/highest package paths and
	// store in map keyed on revisions.
	for _, p := range vf.Package {
		if r, ok := revs[p.Revision]; ok {
			if len(r) < len(p.Path) {
				revs[p.Revision] = p.Path
			}
		} else {
			revs[p.Revision] = p.Path
		}
		vfpkgs[p.Path] = p.Revision
		log.Debugf("%s %s", p.Path, p.Revision)
	}

	// create a sorted list of package paths from the revision map
	pkgs := make([]string, 0)
	for _, v := range revs {
		pkgs = append(pkgs, v)
	}
	sort.Strings(pkgs)

	// construct the absolute path of the Go package and use FindVCSRoot to
	// recursively find the root directory of the project that is `go get`-able.
	var gfb bytes.Buffer
	for _, p := range pkgs {
		gps := path.Join(gp, "src")
		v, err := vendorutils.FindVCSRoot(path.Join(gps, p))
		if err != nil {
			log.Errorf("error finding VCS[%s]: %v", v, err)
			os.Exit(1)
		}
		vt := strings.TrimPrefix(v, gps+"/")
		log.Debugf("%s %s", vt, vfpkgs[p])
		fmt.Fprintf(&gfb, "%s %s\n", vt, vfpkgs[p])
	}

	if writeGlockfile {
		ioutil.WriteFile(dirPath+"/GLOCKFILE", gfb.Bytes(), 0644)
	} else {
		log.Infof("GLOCKFILE deps: \n%s", gfb.Bytes())
	}
}
