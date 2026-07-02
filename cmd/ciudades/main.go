// Entregable 3: aplicación que conecta Go con una base de datos real
// (Supabase / PostgreSQL vía su API REST), carga un dataset de ciudades del
// Perú y las indexa en el Radix Trie para resolver autocompletado por prefijo.
package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/JulioCarrizales/Proyecto_AED_7/internal/db"
	"github.com/JulioCarrizales/Proyecto_AED_7/trie"
)

func main() {
	loadEnv(".env")

	url := os.Getenv("SUPABASE_URL")
	key := os.Getenv("SUPABASE_KEY")
	if url == "" || key == "" {
		log.Fatal("Faltan SUPABASE_URL o SUPABASE_KEY. Copia .env.example a .env y complétalos.")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 1. Conexión Go <-> base de datos real (Supabase por HTTPS).
	fmt.Println("Conectando a la base de datos...")
	store, err := db.Connect(ctx, url, key)
	if err != nil {
		log.Fatalf("No se pudo conectar: %v", err)
	}
	defer store.Close()
	fmt.Println("Conexión exitosa.")

	// 2. Cargar el dataset si la tabla está vacía.
	n, err := store.Count(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if n == 0 {
		fmt.Println("Tabla vacía: cargando dataset de ciudades del Perú...")
		insertadas, err := store.Seed(ctx, db.Dataset())
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Se insertaron %d ciudades en la base de datos.\n", insertadas)
	}

	// 3. Cargar las ciudades DESDE la BD hacia el Radix Trie.
	//    Cada ciudad (clave) guarda su departamento como valor genérico.
	ciudades, err := store.ListCiudades(ctx)
	if err != nil {
		log.Fatal(err)
	}
	t := trie.New[string]()
	for _, c := range ciudades {
		t.Insert(c.Nombre, c.Departamento)
	}
	fmt.Printf("Radix Trie cargado con %d ciudades desde la base de datos.\n\n", t.Len())

	// 4. La estructura resolviendo consultas reales: autocompletado por prefijo.
	fmt.Println("Escribe un prefijo para autocompletar ciudades (vacío para salir):")
	sc := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for sc.Scan() {
		prefijo := strings.TrimSpace(sc.Text())
		if prefijo == "" {
			break
		}
		sugerencias := t.Autocomplete(prefijo)
		if len(sugerencias) == 0 {
			fmt.Println("  (sin coincidencias)")
		} else {
			for _, m := range sugerencias {
				fmt.Printf("  %s (%s)\n", m.Key, m.Value)
			}
		}
		fmt.Print("> ")
	}
}

// loadEnv lee un archivo .env sencillo (CLAVE=valor) y lo carga en el entorno.
func loadEnv(path string) {
	f, err := os.Open(path)
	if err != nil {
		return // sin .env se usan las variables del entorno del sistema
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
