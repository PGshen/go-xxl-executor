package log

import (
	"os"
	"testing"
)

func TestStdStreamLog(t *testing.T) {
	h, _ := NewStreamHandler(os.Stdout)
	s := NewDefault(h)
	s.Info("hello world")

	s.Close()

	s.Info("can not log")

	Info("hello world")

	SetHandler(os.Stderr)

	Infof("%s %d", "Hello", 123)

	SetLevel(LevelError)

	Infof("%s %d", "Hello", 123)

	SetLevelByName("info")

	Infof("%s %d", "Hello", 123)

	Fatalf("%s %d", "Hello", 123)
}

func TestRotatingFileLog(t *testing.T) {
	path := "./test_log"
	os.RemoveAll(path)

	os.Mkdir(path, 0777)
	fileName := path + "/test.log"

	h, err := NewRotatingFileHandler(fileName, 1000, 2)
	if err != nil {
		t.Fatal(err)
	}
	s := NewDefault(h)
	s.Info("hello world")
	s.Warn("Warning...")
	s.Debug("Warning...")
	s.Error("Warning...")
	s.Trace("Warning...")
	s.Fatal("Warning...")
	s.Close()
}
