package validations

type ImageValidation struct {
	MinSize   float32
	MaxSize   float32
	MinWidth  int
	MaxWidth  int
	MinHeight int
	MaxHeight int
	Format    []string
}
