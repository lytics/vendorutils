package vendorutils

import (
	"os"
	"path"
	"testing"

	"github.com/kardianos/govendor/vcs"
)

func TestInvalidDirTest(t *testing.T) {
	path := "@"

	v, err := gitExists(path)
	if err != nil {
		t.Errorf("an invalid path returns no error!")
	}
	t.Logf("%#v %v", v, err)
}

func TestGitDir(t *testing.T) {
	ldir := path.Join(os.Getenv("GOPATH"), "src", "github.com", "lytics", "vendorutils")

	v, err := gitExists(ldir)
	if err != nil {
		t.Errorf("vendorutils dir should return no errors: %#v", err)
	}
	if !v {
		t.Errorf("vendorutils dir should return true!")
	}
}

func TestFalseHgDir(t *testing.T) {
	ldir := path.Join(os.Getenv("GOPATH"), "src", "github.com", "lytics", "vendorutils")

	v, err := hgExists(ldir)
	if err != nil {
		t.Errorf("mag dir should return no errors: %#v", err)
	}
	if v {
		t.Errorf("mag dir should return false; not hg!")
	}
}

func TestHgDir(t *testing.T) {
	ldir := path.Join(os.Getenv("GOPATH"), "src", "github.com", "lytics", "vendorutils", "test", "hgtest")

	v, err := hgExists(ldir)
	if err != nil {
		t.Errorf("hg dir should return no errors: %#v", err)
	}
	if !v {
		t.Errorf("hg dir should return true")
	}
}

func TestBzrDir(t *testing.T) {
	ldir := path.Join(os.Getenv("GOPATH"), "src", "github.com", "lytics", "vendorutils", "test", "bzrtest")

	v, err := bzrExists(ldir)
	if err != nil {
		t.Errorf("bzr dir should return no errors: %#v", err)
	}
	if !v {
		t.Errorf("bzr dir should return true")
	}
}

func TestFindHgDir(t *testing.T) {
	ldir := path.Join(os.Getenv("GOPATH"), "src", "github.com", "lytics", "vendorutils", "test", "hgtest", "src")
	edir := path.Join(os.Getenv("GOPATH"), "src", "github.com", "lytics", "vendorutils", "test", "hgtest")

	v, err := FindVCSRoot(ldir)
	if err != nil {
		t.Errorf("hg dir should return no errors: %#v", err)
	}
	if v != edir {
		t.Errorf("hg dir unexpected: %s", v)
	}
}

func TestFindVCS(t *testing.T) {

	mdir := path.Join(os.Getenv("GOPATH"), "src", "github.com", "lytics", "vendorutils")
	ldir := path.Join(os.Getenv("GOPATH"), "src", "github.com", "lytics", "vendorutils", "test")

	v, err := FindVCSRoot(ldir)
	if err != nil {
		t.Errorf("no error should be returned for a valid path")
	}
	if v != mdir {
		t.Errorf("path returned should match the home directory: %s", v)
	}

}

func TestActualBug(t *testing.T) {
	mdir := "github.com/lytics/gowrapmx4j/cassandra"

	v, err := FindVCSRoot(mdir)
	if err != nil {
		t.Errorf("error should not be returned from GOPATH, however it will not return a valid path")
	}
	if v == mdir {
		t.Errorf("path returned should be empty: %s", v)
	}
}

func TestGVVcsInfo(t *testing.T) {
	mdir := path.Join(os.Getenv("GOPATH"), "src", "github.com", "lytics", "vendorutils")
	vi, err := vcs.FindVcs(os.Getenv("GOPATH"), mdir)
	if err != nil {
		t.Logf("FindVcs err: %v", err)
	}
	t.Logf("%#v", vi)
}
