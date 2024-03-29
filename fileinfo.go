package fileinfo

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"syscall"
	"time"

	"github.com/gonutz/w32"
)

// FileInfo contains extra metadata information on the file. Only .dll and .exe files can contain the extra information.
type FileInfo struct {
	fileName        string
	fileVersionInfo []byte
	translations    []string
	stat            os.FileInfo
}

// New returns a new FileInfo object based on the specified file name path.
func New(path string) (*FileInfo, error) {
	fi := &FileInfo{fileName: path}
	size := w32.GetFileVersionInfoSize(path)
	if size <= 0 {
		return nil, fmt.Errorf(fmt.Sprintf("unable to get File Version Information from %q", path))
	}
	fi.fileVersionInfo = make([]byte, size)
	if ok := w32.GetFileVersionInfo(path, fi.fileVersionInfo); ok {
		fi.translations, ok = w32.VerQueryValueTranslations(fi.fileVersionInfo)
		if !ok {
			return nil, fmt.Errorf("VerQueryValueTranslations failed")
		}
		if len(fi.translations) == 0 {
			return nil, fmt.Errorf("no translations found")
		}
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("unable to open file %q: %v", path, err)
	}
	defer f.Close()

	s, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("unable to get Stat info from file %q: %v", path, err)
	}
	fi.stat = s
	return fi, nil
}

// FileDescription returns the FileDescription property of the file or "-" if no FileDescription is found
func (fi *FileInfo) FileDescription() string {
	fixed, ok := w32.VerQueryValueString(fi.fileVersionInfo, fi.translations[0], w32.FileDescription)
	if ok {
		return fixed
	}
	return "-"
}

// Stat returns the os.FileInfo object from the file
func (fi *FileInfo) Stat() os.FileInfo {
	return fi.stat
}

// FileVersion returns the FileVersion property of the file or "-" if no ProductVersion property is found
func (fi *FileInfo) FileVersion() string {
	fixed, ok := w32.VerQueryValueString(fi.fileVersionInfo, fi.translations[0], w32.FileVersion)
	if ok {
		return fixed
	}
	return "-"
}

// ProductVersion returns the ProductVersion property of the file or "-" if no ProductVersion property is found
func (fi *FileInfo) ProductVersion() string {
	fixed, ok := w32.VerQueryValueString(fi.fileVersionInfo, fi.translations[0], w32.ProductVersion)
	if ok {
		return fixed
	}
	return "-"
}

// ProductName returns the ProductName property of the file or "-" if no ProductName property is found
func (fi *FileInfo) ProductName() string {
	fixed, ok := w32.VerQueryValueString(fi.fileVersionInfo, fi.translations[0], "ProductName")
	if ok {
		return fixed
	}
	return "-"
}

// OriginalFilename returns the OriginalFilename property of the file or "-" if no OriginalFilename property is found
func (fi *FileInfo) OriginalFilename() string {
	fixed, ok := w32.VerQueryValueString(fi.fileVersionInfo, fi.translations[0], w32.OriginalFilename)
	if ok {
		return fixed
	}
	return "-"
}

// Custom returns the specified "propName" property of the file or "-" if the custom "propName" is not found.
// propName should not contain any space and individual words start with a capital letter, eg. "ProductName" for "Product name" property.
func (fi *FileInfo) Custom(propName string) string {
	fixed, ok := w32.VerQueryValueString(fi.fileVersionInfo, fi.translations[0], propName)
	if ok {
		return fixed
	}
	return "-"
}

// Hash returns the SHA256 hash of the file
func (fi *FileInfo) Hash() string {
	hasher := sha256.New()
	s, err := ioutil.ReadFile(fi.fileName)
	hasher.Write(s)
	if err != nil {
		log.Fatal(err)
	}
	return hex.EncodeToString(hasher.Sum(nil))
}

// CreateTime returns the create time on windows systems
func (fi *FileInfo) CreateTime() time.Time {
	d := fi.stat.Sys().(*syscall.Win32FileAttributeData)
	return time.Unix(0, d.CreationTime.Nanoseconds())
}
