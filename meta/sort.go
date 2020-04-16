package meta

import "time"

const baseFormat = "2006-01-02 15:04:05"

type ByUploadTime []Filemeta

func (a ByUploadTime) Len() int {
	return len(a)
}

func (a ByUploadTime) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByUploadTime) Less(i, j int) bool {
	iTime, _ := time.Parse(baseFormat, a[i].UpdateTime)
	jTime, _ := time.Parse(baseFormat, a[j].UpdateTime)
	return iTime.UnixNano() > jTime.UnixNano()
}
