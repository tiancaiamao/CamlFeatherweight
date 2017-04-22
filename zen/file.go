package zen

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/juju/errors"
)

func LoadFile(fileName string) (*VM, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, errors.Trace(err)
	}
	defer f.Close()

	// read file header magic number
	var magic [4]byte
	_, err = f.Read(magic[:])
	if err != nil || magic != [4]byte{'m', 'l', '6', '4'} {
		return nil, errors.Trace(errBadMagic)
	}

	// read global offset and value count
	var (
		globalValueOff uint32
		globalNameNum  uint32
	)
	err1 := binary.Read(f, binary.LittleEndian, &globalValueOff)
	err2 := binary.Read(f, binary.LittleEndian, &globalNameNum)
	if err1 != nil || err2 != nil {
		return nil, errors.Trace(errTruncatedFile)
	}

	// read code
	codeLen := int(globalValueOff - 12)
	code := make([]byte, codeLen)
	n, err := io.ReadFull(f, code)
	if err != nil || n != codeLen {
		return nil, errors.Trace(errTruncatedFile)
	}

	// read global values
	fmt.Println(globalNameNum)
	globalValue := make([]*Value, 0, int(globalNameNum))
	for i := 0; i < int(globalNameNum); i++ {
		value, err := readValue(f)
		if err != nil {
			return nil, errors.Trace(errTruncatedFile)
		}

		globalValue = append(globalValue, value)
	}

	return &VM{
		globalValue: globalValue,
		code:        code,
	}, nil
}

func readValue(r io.Reader) (*Value, error) {
	var t [1]byte
	_, err := r.Read(t[:])
	if err != nil {
		return nil, errors.Trace(err)
	}

	switch t[0] {
	case 1:
		return readIntValue(r)
	case 0:
		return readBlockValue(r)
	}
	return nil, errTruncatedFile
}

func readIntValue(r io.Reader) (*Value, error) {
	var buf [8]byte
	_, err := io.ReadFull(r, buf[:])
	if err != nil {
		return nil, err
	}
	v := binary.LittleEndian.Uint64(buf[:])
	return NewIntegerValue(int(v) >> 1), nil
}

func readBlockValue(r io.Reader) (*Value, error) {
	var head int64
	if err := binary.Read(r, binary.LittleEndian, &head); err != nil {
		return nil, errors.Trace(err)
	}

	if head&0xFF == stringTag {
		blobCount := head >> 9
		buf := make([]byte, blobCount*8)
		_, err := io.ReadFull(r, buf)
		if err != nil {
			return nil, errors.Trace(err)
		}
		return NewStringValue(string(buf)), nil
	} else if head&0xFF == doubleTag {
		var f float64
		err := binary.Read(r, binary.LittleEndian, &f)
		if err != nil {
			return nil, errors.Trace(err)
		}
		return NewFloatValue(f), nil
	} else {
		blobCount := int(head >> 36)
		blob := NewBlockValue(blobCount)
		for i := 0; i < blobCount; i++ {
			v, err := readIntValue(r)
			if err != nil {
				return nil, err
			}
			blob.SetField(i, v)
		}
		return blob, nil
	}
}

const (
	numTags   = 1 << 8
	noScanTag = numTags - 5
	stringTag = noScanTag + 1
	arrayTag  = noScanTag + 2
	doubleTag = noScanTag + 3
)
