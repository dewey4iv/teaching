package main

import (
	"bytes"
	"io"
	"os"
	"strconv"
)

type VisitCounter interface {
	GetVisits() int64
	SetVisits(int64)
}

func SaveVisits(filename string, v VisitCounter) error {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.WriteString(strconv.Itoa(int(v.GetVisits()))); err != nil {
		return err
	}

	return nil
}

func ReadVisits(filename string, v VisitCounter) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	var buf bytes.Buffer

	if _, err := io.Copy(&buf, file); err != nil {
		return err
	}

	visits, err := strconv.Atoi(buf.String())
	if err != nil {
		return err
	}

	v.SetVisits(int64(visits))

	return nil
}
