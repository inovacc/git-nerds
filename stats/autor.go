package stats

import (
	"fmt"
	"sort"
	"strings"

	"github.com/dyammarcano/clonr/pkg/git"
)

// CommitsPerAuthor obtiene la cantidad de commits por autor
func CommitsPerAuthor(b git.Backend) (map[string]int, error) {
	// --format='%an' para obtener solo el nombre del autor
	out, err := b.Log("--pretty=format:%an")
	if err != nil {
		return nil, err
	}

	counts := make(map[string]int)
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			counts[line]++
		}
	}

	return counts, nil
}

// PrintCommitsPerAuthor imprime los commits por autor ordenados por cantidad
func PrintCommitsPerAuthor(b git.Backend) error {
	counts, err := CommitsPerAuthor(b)
	if err != nil {
		return err
	}

	type kv struct {
		Key   string
		Value int
	}
	var sorted []kv
	for k, v := range counts {
		sorted = append(sorted, kv{k, v})
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Value > sorted[j].Value
	})

	for _, kv := range sorted {
		fmt.Printf("%-25s %d\n", kv.Key, kv.Value)
	}
	return nil
}
