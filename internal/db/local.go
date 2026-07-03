package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite" // driver SQLite puro en Go (sin CGO)
)

// LoadSQLite lee las ciudades desde una base de datos SQLite local. Permite
// ejecutar la aplicación sin conexión a internet ni a Supabase.
func LoadSQLite(path string) ([]Ciudad, error) {
	conn, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("abriendo SQLite %q: %w", path, err)
	}
	defer conn.Close()

	rows, err := conn.Query("SELECT nombre, departamento FROM ciudades ORDER BY nombre")
	if err != nil {
		return nil, fmt.Errorf("consultando ciudades locales: %w", err)
	}
	defer rows.Close()

	var ciudades []Ciudad
	for rows.Next() {
		var c Ciudad
		var dep sql.NullString
		if err := rows.Scan(&c.Nombre, &dep); err != nil {
			return nil, fmt.Errorf("leyendo fila: %w", err)
		}
		c.Departamento = dep.String
		ciudades = append(ciudades, c)
	}
	return ciudades, rows.Err()
}
