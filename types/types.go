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

type Output interface {
	Execute(object interface{}, metadata IOMetadata) error
	ResultObject() interface{}
}

type EncodeFunc func(interface{}) ([]byte, error)
type DecodeFunc func([]byte) (interface{}, error)
type InputFunc func(IOMetadata) (IOMetadata, []byte, error)
type InputResultFunc func(interface{}, IOMetadata) error
type OutputFunc func([]byte, IOMetadata) error

type ResultTarget int

const (
	ObjectResultTarget ResultTarget = iota
	MetadataResultTarget

	BatchInputTypeName = "batch"
	FileInputTypeName  = "file"

	FileOutputTypeName    = "file"
	InplaceOutputTypeName = "inplace"
)
