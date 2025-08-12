package unixcompat

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/fatih/color"
)

// Grep filtra líneas que contengan un patrón (substring o regex)
func Grep(lines []string, pattern string, useRegex bool, ignoreCase bool) []string {
	var result []string

	if useRegex {
		flags := ""
		if ignoreCase {
			flags = "(?i)"
		}

		re := regexp.MustCompile(flags + pattern)
		for _, l := range lines {
			if re.MatchString(l) {
				result = append(result, l)
			}
		}
	} else {
		if ignoreCase {
			pattern = strings.ToLower(pattern)
			for _, l := range lines {
				if strings.Contains(strings.ToLower(l), pattern) {
					result = append(result, l)
				}
			}
		} else {
			for _, l := range lines {
				if strings.Contains(l, pattern) {
					result = append(result, l)
				}
			}
		}
	}

	return result
}

// Fields simula `awk` separando por espacios en blanco
func Fields(line string) []string {
	return strings.Fields(line)
}

// SplitN simula `cut o `awk -F` con delimitador y máximo de partes
func SplitN(line, sep string, n int) []string {
	return strings.SplitN(line, sep, n)
}

// SortStrings ordena un slice de strings ascendente
func SortStrings(data []string) {
	sort.Strings(data)
}

type KeyValue struct {
	Key   string
	Value int
}

// SortByValue ordena un mapa[string]int por valor descendente
func SortByValue(m map[string]int) []KeyValue {
	var kvs = make([]KeyValue, 0)

	for k, v := range m {
		kvs = append(kvs, KeyValue{Key: k, Value: v})
	}

	sort.Slice(kvs, func(i, j int) bool {
		return kvs[i].Value > kvs[j].Value
	})

	return kvs
}

// Uniq elimina duplicados preservando orden
func Uniq(lines []string) []string {
	seen := make(map[string]struct{})

	var result []string

	for _, l := range lines {
		if _, ok := seen[l]; !ok {
			seen[l] = struct{}{}
			result = append(result, l)
		}
	}

	return result
}

// Head devuelve las primeras N líneas
func Head(lines []string, n int) []string {
	if n > len(lines) {
		return lines
	}

	return lines[:n]
}

// Basename obtiene el último elemento de una ruta
func Basename(path string) string {
	return filepath.Base(path)
}

// ReplaceAll simula `tr` de reemplazo simple
func ReplaceAll(s, o, n string) string {
	return strings.ReplaceAll(s, o, n)
}

// ToLower convierte a minúsculas
func ToLower(s string) string {
	return strings.ToLower(s)
}

// Reverse invierte una cadena (similar a `rev`)
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}

// DateNow devuelve fecha actual con formato (similar a `date +"%Y-%m-%d"`)
func DateNow(format string) string {
	return time.Now().Format(format)
}

// Printf es un wrapper de fmt.Printf
func Printf(format string, args ...any) {
	_, _ = fmt.Fprintf(os.Stdout, format, args...)
}

// PrintColor imprime texto con color usando fatih/color (simula `tput setaf`)
func PrintColor(c *color.Color, format string, args ...any) {
	_, _ = c.Printf(format, args...)
}

// ReadLines lee un archivo línea por línea
func ReadLines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	var lines []string

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}
