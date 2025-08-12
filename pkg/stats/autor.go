package stats

import (
	"fmt"
	"os"
	"strings"

	"github.com/inovacc/git-nerds/pkg/unixcompat"
)

// CommitsPerAuthor reproduce `git log --pretty=format:%an | sort | uniq -c | sort -nr`
func CommitsPerAuthor(b Backend) error {
	// Obtener autores
	out, err := b.Log("--pretty=format:%an")
	if err != nil {
		return err
	}

	lines := strings.Split(strings.TrimSpace(out), "\n")

	// Contar ocurrencias
	counts := make(map[string]int)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			counts[line]++
		}
	}

	// Ordenar por valor descendente
	sorted := unixcompat.SortByValue(counts)

	// Imprimir alineado
	for _, kv := range sorted {
		_, _ = fmt.Fprintf(os.Stdout, "%-25s %d\n", kv.Key, kv.Value)
	}

	return nil
}
