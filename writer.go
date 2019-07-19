package writer

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	_prefix     = "log"
	_suffix     = "log"
	_dir        = "log"
	_timeLayout = "20060102_150405"
)

// Writer ...
type Writer struct {
	file       *os.File      // the log file
	d          time.Duration // the duration to rotate log file
	s          int64         // the size to rotate log file
	m          int           // the max files of log file
	prefix     string        // the prefix of log file name
	suffix     string        // the suffix of log file name
	dir        string        // the directory of log files
	timeLayout string        // the log file name layout
	btime      time.Time     // current log file create time
	mu         *sync.RWMutex
}

// NewWriter create a new log writer
// with given options
func NewWriter(ops ...*Option) *Writer {
	var (
		prefix     = _prefix
		suffix     = _suffix
		dir        = _dir
		timeLayout = _timeLayout
		d          = time.Duration(0)
		mu         = new(sync.RWMutex)
		s          int64
		m          int
		err        error
	)
	for _, op := range ops {
		switch op.GetKey() {
		case OptionDuration:
			if dr, ok := op.GetValue().(time.Duration); ok {
				d = dr
			}
		case OptionPrefix:
			if pr, ok := op.GetValue().(string); ok {
				prefix = pr
			}
		case OptionSuffix:
			if sr, ok := op.GetValue().(string); ok {
				suffix = sr
			}
		case OptionTimeLayout:
			if tl, ok := op.GetValue().(string); ok {
				timeLayout = tl
			}
		case OptionMaxFiles:
			if mf, ok := op.GetValue().(int); ok {
				m = mf
			}
		case OptionMaxSize:
			if ms, ok := op.GetValue().(int); ok {
				s = int64(ms)
			}
			if ms, ok := op.GetValue().(int64); ok {
				s = ms
			}
		case OptionLogdir:
			if ld, ok := op.GetValue().(string); ok {
				dir = ld
			}
		default:
		}
	}
	dir, err = filepath.Abs(dir)
	if err != nil {
		panic(err)
	}
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		panic(err)
	}
	writer := &Writer{nil, d, s, m, prefix, suffix, dir, timeLayout, time.Now(), mu}
	err = writer.rotate()
	if err != nil {
		panic(err)
	}

	return writer
}

// check check if the writer needs to rotate log file
func (w *Writer) check() bool {

	// check file duration=
	if w.d > 0 && time.Now().Sub(w.btime) >= w.d {
		return true
	}

	// check file size
	if w.s > 0 {
		fi, err := w.file.Stat()
		if err != nil {
			panic(err)
		}
		if fi.Size() >= w.s {
			return true
		}
	}

	return false
}

// rotate rotate the log file
func (w *Writer) rotate() error {
	if w.file != nil {
		w.file.Close()
	}

	fileName := fmt.Sprintf("%s%s.%s", w.prefix, time.Now().Format(w.timeLayout), w.suffix)
	file, err := os.OpenFile(path.Join(w.dir, fileName), os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	w.file = file
	w.btime = time.Now()

	dir, err := os.Open(w.dir)
	if err != nil {
		return err
	}
	files, err := dir.Readdir(-1)
	if err != nil {
		return err
	}

	// check file numbers
	if w.m > 0 {
		buf := make([]os.FileInfo, 0)
		for _, fi := range files {
			if strings.HasSuffix(fi.Name(), w.suffix) && strings.HasPrefix(fi.Name(), w.prefix) {
				buf = append(buf, fi)
			}
		}
		if len(buf) > w.m {
			sort.Sort(fiSort(buf))
			for i := 0; i < len(buf)-w.m; i++ {
				err = os.Remove(path.Join(w.dir, buf[i].Name()))
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// Write write log to file through memory
func (w *Writer) Write(b []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	n, err = w.file.Write(b)
	if err != nil {
		return
	}

	// err = w.file.Sync()
	// if err != nil {
	// 	return
	// }

	if w.check() {
		err = w.rotate()
	}
	return
}
