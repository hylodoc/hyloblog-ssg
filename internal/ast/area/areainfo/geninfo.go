package areainfo

import (
	"github.com/xr0-org/progstack-ssg/internal/assert"
	"github.com/xr0-org/progstack-ssg/internal/ast/page"
)

type GenInfo struct {
	theme, rootdir string
	index          *page.Page
	purpose        Purpose
}

func NewGenInfo(
	theme, rootdir string, index *page.Page, purpose Purpose,
) *GenInfo {
	return &GenInfo{theme, rootdir, index, purpose}
}

func (info *GenInfo) WithNewIndex(index *page.Page) *GenInfo {
	assert.Assert(index != nil)
	return NewGenInfo(info.theme, info.rootdir, index, info.purpose)
}

func (info *GenInfo) GetIndex() (*page.Page, bool) {
	return info.index, info.index != nil
}

func (info *GenInfo) DynamicLinks() bool {
	switch info.purpose {
	case PurposeStaticServe:
		return false
	case PurposeDynamicServe:
		return true
	default:
		assert.Assert(false)
		return false
	}
}

func (info *GenInfo) Theme() string { return info.theme }
func (info *GenInfo) Root() string  { return info.rootdir }

type Purpose int

const (
	PurposeStaticServe Purpose = iota
	PurposeDynamicServe
)
