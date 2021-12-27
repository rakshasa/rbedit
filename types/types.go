package types

type IOMetadata struct {
	Input         Input
	InputFilename string
	Value         interface{}
	InfoHash      string
}

type Input interface {
	Execute(metadata IOMetadata, fn InputResultFunc) error
}

type InputResultFunc func(interface{}, IOMetadata) error

type ResultTarget int

const (
	ObjectResultTarget ResultTarget = iota
	MetadataResultTarget

	BatchInputTypeName = "batch"
	FileInputTypeName  = "file"

	FileOutputTypeName    = "file"
	InplaceOutputTypeName = "inplace"
)
