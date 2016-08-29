package vendorutils

import (
	"os"
	"strings"

	"github.com/kardianos/govendor/vendorfile"
)

// ReadVendorFile is a copy from https://github.com/kardianos/govendor/blob/master/context/vendorFile.go
// This internal function from govendor reads vendor.json and returns an informative
// struct on what is the stored in the vendor/.
func ReadVendorFile(vendorRoot, vendorFilePath string) (*vendorfile.File, error) {
	vf := &vendorfile.File{}
	f, err := os.Open(vendorFilePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	err = vf.Unmarshal(f)
	if err != nil {
		return nil, err
	}
	// Remove any existing origin field if the prefix matches the
	// context package root. This fixes a previous bug introduced in the file,
	// that is now fixed.
	for _, row := range vf.Package {
		row.Origin = strings.TrimPrefix(row.Origin, vendorRoot)
	}

	return vf, nil
}
