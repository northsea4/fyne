//go:build !android && !ios && !mobile && !wasm

package app

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
)

func TestFileChangedNotifiesOnce(t *testing.T) {
	s := &settings{}
	count := 0
	s.AddListener(func(fyne.Settings) { count++ })
	s.fileChanged()
	if count != 1 {
		t.Fatalf("one settings load notified %d times", count)
	}
}

func TestWatchFileIgnoresSiblingChanges(t *testing.T) {
	test.NewTempApp(t)
	dir, err := os.MkdirTemp("testdata", "watch-scope-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "settings.json")
	if err = os.WriteFile(path, []byte("{}"), 0600); err != nil {
		t.Fatal(err)
	}
	called := make(chan struct{}, 16)
	watcher := watchFile(path, func() { called <- struct{}{} })
	if watcher == nil {
		t.Fatal("watch unavailable")
	}
	defer watcher.Close()
	if err = os.Mkdir(filepath.Join(dir, "other-app"), 0700); err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(filepath.Join(dir, "other.json"), []byte("{}"), 0600); err != nil {
		t.Fatal(err)
	}
	select {
	case <-called:
		t.Fatal("unrelated sibling caused settings reload")
	case <-time.After(100 * time.Millisecond):
	}
	if err = os.WriteFile(path, []byte("{\"scale\":2}"), 0600); err != nil {
		t.Fatal(err)
	}
	select {
	case <-called:
	case <-time.After(time.Second):
		t.Fatal("real settings change was ignored")
	}
}
