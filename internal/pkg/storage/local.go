package storage

import ()

type Local struct{}

func (l *Local) Write(path string) ([]byte, error) {

	return []byte{}, nil

}

func (l *Local) Read(path string) ([]byte, error) {

	return []byte{}, nil

}
