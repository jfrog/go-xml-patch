package xmlpatch

import (
	"encoding/xml"
	"fmt"
	"github.com/beevik/etree"
	"golang.org/x/exp/slices"
	"strings"
)

type Diff struct {
	XMLName  xml.Name  `xml:"diff"`
	Replaces []Replace `xml:"replace"`
	Adds     []Add     `xml:"add"`
}

type Replace struct {
	Sel     string `xml:"sel,attr"`
	Text    string `xml:",chardata"`
	Content []byte `xml:",innerxml"`
}

type Add struct {
	Pos       string `xml:"pos,attr"`
	Sel       string `xml:"sel,attr"`
	RejectSel string `xml:"rejectsel,attr"`
	Content   []byte `xml:",innerxml"`
}

type Ops int

const (
	ReplaceAutoCreateMissing Ops = iota + 1
)

func Patch(docData, xmlDiffData []byte, options ...Ops) ([]byte, error) {
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
		if err := doReplace(replace, i, doc, options); err != nil {
			return nil, err
		}
	}

	for i, add := range diff.Adds {
		if err := doAdd(add, i, doc); err != nil {
			return nil, err
		}
	}

	doc.Indent(4)
	return doc.WriteToBytes()
}

func doAdd(add Add, i int, doc *etree.Document) error {
	if add.RejectSel != "" {
		uniqPath, err := etree.CompilePath(add.RejectSel)
		if err != nil {
			return fmt.Errorf("compile sel value %q of add diff entry #%d: %w", add.Sel, i, err)
		}

		exists := doc.FindElementPath(uniqPath)
		if exists != nil {
			return nil
		}
	}

	newDoc := etree.NewDocument()
	if err := newDoc.ReadFromBytes(add.Content); err != nil {
		return fmt.Errorf("read content of add diff entry #%d. Sel value: '%v'. Error: %w", i, add.Sel, err)
	}

	path, err := etree.CompilePath(add.Sel)
	if err != nil {
		return fmt.Errorf("compile sel value %q of add diff entry #%d: %w", add.Sel, i, err)
	}

	elem := doc.FindElementPath(path)

	elem.AddChild(newDoc.Root())

	return nil
}

func doReplace(replace Replace, i int, doc *etree.Document, options []Ops) error {
	xpath := replace.Sel
	attributeRefIndex := strings.LastIndex(xpath, "/@")
	if attributeRefIndex != -1 {
		xpath = xpath[:attributeRefIndex]
	}
	path, err := etree.CompilePath(xpath)
	if err != nil {
		return fmt.Errorf("failed to compile sel value of diff replace entry #%v. Sel value: '%v'. Error: %w", i, xpath, err)
	}
	elems := doc.FindElementsPath(path)
	switch len(elems) {
	case 0:
		if !slices.Contains(options, ReplaceAutoCreateMissing) {
			return fmt.Errorf("expected 1 match for '%v', got 0", xpath)
		}
		createMissing(doc, xpath)
		elem := doc.FindElement(xpath)
		if err := doPatch(attributeRefIndex, elem, replace); err != nil {
			return fmt.Errorf("do patch: %w", err)
		}
	case 1:
		if err := doPatch(attributeRefIndex, elems[0], replace); err != nil {
			return fmt.Errorf("do patch: %w", err)
		}

	default:
		return fmt.Errorf("expected 1 match for '%v', got %v", xpath, len(elems))
	}

	return nil
}

func doPatch(attributeRefIndex int, elem *etree.Element, replace Replace) error {
	if attributeRefIndex != -1 {
		elem.CreateAttr(replace.Sel[attributeRefIndex+2:], replace.Text)
	} else {
		if len(replace.Text) > 0 {
			elem.SetText(replace.Text) // TODO [Max]: test
		}
		if len(replace.Content) > 0 {
			newDoc := etree.NewDocument()
			err := newDoc.ReadFromBytes(replace.Content)
			if err != nil {
				return fmt.Errorf("read replace content: %w\n", err)
			}

			elem.Parent().InsertChildAt(elem.Index(), newDoc.Root())

			elem.Parent().RemoveChild(elem)
		}
	}

	return nil
}

func createMissing(doc *etree.Document, xpath string) {
	xpath = strings.TrimPrefix(xpath, "/")
	xpath = strings.TrimPrefix(xpath, "/")
	parts := strings.Split(xpath, "/") // TODO [Max]: probably not ideal
	lastExisting := doc.Root()
	if lastExisting == nil {
		lastExisting = parseElement(parts[0])
		doc.SetRoot(lastExisting)
	}
	i := 1 // skipping root
	for {
		if next := lastExisting.FindElement(parts[i]); next != nil {
			lastExisting = next
			i++
		} else {
			break
		}
	}
	for curr := lastExisting; i < len(parts); i++ {
		elem := parseElement(parts[i])
		curr.AddChild(elem)
		curr = elem
	}
}

func parseElement(part string) *etree.Element {
	tag, attributeKey, attributeValue := extractTagWithAttribute(part)
	element := etree.NewElement(tag)
	if attributeKey != "" {
		element.CreateAttr(attributeKey, attributeValue)
	}
	return element
}

func extractTagWithAttribute(tag string) (string, string, string) {
	if parenthesesStart := strings.Index(tag, "["); parenthesesStart == -1 {
		return tag, "", ""
	} else {
		attributeKey, attributeValue := extractAttribute(tag[parenthesesStart:])
		return tag[:parenthesesStart], attributeKey, attributeValue
	}

}

func extractAttribute(attributeBlock string) (string, string) {
	equalsIndex := strings.Index(attributeBlock, "=")
	key := attributeBlock[2:equalsIndex]
	value := attributeBlock[equalsIndex+1 : len(attributeBlock)-1]
	if value[0] == '"' || value[0] == '\'' {
		value = value[1 : len(value)-1]
	}
	return key, value
}
