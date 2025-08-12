package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// main es el punto de entrada del programa.
func main() {
	// Reemplaza esta ruta con la ruta a tu repositorio local.
	// O, si lo prefieres, puedes pasarlo como un argumento de la línea de comandos.
	repoPath := "." // Asume el directorio actual.
	if len(os.Args) > 1 {
		repoPath = os.Args[1]
	}

	// Abre el repositorio de Git.
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		log.Fatalf("Error al abrir el repositorio: %s", err)
	}

	// Obtiene un iterador para todos los commits en el historial.
	commitIterator, err := repo.Log(&git.LogOptions{})
	if err != nil {
		log.Fatalf("Error al obtener el historial de commits: %s", err)
	}

	// --- Inicializa los contadores para las estadísticas ---
	commitsByUser := make(map[string]int)
	fileModifications := make(map[string]int)
	linesAdded := 0
	linesDeleted := 0
	commitsByWeekday := make(map[time.Weekday]int)

	// Itera sobre cada commit.
	commitIterator.ForEach(func(commit *object.Commit) error {
		// 1. Commits por usuario
		author := commit.Author.Email
		commitsByUser[author]++

		// 4. Días de la semana que más se hacen commits
		weekday := commit.Author.When.Weekday()
		commitsByWeekday[weekday]++

		// Si es el primer commit (sin padre), no hay cambios para comparar.
		if commit.NumParents() == 0 {
			return nil
		}

		// 2. Archivos más modificados y 3. Líneas de código
		// Obtenemos el padre del commit actual para comparar los cambios.
		parent, err := commit.Parent(0)
		if err != nil {
			return err
		}

		// Creamos un Patch entre el commit y su padre para analizar los cambios.
		patch, err := parent.Patch(commit)
		if err != nil {
			return err
		}

		// Iteramos sobre los FileStats para obtener las estadísticas de cada archivo.
		for _, fs := range patch.Stats() {
			fileModifications[fs.Name]++
			linesAdded += fs.Addition
			linesDeleted += fs.Deletion
		}

		return nil
	})

	// --- Muestra los resultados ---
	fmt.Println("--- Estadísticas del Repositorio ---")

	// 1. Commits por usuario
	fmt.Println("\nCommits por usuario:")

	for user, count := range commitsByUser {
		fmt.Printf("- %s: %d\n", user, count)
	}

	// 2. Archivos más modificados
	fmt.Println("\nArchivos más modificados (Top 5):")

	// Convertimos el mapa a un slice para poder ordenarlo.
	type fileStat struct {
		name  string
		count int
	}

	fileStats := make([]fileStat, 0, len(fileModifications))
	for name, count := range fileModifications {
		fileStats = append(fileStats, fileStat{name, count})
	}

	// Ordenamos el slice de forma descendente.
	sort.Slice(fileStats, func(i, j int) bool {
		return fileStats[i].count > fileStats[j].count
	})

	// Imprimimos el top 5 o todos si hay menos de 5.
	limit := 5
	if len(fileStats) < limit {
		limit = len(fileStats)
	}

	for i := 0; i < limit; i++ {
		fmt.Printf("- %s: %d modificaciones\n", fileStats[i].name, fileStats[i].count)
	}

	// 3. Líneas de código del repo
	fmt.Println("\nLíneas de código:")
	fmt.Printf("- Líneas añadidas: %d\n", linesAdded)
	fmt.Printf("- Líneas eliminadas: %d\n", linesDeleted)
	fmt.Printf("- Total de cambios: %d\n", linesAdded+linesDeleted)

	// 4. Días de la semana que más se hacen commits
	fmt.Println("\nCommits por día de la semana:")
	// Creamos un slice para ordenar los días.
	days := []time.Weekday{
		time.Sunday,
		time.Monday,
		time.Tuesday,
		time.Wednesday,
		time.Thursday,
		time.Friday,
		time.Saturday,
	}
	for _, day := range days {
		count := commitsByWeekday[day]
		fmt.Printf("- %s: %d commits\n", day, count)
	}
}
