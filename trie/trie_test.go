package trie

import (
	"reflect"
	"testing"
)

func TestInsertSearch(t *testing.T) {
	tr := New()
	palabras := []string{"casa", "casaca", "caso", "perro", "pera"}
	for _, w := range palabras {
		tr.Insert(w)
	}

	for _, w := range palabras {
		if !tr.Search(w) {
			t.Errorf("Search(%q) = false; se esperaba true", w)
		}
	}

	for _, w := range []string{"cas", "ca", "per", "gato", ""} {
		if tr.Search(w) {
			t.Errorf("Search(%q) = true; se esperaba false", w)
		}
	}

	if tr.Len() != len(palabras) {
		t.Errorf("Len() = %d; se esperaba %d", tr.Len(), len(palabras))
	}
}

func TestInsertDuplicado(t *testing.T) {
	tr := New()
	tr.Insert("hola")
	tr.Insert("hola")
	if tr.Len() != 1 {
		t.Errorf("Len() = %d tras insertar duplicado; se esperaba 1", tr.Len())
	}
}

func TestAutocomplete(t *testing.T) {
	tr := New()
	for _, w := range []string{"casa", "casaca", "caso", "castor", "perro", "pera"} {
		tr.Insert(w)
	}

	got := tr.Autocomplete("cas")
	want := []string{"casa", "casaca", "caso", "castor"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Autocomplete(\"cas\") = %v; se esperaba %v", got, want)
	}

	got = tr.Autocomplete("per")
	want = []string{"pera", "perro"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Autocomplete(\"per\") = %v; se esperaba %v", got, want)
	}

	if res := tr.Autocomplete("xyz"); res != nil {
		t.Errorf("Autocomplete(\"xyz\") = %v; se esperaba nil", res)
	}
}

func TestKeys(t *testing.T) {
	tr := New()
	for _, w := range []string{"sol", "sal", "luna"} {
		tr.Insert(w)
	}
	got := tr.Keys()
	want := []string{"luna", "sal", "sol"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Keys() = %v; se esperaba %v", got, want)
	}
}
