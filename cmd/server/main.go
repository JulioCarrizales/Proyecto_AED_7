// Entregable 4 (backend): servidor HTTP en Go que expone el Radix Trie como una
// API para que el frontend en Vue.js la consuma. Reutiliza el paquete trie del
// Entregable 2 (no reimplementa la estructura) y las ciudades del Entregable 3.
package main

import (
	"bufio"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/JulioCarrizales/Proyecto_AED_7/internal/db"
	"github.com/JulioCarrizales/Proyecto_AED_7/trie"
)

// cities: trie de solo lectura con las ciudades cargadas desde la base de datos.
// viz: trie que el usuario construye en vivo desde la web (protegido por vizMu).
var (
	cities *trie.Trie
	viz    *trie.Trie
	vizMu  sync.Mutex
)

// vizSeed son las palabras iniciales del árbol interactivo: comparten prefijos
// para que se vea con claridad la compresión y la división de nodos.
var vizSeed = []string{"casa", "caso", "castor", "catedral", "perro", "pera", "persona", "sol", "solar"}

func main() {
	loadEnv(".env")

	cities = trie.New()
	cargarCiudades()

	viz = trie.New()
	resetViz()

	mux := http.NewServeMux()
	mux.HandleFunc("/api/cities", handleCities)
	mux.HandleFunc("/api/tree", handleTree)
	mux.HandleFunc("/api/tree/insert", handleTreeInsert)
	mux.HandleFunc("/api/tree/delete", handleTreeDelete)
	mux.HandleFunc("/api/tree/reset", handleTreeReset)

	// Si el frontend ya está compilado, se sirve desde el mismo servidor.
	if _, err := os.Stat("frontend/dist"); err == nil {
		mux.Handle("/", http.FileServer(http.Dir("frontend/dist")))
	}

	addr := ":8080"
	log.Printf("Servidor escuchando en http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, withCORS(mux)))
}

// cargarCiudades trae las ciudades desde Supabase hacia el trie de solo lectura.
func cargarCiudades() {
	url, key := os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_KEY")
	if url == "" || key == "" {
		log.Print("Aviso: sin SUPABASE_URL/KEY; el autocompletado de ciudades quedará vacío.")
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	store, err := db.Connect(ctx, url, key)
	if err != nil {
		log.Printf("Aviso: no se pudo conectar a la BD (%v); ciudades vacías.", err)
		return
	}
	defer store.Close()

	nombres, err := store.ListNombres(ctx)
	if err != nil {
		log.Printf("Aviso: no se pudieron cargar las ciudades: %v", err)
		return
	}
	for _, n := range nombres {
		cities.Insert(n)
	}
	log.Printf("Cargadas %d ciudades desde la base de datos.", cities.Len())
}

func resetViz() {
	viz = trie.New()
	for _, w := range vizSeed {
		viz.Insert(w)
	}
}

// handleCities: autocompletado de ciudades reales (Entregable 3 sobre HTTP).
func handleCities(w http.ResponseWriter, r *http.Request) {
	prefijo := strings.TrimSpace(r.URL.Query().Get("prefix"))
	var coincidencias []string
	if prefijo != "" {
		coincidencias = cities.Autocomplete(prefijo)
		if len(coincidencias) > 25 {
			coincidencias = coincidencias[:25]
		}
	}
	writeJSON(w, map[string]any{"prefix": prefijo, "matches": coincidencias})
}

// handleTree devuelve el árbol interactivo actual como JSON.
func handleTree(w http.ResponseWriter, r *http.Request) {
	vizMu.Lock()
	defer vizMu.Unlock()
	writeJSON(w, viz.View())
}

// handleTreeInsert inserta una palabra en el árbol interactivo y lo devuelve.
func handleTreeInsert(w http.ResponseWriter, r *http.Request) {
	palabra := normalizar(r.URL.Query().Get("word"))
	vizMu.Lock()
	defer vizMu.Unlock()
	if palabra != "" {
		viz.Insert(palabra)
	}
	writeJSON(w, viz.View())
}

// handleTreeDelete elimina una palabra del árbol interactivo y lo devuelve.
func handleTreeDelete(w http.ResponseWriter, r *http.Request) {
	palabra := normalizar(r.URL.Query().Get("word"))
	vizMu.Lock()
	defer vizMu.Unlock()
	if palabra != "" {
		viz.Delete(palabra)
	}
	writeJSON(w, viz.View())
}

// handleTreeReset restablece el árbol interactivo a sus palabras iniciales.
func handleTreeReset(w http.ResponseWriter, r *http.Request) {
	vizMu.Lock()
	defer vizMu.Unlock()
	resetViz()
	writeJSON(w, viz.View())
}

// normalizar deja la palabra en minúsculas y sin espacios alrededor.
func normalizar(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// withCORS permite que el frontend de Vite (otro puerto en desarrollo) consuma la API.
func withCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		h.ServeHTTP(w, r)
	})
}

// loadEnv lee un archivo .env sencillo (CLAVE=valor) y lo carga en el entorno.
func loadEnv(path string) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		clave, valor, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		os.Setenv(strings.TrimSpace(clave), strings.TrimSpace(valor))
	}
}
