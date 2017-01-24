package csv

import (
	"encoding/csv"
	"fmt"
	"os"
	"sync"
)

type csvTable struct {
	h     *os.File
	r     *csv.Reader
	w     *csv.Writer
	index map[string]int
}

type CsvWriter struct {
	uri   string
	files map[string]*csvTable
	lock  sync.RWMutex
}

func NewCsvWriter(dsn string) *CsvWriter {
	return &CsvWriter{uri: dsn, files: map[string]*csvTable{}}
}

func (w *CsvWriter) InsertStringMap(tblname string, data map[string]interface{}) error {
	t, err := w.getCsvByTableName(tblname, data)
	if err != nil {
		return err
	}

	r := make([]string, len(t.index))
	for k, i := range t.index {
		if v, ok := data[k]; ok {
			r[i] = fmt.Sprint(v)
		}
	}
	t.w.Write(r)
	return nil
}

func (w *CsvWriter) InsertKVData(tblname string, data map[string]interface{}) {

}

func (w *CsvWriter) getCsvByTableName(tblname string, data map[string]interface{}) (*csvTable, error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	t, ok := w.files[tblname]
	resetHeader := false
	if ok {
		for k, _ := range data {
			if _, ok := t.index[k]; !ok {
				resetHeader = true
				t.index[k] = len(t.index)
			}
		}

		if !resetHeader {
			return t, nil
		}
	} else {
		name := fmt.Sprintf(w.uri, tblname)
		h, err := os.OpenFile(name, os.O_RDWR|os.O_APPEND|os.O_CREATE, os.ModePerm)
		if err != nil {
			return nil, err
		}

		t = &csvTable{h: h}
		t.r = csv.NewReader(t.h)
		t.w = csv.NewWriter(t.h)
		t.index = map[string]int{}
		w.files[tblname] = t
		h.Seek(0, os.SEEK_SET)
		if head, err := t.r.Read(); err == nil {
			for i, k := range head {
				t.index[k] = i
			}
			h.Seek(0, os.SEEK_END)
		}

		for k, _ := range data {
			if _, ok := t.index[k]; !ok {
				resetHeader = true
				t.index[k] = len(t.index)
			}
		}
	}

	if resetHeader {
		t.h.Seek(0, os.SEEK_SET)
		rs, _ := t.r.ReadAll()
		if len(rs) == 0 {
			r := make([]string, len(t.index))
			for k, i := range t.index {
				r[i] = k
			}
			rs = append(rs, r)
			t.w.WriteAll(rs)

		} else {
			ap := make([]string, len(t.index)-len(rs[0]))
			for i, r := range rs {
				rs[i] = append(r, ap...)
			}
			for k, i := range t.index {
				rs[0][i] = k
			}
			t.h.Seek(0, os.SEEK_SET)
			t.w.WriteAll(rs)
		}
	}

	return t, nil
}

func (w *CsvWriter) Close() error {
	for _, t := range w.files {
		t.h.Close()
	}
	w.files = map[string]*csvTable{}
	return nil
}
