package internal

import (
	"fmt"
	"go/format"
	"os"
)

type Outputer interface {
	Write(s string) error
	Writeln(s string) error
	Close() error
	Flush() error
}

type STDOutput struct{
	buff []byte
}

func NewSTDOutput() *STDOutput {
	return &STDOutput{buff: []byte{}}
}

func (o *STDOutput) Write(s string) error {
	o.buff = append(o.buff, []byte(s)...)
	return nil
}
func (o *STDOutput) Writeln(s string) error {
	o.buff = append(o.buff, []byte(s)...)
	o.buff = append(o.buff, []byte("\n")...)
	return nil
}

func (o *STDOutput) Flush() error {
	out, err := format.Source(o.buff)
	if err != nil {
		return err
	}
	fmt.Print(string(out))
	return nil
}

func (o *STDOutput) Close() error {
	return o.Flush()
}

type FileOutput struct {
	OutputFile *os.File
	buff []byte
}

func NewFileOutput(fName string, ) (*FileOutput, error) {
	file, err := os.Create(fName)
	if err != nil {
		return nil, err
	}
	return &FileOutput{OutputFile: file, buff: []byte{}}, nil
}
func (o *FileOutput) Write(s string) error {
	o.buff = append(o.buff, []byte(s)...)
	return nil
}
func (o *FileOutput) Writeln(s string) error {
	o.buff = append(o.buff, []byte(s)...)
	o.buff = append(o.buff, []byte("\n")...)
	return nil
}

func (o *FileOutput) Close() error {
	if err := o.Flush(); err != nil {
		return err
	}
	return o.OutputFile.Close()
}

func (o *FileOutput) Flush() error {
	out, err := format.Source(o.buff)
	if err != nil {
		return err
	}

	_, err = o.OutputFile.Write(out)
	return err
}
