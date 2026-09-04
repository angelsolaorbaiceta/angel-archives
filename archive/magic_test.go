package archive

import (
	"bytes"
	"testing"
)

func TestReadMagic(t *testing.T) {
	t.Run("Not enough bytes to read", func(t *testing.T) {
		reader := bytes.NewBuffer([]byte{0x1})
		if got := ReadMagic(reader); got != FileMagicUnknown {
			t.Fatalf("Want Unknown, got %s", got)
		}
	})

	t.Run("Unknown magic", func(t *testing.T) {
		reader := bytes.NewBuffer([]byte{0x1, 0x2, 0x3, 0x4})
		if got := ReadMagic(reader); got != FileMagicUnknown {
			t.Fatalf("Want Unknown, got %s", got)
		}
	})

	t.Run("Archive magic", func(t *testing.T) {
		reader := bytes.NewBuffer(magic)
		if got := ReadMagic(reader); got != FileMagicArchive {
			t.Fatalf("Want Archive, got %s", got)
		}
	})

	t.Run("Encrypted archive magic", func(t *testing.T) {
		reader := bytes.NewBuffer(encMagic)
		if got := ReadMagic(reader); got != FileMagicEncArchive {
			t.Fatalf("Want EncArchive, got %s", got)
		}
	})
}
