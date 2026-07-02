<script setup>
import { ref, onMounted } from 'vue'
import TrieTree from './components/TrieTree.vue'

// --- Panel 1: autocompletado de ciudades (datos reales desde la BD) ---
const cityPrefix = ref('')
const cityMatches = ref([])
const cityLoaded = ref(false)

async function buscarCiudades() {
  const p = cityPrefix.value.trim()
  if (p === '') {
    cityMatches.value = []
    cityLoaded.value = false
    return
  }
  const res = await fetch(`/api/cities?prefix=${encodeURIComponent(p)}`)
  const data = await res.json()
  cityMatches.value = data.matches || []
  cityLoaded.value = true
}

// --- Panel 2: árbol interactivo (se construye en vivo) ---
const tree = ref(null)
const newWord = ref('')
const highlight = ref('')

async function cargarArbol() {
  const res = await fetch('/api/tree')
  tree.value = await res.json()
}

async function insertar() {
  const w = newWord.value.trim()
  if (w === '') return
  const res = await fetch(`/api/tree/insert?word=${encodeURIComponent(w)}`, { method: 'POST' })
  tree.value = await res.json()
  highlight.value = w.toLowerCase()
  newWord.value = ''
}

async function eliminar() {
  const w = newWord.value.trim()
  if (w === '') return
  const res = await fetch(`/api/tree/delete?word=${encodeURIComponent(w)}`, { method: 'POST' })
  tree.value = await res.json()
  highlight.value = ''
  newWord.value = ''
}

async function reiniciar() {
  const res = await fetch('/api/tree/reset', { method: 'POST' })
  tree.value = await res.json()
  highlight.value = ''
}

onMounted(cargarArbol)
</script>

<template>
  <header class="topbar">
    <img class="logo" src="/logo-esan.png" alt="Universidad ESAN" />
    <div>
      <h1>Radix Trie — Simulación interactiva</h1>
      <p>Proyecto AED 7 · backend en Go + frontend en Vue.js</p>
    </div>
  </header>

  <div class="grid">
    <!-- Panel 1 -->
    <section class="card">
      <h2>Autocompletado de ciudades</h2>
      <p class="hint">Datos reales cargados desde la base de datos (Supabase).</p>
      <input
        type="text"
        v-model="cityPrefix"
        @input="buscarCiudades"
        placeholder="Escribe un prefijo, p. ej. Cu"
      />
      <ul class="matches" v-if="cityMatches.length">
        <li v-for="c in cityMatches" :key="c.nombre">
          {{ c.nombre }}<span class="dep" v-if="c.departamento"> · {{ c.departamento }}</span>
        </li>
      </ul>
      <p class="empty" v-else-if="cityLoaded">Sin coincidencias para ese prefijo.</p>
      <p class="empty" v-else>Empieza a escribir para ver las sugerencias.</p>
    </section>

    <!-- Panel 2 -->
    <section class="card">
      <h2>Árbol interactivo</h2>
      <p class="hint">Inserta o elimina palabras y observa cómo el árbol se comprime y divide sus nodos.</p>
      <div class="row">
        <input
          type="text"
          v-model="newWord"
          @keyup.enter="insertar"
          @input="highlight = newWord.toLowerCase()"
          placeholder="Palabra a insertar o eliminar"
        />
        <button class="primary" @click="insertar">Insertar</button>
        <button class="danger" @click="eliminar">Eliminar</button>
        <button @click="reiniciar">Reiniciar</button>
      </div>

      <div class="tree-wrap">
        <TrieTree :tree="tree" :highlight="highlight" />
      </div>

      <div class="legend">
        <span><i class="dot" style="background:#d8231f"></i> nodo que termina una palabra</span>
        <span><i class="dot" style="background:#fff;border:2px solid #94a3b8"></i> nodo intermedio</span>
        <span><i class="dot" style="background:#e8a400"></i> resaltado por el texto</span>
      </div>
    </section>
  </div>
</template>
