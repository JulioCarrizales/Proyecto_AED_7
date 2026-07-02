package trie

import (
	"strings"
	"testing"
)

// benchWords genera un conjunto determinista de ~2000 palabras para medir el
// rendimiento de las operaciones del trie de forma reproducible.
func benchWords() []string {
	const letras = "abcdefghijklmnopqrstuvwxyz"
	palabras := make([]string, 0, 2000)
	for i := 0; i < 2000; i++ {
		largo := 3 + i%6
		var sb strings.Builder
		x := i
		for j := 0; j < largo; j++ {
			sb.WriteByte(letras[x%26])
			x = x*31 + 7
		}
		palabras = append(palabras, sb.String())
	}
	return palabras
}

func BenchmarkInsert(b *testing.B) {
	palabras := benchWords()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tr := New()
		for _, w := range palabras {
			tr.Insert(w)
		}
	}
}

func BenchmarkSearch(b *testing.B) {
	palabras := benchWords()
	tr := New()
	for _, w := range palabras {
		tr.Insert(w)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tr.Search(palabras[i%len(palabras)])
	}
}

func BenchmarkAutocomplete(b *testing.B) {
	palabras := benchWords()
	tr := New()
	for _, w := range palabras {
		tr.Insert(w)
	}
	prefijos := []string{"ab", "ca", "lo", "pe", "ma", "za"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tr.Autocomplete(prefijos[i%len(prefijos)])
	}
}

func BenchmarkDelete(b *testing.B) {
	palabras := benchWords()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		tr := New()
		for _, w := range palabras {
			tr.Insert(w)
		}
		b.StartTimer()
		for _, w := range palabras {
			tr.Delete(w)
		}
	}
}
