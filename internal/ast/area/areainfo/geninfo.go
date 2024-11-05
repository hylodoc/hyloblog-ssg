package areainfo

import (
	"github.com/xr0-org/progstack-ssg/internal/assert"
	"github.com/xr0-org/progstack-ssg/internal/ast/page"
)

type GenInfo struct {
	theme, rootdir string
	index          page.Page
	purpose        Purpose
	head, foot     string
}

func (info *GenInfo) copy() *GenInfo {
	return &GenInfo{
		theme:   info.theme,
		rootdir: info.rootdir,
		index:   info.index,
		purpose: info.purpose,
		head:    info.head,
		foot:    info.foot,
	}
}

func NewGenInfo(theme, rootdir string, purpose Purpose) *GenInfo {
	return &GenInfo{
		theme:   theme,
		rootdir: rootdir,
		purpose: purpose,
	}
}

func (info *GenInfo) WithNewIndex(index page.Page) *GenInfo {
	assert.Assert(index != nil)
	gi := info.copy()
	gi.index = index
	return gi
}

func (info *GenInfo) WithHeadFoot(head, foot string) *GenInfo {
	gi := info.copy()
	gi.head = head
	gi.foot = foot
	return gi
}

func (info *GenInfo) GetIndex() (page.Page, bool) {
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

func (info *GenInfo) Root() string  { return info.rootdir }
func (info *GenInfo) Theme() string { return info.theme }
func (info *GenInfo) Head() string  { return info.head }
func (info *GenInfo) Foot() string  { return info.foot }

type Purpose int

const (
	PurposeStaticServe Purpose = iota
	PurposeDynamicServe
)
