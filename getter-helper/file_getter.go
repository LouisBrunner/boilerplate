package getter_helper

import (
	"fmt"
	"net/url"
	"os"

	"github.com/hashicorp/go-getter"

	"github.com/louisbrunner/boilerplate/util"
)

// A custom getter.Getter implementation that uses file copying instead of symlinks. Symlinks are
// faster and use less disk space, but they cause issues in Windows and with infinite loops, so we copy files/folders
// instead.
type FileCopyGetter struct {
	getter.FileGetter
}

// The original FileGetter does NOT know how to do folder copying (it only does symlinks), so we provide a copy
// implementation here
func (g *FileCopyGetter) Get(dst string, u *url.URL) error {
	path := u.Path
	if u.RawPath != "" {
		path = u.RawPath
	}

	// The source path must exist and be a directory to be usable.
	if fi, err := os.Stat(path); err != nil {
		return fmt.Errorf("source path error: %s", err)
	} else if !fi.IsDir() {
		return fmt.Errorf("source path must be a directory")
	}

	return util.CopyFolder(path, dst)
}

// The original FileGetter already knows how to do file copying so long as we set the Copy flag to true, so just
// delegate to it
func (g *FileCopyGetter) GetFile(dst string, u *url.URL) error {
	underlying := &getter.FileGetter{Copy: true}
	return underlying.GetFile(dst, u)
}
