// Package db maneja la conexión a la base de datos real (Supabase / PostgreSQL)
// a través de su API REST sobre HTTPS. Es el núcleo del Entregable 3: demuestra
// la conexión Go<->BD, la carga de un dataset real y la estructura resolviendo
// consultas concretas (autocompletado).
//
// Se usa la API REST (puerto 443) en lugar del driver de PostgreSQL (puerto
// 5432) porque muchas redes bloquean el puerto de la base de datos; el 443
// (HTTPS) casi siempre está disponible.
package db

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Ciudad es un registro de la tabla "ciudades".
type Ciudad struct {
	Nombre       string `json:"nombre"`
	Departamento string `json:"departamento"`
}

// Store habla con Supabase por su API REST.
type Store struct {
	baseURL string
	key     string
	http    *http.Client
}

// Connect prepara el cliente y verifica que la BD responda con una consulta mínima.
func Connect(ctx context.Context, url, key string) (*Store, error) {
	if url == "" || key == "" {
		return nil, fmt.Errorf("faltan SUPABASE_URL o SUPABASE_KEY en el .env")
	}
	s := &Store{
		baseURL: strings.TrimRight(url, "/") + "/rest/v1",
		key:     key,
		http:    &http.Client{Timeout: 15 * time.Second},
	}
	if _, err := s.Count(ctx); err != nil {
		return nil, fmt.Errorf("no se pudo verificar la conexión: %w", err)
	}
	return s, nil
}

// Close se mantiene por simetría; con HTTP no hay conexión persistente que cerrar.
func (s *Store) Close() {}

// newRequest arma una petición con los encabezados de autenticación de Supabase.
func (s *Store) newRequest(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, s.baseURL+path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("apikey", s.key)
	req.Header.Set("Authorization", "Bearer "+s.key)
	return req, nil
}

// Count devuelve cuántas ciudades hay almacenadas, leyendo el encabezado
// Content-Range que PostgREST incluye cuando se pide count=exact.
func (s *Store) Count(ctx context.Context) (int, error) {
	req, err := s.newRequest(ctx, http.MethodGet, "/ciudades?select=id", nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Prefer", "count=exact")
	req.Header.Set("Range", "0-0")

	resp, err := s.http.Do(req)
	if err != nil {
		return 0, fmt.Errorf("consultando la BD: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
		b, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("respuesta inesperada (%d): %s", resp.StatusCode, strings.TrimSpace(string(b)))
	}

	// Content-Range tiene forma "0-0/115" o "*/115"; el total va tras la barra.
	cr := resp.Header.Get("Content-Range")
	_, total, ok := strings.Cut(cr, "/")
	if !ok {
		return 0, nil
	}
	n, err := strconv.Atoi(total)
	if err != nil {
		return 0, nil
	}
	return n, nil
}

// Seed inserta las ciudades dadas, ignorando las que ya existan (por nombre).
// Devuelve cuántas filas nuevas se insertaron.
func (s *Store) Seed(ctx context.Context, ciudades []Ciudad) (int, error) {
	body, err := json.Marshal(ciudades)
	if err != nil {
		return 0, err
	}
	req, err := s.newRequest(ctx, http.MethodPost, "/ciudades?on_conflict=nombre", bytes.NewReader(body))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	// ignore-duplicates: no falla si el nombre ya existe.
	// return=representation: la respuesta trae las filas realmente insertadas.
	req.Header.Set("Prefer", "resolution=ignore-duplicates,return=representation")

	resp, err := s.http.Do(req)
	if err != nil {
		return 0, fmt.Errorf("insertando ciudades: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("respuesta inesperada al insertar (%d): %s", resp.StatusCode, strings.TrimSpace(string(b)))
	}

	var insertadas []Ciudad
	if err := json.NewDecoder(resp.Body).Decode(&insertadas); err != nil {
		return 0, fmt.Errorf("leyendo respuesta de inserción: %w", err)
	}
	return len(insertadas), nil
}

// ListNombres devuelve todos los nombres de ciudad ordenados. Es lo que se
// carga en el Radix Trie para el autocompletado.
func (s *Store) ListNombres(ctx context.Context) ([]string, error) {
	req, err := s.newRequest(ctx, http.MethodGet, "/ciudades?select=nombre&order=nombre.asc", nil)
	if err != nil {
		return nil, err
	}
	resp, err := s.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("consultando nombres: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("respuesta inesperada (%d): %s", resp.StatusCode, strings.TrimSpace(string(b)))
	}

	var filas []struct {
		Nombre string `json:"nombre"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&filas); err != nil {
		return nil, fmt.Errorf("leyendo nombres: %w", err)
	}

	nombres := make([]string, 0, len(filas))
	for _, f := range filas {
		nombres = append(nombres, f.Nombre)
	}
	return nombres, nil
}
