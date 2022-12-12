package document

import (
	"io"

	"github.com/NTSka/dxf-go/core"
	"github.com/NTSka/dxf-go/sections"
	"github.com/NTSka/dxf-go/tables"
)

// DxfDocument the representation of a full dxf document.
type DxfDocument struct {
	Header   *sections.HeaderSection
	Tables   *tables.TablesSection
	Entities *sections.EntitiesSection
	Blocks   sections.BlocksSection
}

// Equals compares against the other DxfDocument for equality.
func (doc DxfDocument) Equals(other *DxfDocument) bool {
	return doc.Header.Equals(other.Header) &&
		doc.Tables.Equals(other.Tables) &&
		doc.Entities.Equals(other.Entities) &&
		doc.Blocks.Equals(other.Blocks)
}

// DxfDocumentFromStream reads a DxfDocument from the stream.
func DxfDocumentFromStream(stream io.Reader) (*DxfDocument, error) {
	doc := new(DxfDocument)

	doc.Header = new(sections.HeaderSection)
	doc.Tables = new(tables.TablesSection)
	doc.Entities = new(sections.EntitiesSection)
	doc.Blocks = make(sections.BlocksSection)

	sectionParsers := map[string]func(slice core.TagSlice) error{
		"HEADER": func(slice core.TagSlice) error {
			header := sections.NewHeaderSection(slice)
			doc.Header = header
			return nil
		},
		"TABLES": func(slice core.TagSlice) error {
			section, err := tables.NewTablesSection(slice)
			doc.Tables = section
			return err
		},
		"ENTITIES": func(slice core.TagSlice) error {
			section, err := sections.NewEntitiesSection(slice)
			doc.Entities = section
			return err
		},
		"BLOCKS": func(slice core.TagSlice) error {
			section, err := sections.NewBlocksSection(slice)
			doc.Blocks = section
			return err
		},
	}

	next := core.Tagger(stream)
	tags := core.TagSlice(core.AllTags(next))

	stopTag := core.NewTag(0, core.NewStringValue("EOF"))
	endOfChunk := core.NewTag(0, core.NewStringValue("ENDSEC"))
	for _, sectionTags := range tables.SplitTagChunks(tags, stopTag, endOfChunk) {
		sectionType := sectionTags[1].Value.ToString()

		if parserFunc, ok := sectionParsers[sectionType]; ok {
			err := parserFunc(sectionTags)
			if err != nil {
				return nil, err
			}
		} else {
			core.Log.Printf("Ignoring unsupported Section type: %+v\n", sectionType)
		}
	}

	return doc, nil
}
