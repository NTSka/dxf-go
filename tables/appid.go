package tables

import (
	"github.com/NTSka/dxf-go/core"
)

// AppID representation.
type AppID struct {
	core.DxfParseable
	Name string
}

// Equals tests equality against another AppID. It only considers the values of the attributes
// on AppID struct, not on parent core.DxfParseable.
func (l AppID) Equals(other core.DxfElement) bool {
	if otherAppID, ok := other.(*AppID); ok {
		return l.Name == otherAppID.Name
	}
	return false
}

// NewAppID builds a new AppID from a tag slice.
func NewAppID(tags core.TagSlice) (*AppID, error) {
	appID := new(AppID)

	appID.Init(map[int]core.TypeParser{
		2: core.NewStringTypeParserToVar(&appID.Name),
	})

	err := appID.Parse(tags)
	return appID, err
}

// NewAppIDTable parses the slice of tags into a table that maps the AppID name to
// the parsed AppID object.
func NewAppIDTable(tags core.TagSlice) (Table, error) {
	table := make(Table)

	tableSlices, err := TableEntryTags(tags)
	if err != nil {
		return table, err
	}

	for _, slice := range tableSlices {
		appID, err := NewAppID(slice)
		if err != nil {
			return nil, err
		}
		table[appID.Name] = appID
	}

	return table, nil
}
