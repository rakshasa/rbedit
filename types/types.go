package types

type Input interface {
	Execute(metadata IOMetadata, fn InputResultFunc) error
}

type Output interface {
	Execute(metadata IOMetadata, object interface{}) error
	ResultObject() interface{}
}

type EncodeFunc func(IOMetadata, interface{}) (IOMetadata, []byte, error)
type DecodeFunc func(IOMetadata, []byte) (IOMetadata, interface{}, error)
type InputFunc func(IOMetadata) (IOMetadata, []byte, error)
type InputResultFunc func(IOMetadata, interface{}) error
type OutputFunc func(IOMetadata, []byte) error

type ResultTarget int

const (
	ObjectResultTarget ResultTarget = iota
	MetadataResultTarget
)
