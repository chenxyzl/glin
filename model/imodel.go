package model

type IModel interface {
	CoName() string
	IsDirty() bool
	MarkDirty()
	CleanDirty()
	Load() error
	Save() error
	Delete() error
}

type Dirty struct {
	dirty bool
}

func (b *Dirty) IsDirty() bool {
	if b == nil {
		return false
	}
	return b.dirty
}

func (b *Dirty) MarkDirty() {
	b.dirty = true
}

func (b *Dirty) CleanDirty() {
	if b != nil {
		b.dirty = false
	}
}
