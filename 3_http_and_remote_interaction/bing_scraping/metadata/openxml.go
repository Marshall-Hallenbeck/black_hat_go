package metadata

import (
	"archive/zip"
	"encoding/xml"
	"strings"
)

type OfficeCoreProperty struct {
	XMLName        xml.Name `xml:"coreProperties"`
	Creator        string   `xml:"creator"`
	LastModifiedBy string   `xml:"lastModifiedBy"`
}

type OfficeAppProperty struct {
	XMLName     xml.Name `xml:"Properties"`
	Application string   `xml:"Application"`
	Company     string   `xml:"Company"`
	Version     string   `xml:"AppVersion"`
}

// maintains a relationship of major version numbers to recognizable release years
var OfficeVersions = map[string]string{
	"16": "2016",
	"15": "2013",
	"14": "2010",
	"12": "2007",
	"11": "2003",
}

func (a *OfficeAppProperty) GetMajorVersion() string {
	// splits the XML data’s AppVersion value to retrieve the major version number
	tokens := strings.Split(a.Version, ".")

	if len(tokens) < 2 {
		return "Unknown"
	}
	// using the split AppVersion value and the OfficeVersions map to retrieve the release year
	v, ok := OfficeVersions[tokens[0]]
	if !ok {
		return "Unknown"
	}
	return v
}

// accepts a *zip.Reader, which represents an io.Reader for ZIP archives
func NewProperties(r *zip.Reader) (*OfficeCoreProperty, *OfficeAppProperty, error) {
	var coreProps OfficeCoreProperty
	var appProps OfficeAppProperty

	// using "r", a zip.Reader, iterate through all the files in the archive
	for _, f := range r.File {
		// check the filenames (f.Name)
		switch f.Name {
		// If a filename matches one of the two property filenames, call the process() function
		case "docProps/core.xml":
			if err := process(f, &coreProps); err != nil {
				return nil, nil, err
			}
		case "docProps/app.xml":
			if err := process(f, &appProps); err != nil {
				return nil, nil, err
			}
		default:
			continue
		}
	}
	return &coreProps, &appProps, nil
}

// accepts a zip.File and a generic interface{} type to allow for the file contents to be assigned into any data type
func process(f *zip.File, prop interface{}) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()
	// read the content of the file and unmarshal the XML into the "prop" generic interface{} struct
	if err := xml.NewDecoder(rc).Decode(&prop); err != nil {
		return err
	}
	return nil
}
