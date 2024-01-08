package svg

import (
	"strings"

	"github.com/fpagyu/illustrator"
)

type SvgPath struct {
	id     string
	indent int // 层级

	pathOp illustrator.PathOp

	d     string
	attrs map[string]string
}

type SvgCompoundPath struct {
	*SvgPath
}

func (sp *SvgPath) Id() string {
	return sp.id
}

func (sp *SvgPath) Indent() int {
	return sp.indent
}

func (sp *SvgPath) SetStyle(v string) {
	sp.SetAttr("style", v)
}

func (sp *SvgPath) SetAttr(k, v string) {
	if len(k) == 0 || len(v) == 0 {
		return
	}

	if sp.attrs == nil {
		sp.attrs = make(map[string]string)
	}
	sp.attrs[k] = v
}

func (sp *SvgPath) Attrs() []string {
	r := make([]string, 0, len(sp.attrs)+1)
	if len(sp.id) > 0 {
		r = append(r, Attr("id", sp.id))
	}

	for k, v := range sp.attrs {
		r = append(r, Attr(k, v))
	}

	return r
}

type PathBuilder struct {
	ptype        uint8 // path type
	compoundPath *SvgPath

	strings.Builder
}

func (b *PathBuilder) Reset() {
	b.ptype = 0
	b.Builder.Reset()
}

func (b *PathBuilder) SetAsClip() {
	b.ptype |= illustrator.AI_ClipPath
}

func (b *PathBuilder) SetAsCompound() {
	b.ptype |= illustrator.AI_CompoundPath
	b.compoundPath = &SvgPath{id: "<Compound Path>"}
}

func (b *PathBuilder) UnsetClip() {
	b.ptype &= (^illustrator.AI_ClipPath)
}

func (b *PathBuilder) UnsetCompound() {
	b.compoundPath = nil
	b.ptype &= (^illustrator.AI_CompoundPath)
}

func (b *PathBuilder) IsClip() bool {
	return b.ptype&illustrator.AI_ClipPath == illustrator.AI_ClipPath
}

func (b *PathBuilder) IsCompound() bool {
	return b.ptype&illustrator.AI_CompoundPath == illustrator.AI_CompoundPath
}
