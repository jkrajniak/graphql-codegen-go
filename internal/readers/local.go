package readers

import "io/ioutil"

type LocalReader struct {
	fName string
}

func NewLocalReader(fName string) *LocalReader {
	return &LocalReader{fName: fName}
}

func (l *LocalReader) Read() ([]byte, error) {
	return ioutil.ReadFile(l.fName)
}
