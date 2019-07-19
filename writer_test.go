package writer

import (
	"os"
	"testing"
	"time"
)

func TestWrite(t *testing.T) {
	dir := "log"
	size := 1024 * 1024
	d := time.Second * 100
	m := 2
	content := "------test---------\n"
	w := NewWriter(NewOption(OptionLogdir, dir), NewOption(OptionMaxFiles, m),
		NewOption(OptionMaxSize, size), NewOption(OptionDuration, d))

	tm := time.NewTimer(time.Minute)
	tk := time.NewTicker(time.Millisecond)

loop:
	for {
		select {
		case <-tk.C:
			_, err := w.Write([]byte(content))
			if err != nil {
				t.Fatal(err)
			}
		case <-tm.C:
			tk.Stop()
			tm.Stop()
			break loop
		}
	}

	file, err := os.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	fis, err := file.Readdir(-1)
	file.Close()
	if err != nil {
		t.Fatal(err)
	}

	if len(fis) > m {
		t.Error("file count greater than limit")
	}

	for _, fi := range fis {
		if fi.Size() > int64(size+len(content)) {
			t.Errorf("file %s size greater than limit", fi.Name())
		}
	}

	if time.Now().Sub(w.btime) > d*2 {
		t.Error("create date invalid")
	}
}
