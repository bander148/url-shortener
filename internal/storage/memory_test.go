package storage

import (
	"log"
	"testing"
	"url-shortener/internal/models"
)

func TestSaveAndGet(t *testing.T) {
	pointer := NewMemoryStorage()
	want := models.UrlData{LongUrl: "vanyaglovesha224"}
	_, err := pointer.Save(want)
	if err != nil {
		log.Fatal(err)
	}
	got, err := pointer.Get("1")
	if err != nil {
		log.Fatal(err)
	}
	if got != want {
		t.Errorf("Test FAIL! want : %s , got : %s ", want, got)
	}
}
func TestNotFound(t *testing.T) {
	pointer := NewMemoryStorage()
	exdata := models.UrlData{LongUrl: "vanyaglovesha224"}
	_, err := pointer.Save(exdata)
	if err != nil {
		log.Fatal(err)
	}
	data, err := pointer.Get("ggvp")
	if err != ErrNotFound {
		t.Errorf("Does not return the desired error.Expected : %#v , Return ; %#v ,Data : %s", ErrNotFound, err, data)
	}
}
