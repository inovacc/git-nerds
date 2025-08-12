package unixcompat

import (
	"log"
	"testing"
)

func TestUnixCompat(t *testing.T) {
	lines := []string{"foo", "bar", "foo", "baz"}

	log.Println("Uniq:", Uniq(lines))
	log.Println("Head 2:", Head(lines, 2))
	log.Println("Basename:", Basename("/path/to/file.txt"))
	log.Println("Reverse:", Reverse("abc"))
}
