package types

type IOMetadata struct {
	InputFilename string
	Value         interface{}
	InfoHash      string
}

type ResultTarget int

const (
	ObjectResultTarget ResultTarget = iota
	MetadataResultTarget
)
