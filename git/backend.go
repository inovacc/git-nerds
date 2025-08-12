package git

// Backend define la interfaz para consultar datos de Git.
// Así en el futuro podemos añadir otros backends (go-git, mocks).
type Backend interface {
	Log(args ...string) (string, error)
	Branches(args ...string) (string, error)
}
