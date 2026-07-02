# Proyecto_AED_7 — Radix / Patricia Trie

Proyecto del curso **Algoritmos y Estructura de Datos** (ESAN, 2026-1).
Estructura asignada (N.º 7): **Radix / Patricia Trie**.

> Paper de referencia: Morrison (1968), *"PATRICIA — Practical Algorithm To
> Retrieve Information Coded In Alphanumeric"*, JACM.

**Caso de uso:** autocompletado de nombres de ciudades del Perú conforme se
teclea un prefijo, con los datos cargados desde una base de datos real.

---

## ¿Qué es un Radix / Patricia Trie?

Un **trie** es un árbol que indexa cadenas carácter por carácter. Un **Radix
Trie (PATRICIA)** es un trie **comprimido**: cuando una cadena de nodos tiene un
solo hijo, esos nodos se fusionan en una sola arista etiquetada con varias
letras. Eso ahorra memoria y acelera la búsqueda.

```
Trie clásico para {casa, caso}:        Radix Trie equivalente:

        (raíz)                              (raíz)
          |c                                  |"cas"
          a                                  (•)
          |s                              a /     \ o
          (•)                            (casa)  (caso)
        a/    \o
   (casa)    (caso)
```

Operaciones (con `L` = longitud de la palabra):

| Operación | Qué hace | Complejidad |
|-----------|----------|-------------|
| `Insert(word)` | Inserta una palabra, partiendo aristas cuando hace falta | `O(L)` |
| `Search(word)` | Indica si la palabra exacta existe | `O(L)` |
| `Delete(word)` | Elimina una palabra y re-comprime el árbol | `O(L)` |
| `Autocomplete(prefix)` | Devuelve las palabras que empiezan con el prefijo | `O(P + K)` |
| `Keys()` | Devuelve todas las palabras ordenadas | `O(N·L)` |

`P` = longitud del prefijo, `K` = caracteres recorridos para reunir las
coincidencias, `N` = número de palabras.

---

## Estado de los entregables

- [x] **Entregable 2 — Codificación en Go**
  - [x] Estructura implementada desde cero en `trie/` (`Insert`, `Search`, `Delete`, `Autocomplete`, `Keys`)
  - [x] Pruebas unitarias (`trie/trie_test.go`)
  - [x] Benchmarks (`trie/trie_bench_test.go`)
  - [x] Repositorio con historial de commits y este README
- [x] **Entregable 3 — Aplicación con base de datos real**
  - [x] Conexión Go ↔ Supabase (PostgreSQL) por su API REST sobre HTTPS
  - [x] Dataset real de ciudades del Perú cargado en la BD
  - [x] El trie resuelve autocompletado sobre esos datos (`cmd/ciudades`)
- [x] **Entregable 4 — Simulación con Go + Vue.js**
  - [x] Backend en Go que expone el trie como API (`cmd/server`)
  - [x] Frontend en Vue.js con autocompletado de ciudades y árbol interactivo (`frontend/`)
  - [x] Reutiliza el paquete `trie` del Entregable 2 (no lo reimplementa)
- [ ] **Entregable 1 — PPTX tipo clase + video explicativo**

---

## Estructura del repositorio

```
Proyecto_AED_7/
├── go.mod
├── .env.example              # plantilla de configuración (copiar a .env)
├── trie/                     # Entregable 2 — la estructura de datos
│   ├── trie.go               #   implementación del Radix Trie
│   ├── trie_test.go          #   pruebas unitarias
│   └── trie_bench_test.go    #   benchmarks
├── internal/db/              # Entregable 3 — capa de base de datos
│   ├── db.go                 #   conexión a Supabase por REST/HTTPS
│   └── seed.go               #   dataset de ciudades del Perú
├── cmd/
│   ├── demo/main.go          # demo del trie por consola (sin base de datos)
│   ├── ciudades/main.go      # app que carga las ciudades de la BD al trie
│   └── server/main.go        # Entregable 4 — API en Go que expone el trie
└── frontend/                 # Entregable 4 — interfaz en Vue.js (Vite)
    ├── index.html
    ├── package.json
    └── src/
        ├── App.vue           #   paneles: autocompletado + árbol interactivo
        └── components/
            └── TrieTree.vue  #   dibuja el árbol en SVG
```

El paquete `trie` es **independiente**: no sabe de bases de datos ni de
interfaces. Por eso puede reutilizarse tal cual en los entregables 3 y 4, sin
reimplementar la estructura.

---

## Cómo ejecutar

Requisitos: **Go 1.25+**.

### Pruebas y benchmarks (Entregable 2)

```bash
# Ejecutar las pruebas unitarias
go test ./...

# Ejecutar las pruebas del trie con detalle
go test ./trie/ -v

# Ejecutar los benchmarks
go test ./trie/ -bench=. -benchmem -run=^$
```

### Demo del trie por consola (sin base de datos)

```bash
go run ./cmd/demo
```

### App con base de datos real (Entregable 3)

1. Copia `.env.example` a `.env` y completa tus datos de Supabase:
   ```
   SUPABASE_URL=https://TU-PROYECTO.supabase.co
   SUPABASE_KEY=sb_secret_tu_llave_secreta
   ```
2. Crea la tabla en Supabase (SQL Editor):
   ```sql
   create table if not exists ciudades (
     id serial primary key,
     nombre text not null unique,
     departamento text
   );
   ```
3. Ejecuta la app (carga las ciudades y ofrece autocompletado):
   ```bash
   go run ./cmd/ciudades
   ```

Escribe un prefijo (por ejemplo `Cu`) y verás las ciudades que lo completan.

> El `.env` está en `.gitignore`: las credenciales nunca se suben al repositorio.

### Simulación web con Vue.js (Entregable 4)

Necesita **Node.js 18+** además de Go.

```bash
# 1. Instalar dependencias del frontend (una sola vez)
cd frontend
npm install

# 2. Compilar el frontend
npm run build
cd ..

# 3. Levantar el servidor (sirve la web y la API en el mismo puerto)
go run ./cmd/server
```

Luego abre **http://localhost:8080** en el navegador. Verás dos paneles:

- **Autocompletado de ciudades**: escribe un prefijo y consulta los datos reales de la BD.
- **Árbol interactivo**: inserta o elimina palabras y observa cómo el Radix Trie
  se comprime y divide sus nodos; al escribir, se resalta la rama correspondiente.

> Para desarrollo del frontend con recarga en caliente: en una terminal
> `go run ./cmd/server` y en otra `cd frontend && npm run dev` (Vite en el
> puerto 5173 redirige las llamadas `/api` al backend).

---

## Detalle técnico: la conexión a la base de datos

Se usa la **API REST de Supabase sobre HTTPS (puerto 443)** en lugar del driver
de PostgreSQL (puerto 5432), porque muchas redes bloquean el puerto de la base
de datos, mientras que el 443 casi siempre está disponible. El flujo es:

```
Supabase (tabla ciudades)  →  internal/db (HTTP)  →  Radix Trie  →  autocompletado
```

---

## Autor

- Julio Carrizales
