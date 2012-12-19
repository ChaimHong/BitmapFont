package bf

import (
	"os"
	"testing"
)

func TestPrint(t *testing.T) {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	bf := NewBitmapFont(path+"/HZK16", path+"/ASC16")
	defer bf.CloseFd()

	str := "2B你好吗？"

	_ = bf.GetBitmap(str, true)
}
