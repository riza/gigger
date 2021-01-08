package git

import (
	"bufio"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"os"

	"github.com/riza/gigger/pkg/config"
)

const (
	errNullBytesNotNull = "null bytes not null"
)

type Git struct {
	Index index

	pos int //reader cursor

	reader *bufio.Reader
	config *config.Config
}

type index struct {
	Header    header
	Entries   []entry
	Extension extension
	Checksum  checksum
}

type header struct {
	Signature string
	Version   uint32
	Entries   uint32
}

type entry struct {
	CTime       uint32
	MTime       uint32
	Dev         uint32
	Ino         uint32
	Mode        uint32
	UID         uint32
	GID         uint32
	Size        uint32
	SHA1        string
	Flags       uint16
	AssumeValid bool
	Extended    bool
	Stage       []bool
	Name        string
}

type extension struct {
	Extension int
	Signature string
	Size      uint32
	Data      string
}

type checksum struct {
	Checksum bool
	SHA1     string
}

func New(path string) (*Git, error) {
	g := &Git{}

	data, err := os.Open(path)
	if err != nil {
		return g, err
	}

	dataStat, err := data.Stat()
	if err != nil {
		return g, err
	}

	g.reader = bufio.NewReaderSize(data, int(dataStat.Size()))
	g.Index, err = g.parseIndex()
	if err != nil {
		return g, err
	}

	return g, err
}

func (g *Git) parseIndex() (index, error) {

	index := index{}
	index.Header = header{
		string(g.readBytes(4)),
		binary.BigEndian.Uint32(g.readBytes(4)),
		binary.BigEndian.Uint32(g.readBytes(4)),
	}
	index.Entries = []entry{}

	for i := 0; i < int(index.Header.Entries); i++ {
		item := entry{
			CTime: binary.BigEndian.Uint32(g.readBytes(8)),
			MTime: binary.BigEndian.Uint32(g.readBytes(8)),
			Dev:   binary.BigEndian.Uint32(g.readBytes(4)),
			Ino:   binary.BigEndian.Uint32(g.readBytes(4)),
			Mode:  binary.BigEndian.Uint32(g.readBytes(4)),
			UID:   binary.BigEndian.Uint32(g.readBytes(4)),
			GID:   binary.BigEndian.Uint32(g.readBytes(4)),
			Size:  binary.BigEndian.Uint32(g.readBytes(4)),
		}

		item.SHA1 = hex.EncodeToString(g.readBytes(20))
		item.Flags = binary.BigEndian.Uint16(g.readBytes(2))
		item.AssumeValid = (item.Flags & 0b10000000 << 8) != 0
		item.Extended = (item.Flags & 0b01000000 << 8) != 0
		item.Stage = append(item.Stage, (item.Flags&0b00100000<<8) != 0)
		item.Stage = append(item.Stage, (item.Flags&0b00010000<<8) != 0)
		nameLen := item.Flags & 0xFFF
		entryLen := 62

		if item.Extended && index.Header.Version == 3 {
			//TODO for index v3
		}

		if nameLen < 0xFFF {
			item.Name = string(g.readBytes(int(nameLen)))
			entryLen += int(nameLen)
		} else {
			name := []byte{}
			for {
				b := g.readBytes(1)
				if string(b) == "\x00" {
					break
				}
				name = append(name, b[0])
			}
			item.Name = string(name)
			entryLen++
		}

		padLen := 8 - (entryLen % 8)
		if padLen == 0 {
			padLen = 8
		}

		nullBytes := g.readBytes(padLen)
		if !g.checkNullBytes(nullBytes) {
			return index, errors.New(errNullBytesNotNull)
		}

		index.Entries = append(index.Entries, item)

	}

	indexLen := g.reader.Size()
	extNum := 1

	for g.pos < (indexLen - 20) {
		index.Extension.Extension = extNum
		index.Extension.Signature = string(g.readBytes(4))
		index.Extension.Size = binary.BigEndian.Uint32(g.readBytes(4))
		index.Extension.Data = string(g.readBytes(int(index.Extension.Size)))
		extNum++
	}

	index.Checksum.Checksum = true
	index.Checksum.SHA1 = hex.EncodeToString(g.readBytes(20))

	return index, nil
}

func (g *Git) readBytes(size int) []byte {
	tmp := []byte{}
	for i := 1; i <= size; i++ {
		data, _ := g.reader.ReadByte()
		tmp = append(tmp, data)
	}
	g.pos += size
	return tmp
}

func (g *Git) checkNullBytes(bytes []byte) bool {
	for _, b := range bytes {
		if b != 0 {
			return false
		}
	}
	return true
}
