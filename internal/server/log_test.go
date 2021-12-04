package server

import "testing"

func TestAppend(t *testing.T) {
	l := NewLog()

	r := Record{
		Value: []byte("Test Record-001"),
	}

	off, err := l.Append(r)
	if err != nil {
		t.Error(err)
	}

	if off != 0 {
		t.Error("invalid offset value")
	}
}

func TestRead(t *testing.T) {
	l := NewLog()
	l.Append(Record{
		Value: []byte("Test Record-001"),
	})
	l.Append(Record{
		Value: []byte("Test Record-002"),
	})

	rec, err := l.Read(1)
	if err != nil {
		t.Error(err)
	}

	if string(rec.Value) != "Test Record-002" {
		t.Error("invalid record value")
	}
}
