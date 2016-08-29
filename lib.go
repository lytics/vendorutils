package vendorutils

import (
	"errors"
	"os"
	"path"
)

// scmCheck executes tests for major SCM directories given a path
// returns true if an SCM directory is found
// returns false if nothing matched
func scmCheck(p string) (bool, error) {
	var err error = nil
	var c bool = false

	check := func(p string, t func(string) (bool, error)) {
		if err != nil {
			return
		}
		if c == true {
			return
		}
		c, err = t(p)
	}

	check(p, gitExists)
	check(p, hgExists)
	check(p, bzrExists)

	return c, err
}

// gitExists detects if .git/ directory exists in a given directory path
func gitExists(p string) (bool, error) {
	gpath := path.Join(p, ".git")

	gdir, err := os.Stat(gpath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, err
		}
	}
	if gdir.IsDir() {
		return true, nil
	} else {
		return false, nil
	}
}

// hgExists detects if .git/ directory exists in a given directory path
func hgExists(p string) (bool, error) {
	hgpath := path.Join(p, ".hg")

	hgdir, err := os.Stat(hgpath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, err
		}
	}
	if hgdir.IsDir() {
		return true, nil
	} else {
		return false, nil
	}
}

// bzrExists detects if .bzr/ directory exists in a given directory path
func bzrExists(p string) (bool, error) {
	hgpath := path.Join(p, ".bzr")

	hgdir, err := os.Stat(hgpath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, err
		}
	}
	if hgdir.IsDir() {
		return true, nil
	} else {
		return false, nil
	}
}

var VCSRootNotFound = errors.New("scm root was not found")

// FindVCSRoot recurses up the path to find the project's root
// based on the existance of an VCS trigger
func FindVCSRoot(p string) (string, error) {
	v, err := scmCheck(p)
	if err != nil {
		return "", err
	}
	if v {
		return p, nil
	} else if p == "/" || p == "." {
		return "", VCSRootNotFound
	}
	return FindVCSRoot(path.Dir(p))
}
