package xmlpatch

import (
	"encoding/xml"
	"fmt"
	"github.com/beevik/etree"
	"strings"
)

type Diff struct {
	XMLName  xml.Name  `xml:"diff"`
	Replaces []Replace `xml:"replace"`
	Adds     []Add     `xml:"add"`
}

type Replace struct {
	Sel  string `xml:"sel,attr"`
	Text string `xml:",chardata"`
}

type Add struct {
	Pos     string `xml:"pos,attr"`
	Sel     string `xml:"sel,attr"`
	Content []byte `xml:",innerxml"`
}

func Patch(docData, xmlDiffData []byte) ([]byte, error) {
	var diff Diff
	if err := xml.Unmarshal(xmlDiffData, &diff); err != nil {
		return nil, fmt.Errorf("failed to parse xml diff with error: %w", err)
	}
	doc := etree.NewDocument()
	err := doc.ReadFromBytes(docData)
	if err != nil {
		return nil, fmt.Errorf("failed to read doc data with error: %w", err)
	}

	for i, replace := range diff.Replaces {
		xpath := replace.Sel
		attributeRefIndex := strings.LastIndex(xpath, "/@")
		if attributeRefIndex != -1 {
			xpath = xpath[:attributeRefIndex]
			fmt.Printf("\n\n--\n\nnew xpath: %v\n\n--\n\n", xpath)
		}
		path, err := etree.CompilePath(xpath)
		if err != nil {
			return nil, fmt.Errorf("failed to compile sel value of diff replace entry #%v. Sel value: '%v'. Error: %w", i, xpath, err)
		}
		elems := doc.FindElementsPath(path)
		if len(elems) != 1 {
			return nil, fmt.Errorf("expected 1 match for '%v' bot got %v", xpath, len(elems))
		}
		elem := elems[0]
		if attributeRefIndex != -1 {
			elem.CreateAttr(replace.Sel[attributeRefIndex+2:], replace.Text)
		} else {
			elem.SetText(replace.Text) // TODO [Max]: test
		}
	}
	return doc.WriteToBytes()
}
