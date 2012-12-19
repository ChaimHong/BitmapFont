package bf

import (
	"fmt"
	"github.com/go-iconv/src"
	"os"
)

const (
	HZ = 0
	YW = 1
)

type Word struct {
	b      []byte
	T      int8
	Bitmap []byte
}

type BitMapFont struct {
	Words []*Word

	fdHZ   *os.File
	fdYW   *os.File
	hzFile string
	ywFile string
	err    error
	print  bool
}

func NewBitmapFont(hzFile, ywFile string) *BitMapFont {
	return &BitMapFont{
		hzFile: hzFile,
		ywFile: ywFile,
	}
}

func (bf *BitMapFont) GetBitmap(str string, p bool) []*Word {
	bf.print = p

	a, b := os.Getwd()
	fmt.Println(a, b)

	cd, err := iconv.Open("gbk", "utf-8")
	if err != nil {
		panic(err)
		return nil
	}
	defer cd.Close()

	bf.Words = make([]*Word, 0)
	gbk := cd.ConvString(str)

	gbkBytes := []byte(gbk)

	for i := 0; i < len(gbkBytes); {
		if uint8(gbkBytes[i]) < 0xa1 {
			bf.Words = append(bf.Words, &Word{gbkBytes[i : i+1], YW, nil})
			i++
		} else {
			bf.Words = append(bf.Words, &Word{gbkBytes[i : i+2], HZ, nil})
			i += 2
		}
	}

	for _, v := range bf.Words {
		switch v.T {
		case YW:
			v.Bitmap = bf.getYW(v.b)
		case HZ:
			v.Bitmap = bf.getHZ(v.b)
		}
	}

	return bf.Words
}

func (bf *BitMapFont) getYW(yw []byte) (bitmap []byte) {
	if bf.fdYW == nil {
		bf.fdYW, bf.err = os.OpenFile(bf.ywFile, os.O_RDONLY, 0777)
		bf.printError()
	}

	_, bf.err = bf.fdYW.Seek(int64(yw[0])*16, os.SEEK_SET)
	bf.printError()

	bitmap = make([]byte, 16)
	_, bf.err = bf.fdYW.Read(bitmap)
	bf.printError()

	if !bf.print {
		return
	}

	for i := 0; i < 16; i++ {
		fmt.Println("")
		for k := uint(0); k < 8; k++ {
			if bitmap[i]&(0x80>>k) == (0x80 >> k) {
				fmt.Print(1)
			} else {
				fmt.Print(" ")
			}
		}
	}
	fmt.Println("")
	return
}

func (bf *BitMapFont) getHZ(gbkBytes []byte) (bitmap []byte) {
	if bf.fdHZ == nil {
		bf.fdHZ, bf.err = os.OpenFile(bf.hzFile, os.O_RDONLY, 0777)
		bf.printError()
	}

	seek := int64((94*(int(gbkBytes[0])-161) + int(gbkBytes[1]) - 161) * 32)
	_, bf.err = bf.fdHZ.Seek(seek, os.SEEK_SET)
	bf.printError()

	bitmap = make([]byte, 32)
	_, bf.err = bf.fdHZ.Read(bitmap)
	bf.printError()

	if !bf.print {
		return
	}

	for i := 0; i < 32; i++ {
		if i%2 == 0 {
			fmt.Println("")
		}
		for k := uint(0); k < 8; k++ {
			if bitmap[i]&(0x80>>k) == (0x80 >> k) {
				fmt.Print(1)
			} else {
				fmt.Print(" ")
			}
		}
	}
	fmt.Println("")
	return
}

func (bf *BitMapFont) printError() {
	if bf.err != nil {
		panic(bf.err)
	}
}

func (bf *BitMapFont) CloseFd() {
	if bf.fdHZ != nil {
		bf.fdHZ.Close()
	}
	if bf.fdYW != nil {
		bf.fdYW.Close()
	}
}
