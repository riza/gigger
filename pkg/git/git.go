package git

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"io/ioutil"
	"unsafe"
)

const (
	errNullBytesNotNull = "null bytes not null"
)

type Git struct {
	Index index

	pos int //reader cursor

	reader *bufio.Reader
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

func NewGit() *Git {
	return &Git{}
}

func (g *Git) ParseIndex(data []byte) (index, error) {
	dataReader := bytes.NewReader(data) // weirdo..
	g.reader = bufio.NewReaderSize(dataReader, int(dataReader.Size()))

	index := index{}
	index.Header = header{
		g.b2s(g.readBytes(4)),
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
			item.Name = g.b2s(g.readBytes(int(nameLen)))
			entryLen += int(nameLen)
		} else {
			name := []byte{}
			for {
				b := g.readBytes(1)
				if g.b2s(b) == "\x00" {
					break
				}
				name = append(name, b[0])
			}
			item.Name = g.b2s(name)
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
		index.Extension.Signature = g.b2s(g.readBytes(4))
		index.Extension.Size = binary.BigEndian.Uint32(g.readBytes(4))
		index.Extension.Data = g.b2s(g.readBytes(int(index.Extension.Size)))
		extNum++
	}

	index.Checksum.Checksum = true
	index.Checksum.SHA1 = hex.EncodeToString(g.readBytes(20))

	return index, nil
}

func (g *Git) ParseObject(data []byte) (string, error) {
	b := bytes.NewReader(data)

	z, err := zlib.NewReader(b)
	if err != nil {
		return "", err
	}

	defer z.Close()

	p, err := ioutil.ReadAll(z)
	if err != nil {
		return "", err
	}

	return g.b2s(g.blobRemover(p)), nil
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

func (g *Git) blobRemover(data []byte) []byte {
	var limit int
	for i, b := range data {
		if b == 0 {
			limit = i
			break
		}
	}

	if limit <= 0 {
		return data
	}

	return data[limit+1:]
}

func (g *Git) b2s(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
