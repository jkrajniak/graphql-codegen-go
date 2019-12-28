package internal

import (
	"fmt"
	"os"
)

type Outputer interface {
	Write(s string) error
	Close() error
}

type STDOutput struct{}

func (o *STDOutput) Write(s string) error {
	fmt.Print(s)
	return nil
}
func (o *STDOutput) Close() error { return nil }

type FileOutput struct {
	OutputFile *os.File
}

func NewFileOutput(fName string, ) (*FileOutput, error) {
	file, err := os.Create(fName)
	if err != nil {
		return nil, err
	}
	return &FileOutput{file}, nil
}
func (o *FileOutput) Write(s string) error {
	_, err := o.OutputFile.WriteString(s)
	return err
}

func (o *FileOutput) Close() error {
	return o.OutputFile.Close()
}
