package fileinfo

import (
	"fmt"
	"log"

	"github.com/gonutz/w32"
)

// FileInfo contains extra metadata information on the file. Only .dll and .exe files can contain the extra information.
type FileInfo struct {
	fileName        string
	fileVersionInfo []byte
	translations    []string
}

// New returns a new FileInfo object based on the specified file name path.
func New(path string) *FileInfo {
	fi := &FileInfo{fileName: path}
	size := w32.GetFileVersionInfoSize(path)
	if size <= 0 {
		panic(fmt.Sprintf("unable to get File Version Information from %q", path))
	}
	fi.fileVersionInfo = make([]byte, size)
	if ok := w32.GetFileVersionInfo(path, fi.fileVersionInfo); ok {
		fi.translations, ok = w32.VerQueryValueTranslations(fi.fileVersionInfo)
		if !ok {
			panic("VerQueryValueTranslations failed")
		}
		if len(fi.translations) == 0 {
			log.Fatalf("no translations found!")
		}
	}
	return fi
}

// GetFileDesc returns the FileDescription property of the file or "-" if no FileDescription is found
func (fi *FileInfo) GetFileDesc() string {
	fixed, ok := w32.VerQueryValueString(fi.fileVersionInfo, fi.translations[0], w32.FileDescription)
	if ok {
		return fixed
	}
	return "-"
}

// GetFileVer returns the FileVersion property of the file or "-" if no ProductVersion property is found
func (fi *FileInfo) GetFileVer() string {
	fixed, ok := w32.VerQueryValueString(fi.fileVersionInfo, fi.translations[0], w32.FileVersion)
	if ok {
		return fixed
	}
	return "-"
}

// GetProdVer returns the ProductVersion property of the file or "-" if no ProductVersion property is found
func (fi *FileInfo) GetProdVer() string {
	fixed, ok := w32.VerQueryValueString(fi.fileVersionInfo, fi.translations[0], w32.ProductVersion)
	if ok {
		return fixed
	}
	return "-"
}

// GetProdName returns the ProductName property of the file or "-" if no ProductName property is found
func (fi *FileInfo) GetProdName() string {
	fixed, ok := w32.VerQueryValueString(fi.fileVersionInfo, fi.translations[0], "ProductName")
	if ok {
		return fixed
	}
	return "-"
}

// GetOrgName returns the OriginalFilename property of the file or "-" if no OriginalFilename property is found
func (fi *FileInfo) GetOrgName() string {
	fixed, ok := w32.VerQueryValueString(fi.fileVersionInfo, fi.translations[0], w32.OriginalFilename)
	if ok {
		return fixed
	}
	return "-"
}

// GetCustom returns the specified "propName" property of the file or "-" if the custom "propName" is not found.
// propName should not contain any space and individual words start with a capital letter, eg. "ProductName" for "Product name" property.
func (fi *FileInfo) GetCustom(propName string) string {
	fixed, ok := w32.VerQueryValueString(fi.fileVersionInfo, fi.translations[0], propName)
	if ok {
		return fixed
	}
	return "-"
}
