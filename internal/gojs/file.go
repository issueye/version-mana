package gojs

import (
	"bufio"
	"fmt"
	"os"
)

type FileOperation struct {
	file *os.File
	w    *bufio.Writer
}

func CreateFileOp(file *os.File) *FileOperation {
	fo := new(FileOperation)
	fo.w = bufio.NewWriter(file)
	return fo
}

func (f *FileOperation) WriteString(content string) {
	f.w.WriteString(fmt.Sprintf("%s\n", content))
}

func (f *FileOperation) Flush() {
	f.w.Flush()
}
