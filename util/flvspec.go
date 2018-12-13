package util

type FlvFileHeader struct {
	Signature  string
	Version    uint8
	TypeFlags  uint8
	DataOffset uint32
}

type FlvFiletags struct {
	TagType           uint8
	DataSize          uint32
	Timestamp         uint32
	TimestampExtended uint8
	StreamID          uint32
}

type FlvFileBody struct {
	PreviousTagSize uint32
	Tag             FlvFiletags
}

type FlvFileSpec struct {
	TotalSize int
	Header    FlvFileHeader
	Bodys     []*FlvFileBody
}

func NewFlvFileSpec() *FlvFileSpec {
	fs := &FlvFileSpec{
		TotalSize: 0,

		Bodys: make([]*FlvFileBody, 0),
	}
	return fs
}
