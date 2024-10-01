package areainfo

import (
	"github.com/xr0-org/progstack-ssg/internal/assert"
	"github.com/xr0-org/progstack-ssg/internal/ast/page"
)

type AreaInfo struct {
	theme, rootdir string
	index          *page.Page
	purpose        Purpose
}

func Create(
	theme, rootdir string, index *page.Page, purpose Purpose,
) *AreaInfo {
	return &AreaInfo{theme, rootdir, index, purpose}
}

func (info *AreaInfo) WithNewIndex(index *page.Page) *AreaInfo {
	assert.Assert(index != nil)
	return Create(info.theme, info.rootdir, index, info.purpose)
}

func (info *AreaInfo) GetIndex() (*page.Page, bool) {
	return info.index, info.index != nil
}

func (info *AreaInfo) LinksWithHtmlExt() bool {
	switch info.purpose {
	case PurposeStaticServe:
		return true
	case PurposeDynamicServe:
		return false
	default:
		assert.Assert(false)
		return false
	}
}

func (info *AreaInfo) Theme() string { return info.theme }
func (info *AreaInfo) Root() string  { return info.rootdir }

type Purpose int

const (
	PurposeStaticServe Purpose = iota
	PurposeDynamicServe
)
