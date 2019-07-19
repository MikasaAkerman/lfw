package writer

import "os"

type fiSort []os.FileInfo

func (a fiSort) Len() int           { return len(a) }
func (a fiSort) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a fiSort) Less(i, j int) bool { return a[i].ModTime().Before(a[j].ModTime()) }
