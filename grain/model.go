package grain

type IModel interface {
	Load() error
	Save() error
	MarkDirty() error
	CleanDirty() error
}
