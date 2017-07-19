package main

import (
	"flag"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/lytics/vendorutils"
	"github.com/naoina/toml"
)

var (
	dirPath       string
	writetomlfile bool
	loglvl        string
)

func main() {

	flag.StringVar(&dirPath, "dirPath", "", "path to Go project above vendor/  eg: $GOPATH/src/github.com/lytics/gowrapmx4j")
	flag.BoolVar(&writetomlfile, "tfile", true, "write the toml file to the Go directory specified")
	flag.StringVar(&loglvl, "loglvl", "info", "logrus log level to use")
	flag.Parse()

	lvl, err := log.ParseLevel(loglvl)
	if err != nil {
		log.Errorf("error parsing log level: %v", err)
		os.Exit(1)
	}
	log.SetLevel(lvl)

	log.Debug("derp: a program to initialize dep Gopkg.tomls from govendor")
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
	versionedpkgs := make(map[string]string)

	// iterate over packages from vendorfile and use VCS revisions to remove
	// duplicate package paths. Perfer the shortest/highest package paths and
	// store in map keyed on revisions.
	for _, p := range vf.Package {
		if r, ok := revs[p.Version]; ok && p.Version != "" {
			if len(r) < len(p.Path) {
				revs[p.Version] = p.Path
			}
		} else if p.Version != "" {
			revs[p.Version] = p.Path
		} else {

			// Use revision if no version found
			if r, ok := revs[p.Revision]; ok {
				if len(r) < len(p.Path) {
					revs[p.Revision] = p.Path
				}
			} else {
				revs[p.Revision] = p.Path
			}
		}

		if p.Version != "" {
			vfpkgs[p.Path] = p.Version
			versionedpkgs[p.Path] = p.Version
		} else {
			vfpkgs[p.Path] = p.Revision
		}
		log.Debugf("%s %s %s", p.Path, p.Revision, p.Version)
	}

	// create a sorted list of package paths from the revision map
	pkgs := make([]string, 0)
	for _, v := range revs {
		pkgs = append(pkgs, v)
	}
	sort.Strings(pkgs)

	// construct the absolute path of the Go package and use FindVCSRoot to
	// recursively find the root directory of the project that is `go get`-able.
	/*
		var gfb bytes.Buffer
		for _, p := range pkgs {
			gps := path.Join(gp, "src")
			v, err := vendorutils.FindVCSRoot(path.Join(gps, p))
			log.Debugf("checking path: '%s/%s' '%s'", gps, p, v)
			if err != nil {
				log.Errorf("error finding VCS[%s]: %v", gps, err)
				os.Exit(1)
			}
			vt := strings.TrimPrefix(v, gps+"/")
			log.Debugf("%s %s", vt, vfpkgs[p])
			fmt.Fprintf(&gfb, "%s %s\n", vt, vfpkgs[p])
		}
	*/
	type constraint struct {
		Name     string `toml:"name"`
		Revision string `toml:"revision"`
		Version  string `toml:"version"`
	}
	type Config struct {
		Constraints []constraint `toml:"constraint"`
	}

	constraints := make([]constraint, 0)

	for _, p := range pkgs {
		gps := path.Join(gp, "src")
		v, err := vendorutils.FindVCSRoot(path.Join(gps, p))
		log.Debugf("checking path: '%s/%s' '%s'", gps, p, v)
		if err != nil {
			log.Errorf("error finding VCS[%s]: %v", gps, err)
			os.Exit(1)
		}
		vt := strings.TrimPrefix(v, gps+"/")
		log.Debugf("%s %s", vt, vfpkgs[p])

		var c constraint
		if _, ok := versionedpkgs[p]; ok {
			c = constraint{Name: vt, Version: versionedpkgs[p]}
		} else {
			c = constraint{Name: vt, Revision: vfpkgs[p]}
		}
		constraints = append(constraints, c)
	}
	conf := &Config{Constraints: constraints}
	tomlbytes, err := toml.Marshal(conf)
	if err != nil {
		log.Errorf("error marshaling toml: %v", err)
		os.Exit(1)
	}

	if writetomlfile {
		ioutil.WriteFile(dirPath+"/Gopkg.toml", tomlbytes, 0644)
	} else {
		log.Infof("toml deps: \n%s", tomlbytes)
	}
}
