package model

// Access модель доступа в репозитории
type Access struct {
	Endpoint string   `db:"endpoint"`
	Roles    []string `db:"roles"`
}
