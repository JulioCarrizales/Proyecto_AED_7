# Proyecto_AED_7 — Radix / Patricia Trie

Proyecto grupal del curso **Algoritmos y Estructura de Datos** (ESAN, 2026-1).
Estructura asignada (N.º 7): **Radix / Patricia Trie**.

> Paper de referencia: Morrison (1968), *"PATRICIA — Practical Algorithm To
> Retrieve Information Coded In Alphanumeric"*, JACM.

**Caso de uso:** autocompletado de palabras conforme se teclea un prefijo.

---

## ¿Qué es un Radix / Patricia Trie?

Un **trie** es un árbol que indexa cadenas carácter por carácter: cada arista
representa una letra y cada camino de la raíz a un nodo "palabra" representa una
cadena almacenada.

Un **Radix Trie (PATRICIA)** es un trie **comprimido**: cuando una cadena de
nodos tiene un solo hijo, esos nodos se fusionan en una sola arista etiquetada
con varias letras. Eso ahorra memoria y acelera la búsqueda.

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

Operaciones principales (todas en `O(L)`, con `L` = longitud de la palabra):

| Operación | Qué hace |
|-----------|----------|
| `Insert(word)` | Inserta una palabra, partiendo aristas cuando hace falta |
| `Search(word)` | Indica si la palabra exacta existe |
| `Autocomplete(prefix)` | Devuelve todas las palabras que empiezan con el prefijo |
| `Keys()` | Devuelve todas las palabras ordenadas |

---

## Estado actual del proyecto (avance)

- [x] Esqueleto del repositorio en Go
- [x] **Entregable 2 (inicio):** estructura implementada desde cero en `trie/`
      con `Insert`, `Search`, `Autocomplete` y `Keys`
- [x] Pruebas unitarias (`trie/trie_test.go`) — **pasan**
- [x] Demo de consola en `cmd/demo`
- [ ] Benchmarks de las operaciones clave
- [ ] **Entregable 1:** PPTX tipo clase + video explicativo
- [ ] **Entregable 3:** app en Go conectada a una base de datos real
- [ ] **Entregable 4:** simulación interactiva (backend Go + frontend Vue.js)

---

## Estructura del repositorio

```
Proyecto_AED_7/
├── go.mod
├── README.md
├── trie/                 # Entregable 2 — la estructura de datos
│   ├── trie.go           #   implementación del Radix Trie
│   └── trie_test.go      #   pruebas unitarias
└── cmd/
    └── demo/
        └── main.go       # demo de consola del autocompletado
```

Más adelante se añadirán:

```
├── internal/db/          # Entregable 3 — carga del dataset desde la BD
├── cmd/server/           # Entregable 4 — API en Go que expone el trie
└── frontend/             # Entregable 4 — interfaz en Vue.js
```

---

## Cómo ejecutar

Requisitos: **Go 1.24+**.

```bash
# Ejecutar las pruebas
go test ./...

# Ejecutar la demo de autocompletado por consola
go run ./cmd/demo
```

En la demo, escribe un prefijo (por ejemplo `cas`) y verás las palabras que lo
completan.

---

## Análisis de complejidad (Big-O)

| Operación | Tiempo | Espacio |
|-----------|--------|---------|
| `Insert` | `O(L)` | `O(L)` en el peor caso (nueva rama) |
| `Search` | `O(L)` | `O(1)` |
| `Autocomplete` | `O(P + K)` | `O(K)` para el resultado |

`L` = longitud de la palabra, `P` = longitud del prefijo, `K` = caracteres
recorridos para reunir las coincidencias.

---

## Autor

- Julio Carrizales
