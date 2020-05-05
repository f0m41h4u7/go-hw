package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath string, toPath string, offset, limit int64) error {
	inFile, err := os.Open(fromPath)
	if err != nil {
		return ErrUnsupportedFile
	}
	defer inFile.Close()

	info, err := inFile.Stat()
	if err != nil {
		return err
	}

	switch {
	case (info.Size() == 0) || (info.IsDir()):
		return ErrUnsupportedFile
	case info.Size() < offset:
		return ErrOffsetExceedsFileSize
	case (limit == 0) || (limit > info.Size()-offset):
		limit = info.Size() - offset
	default:
	}

	_, err = inFile.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}

	outFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	bar := pb.Full.Start64(limit)
	defer bar.Finish()
	barReader := bar.NewProxyReader(inFile)
	_, err = io.CopyN(outFile, barReader, limit)
	if err != nil {
		return err
	}

	return nil
}
