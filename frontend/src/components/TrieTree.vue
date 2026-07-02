<script setup>
import { computed } from 'vue'

// Recibe el árbol (raíz NodeView del backend) y un texto a resaltar.
const props = defineProps({
  tree: { type: Object, default: null },
  highlight: { type: String, default: '' }
})

const DX = 64 // separación horizontal entre hojas
const DY = 78 // separación vertical entre niveles

// Calcula posiciones (x, y) de cada nodo con un recorrido en profundidad:
// las hojas se colocan de izquierda a derecha y cada nodo interno se centra
// sobre sus hijos.
const layout = computed(() => {
  const nodes = []
  const edges = []
  if (!props.tree) return { nodes, edges, width: 200, height: 120 }

  const q = (props.highlight || '').toLowerCase()
  let leaf = 0
  let idc = 0

  // Un nodo está "en el camino" del texto buscado si su ruta acumulada es
  // prefijo del texto, o el texto es prefijo de su ruta (resalta la rama).
  const enCamino = (acc) => q !== '' && (acc.startsWith(q) || q.startsWith(acc))

  const walk = (node, depth, acc) => {
    const id = idc++
    const hijos = node.children || []
    const xsHijos = []
    for (const c of hijos) {
      const hijo = walk(c, depth + 1, acc + c.label)
      xsHijos.push(hijo.x)
      edges.push({ from: id, to: hijo.id, label: c.label, hi: enCamino(acc + c.label) })
    }
    let x
    if (xsHijos.length === 0) {
      x = leaf * DX + DX
      leaf++
    } else {
      x = (xsHijos[0] + xsHijos[xsHijos.length - 1]) / 2
    }
    const rec = { id, x, y: depth * DY + 40, label: node.label, isWord: node.isWord, hi: enCamino(acc), depth }
    nodes.push(rec)
    return rec
  }
  walk(props.tree, 0, '')

  const byId = {}
  for (const n of nodes) byId[n.id] = n
  for (const e of edges) {
    e.x1 = byId[e.from].x
    e.y1 = byId[e.from].y
    e.x2 = byId[e.to].x
    e.y2 = byId[e.to].y
  }

  const width = Math.max(220, leaf * DX + DX)
  const maxDepth = nodes.reduce((m, n) => Math.max(m, n.depth), 0)
  const height = (maxDepth + 1) * DY + 20
  return { nodes, edges, width, height }
})
</script>

<template>
  <svg
    :viewBox="`0 0 ${layout.width} ${layout.height}`"
    :width="layout.width"
    :height="layout.height"
    class="trie-svg"
  >
    <!-- aristas -->
    <line
      v-for="(e, i) in layout.edges"
      :key="'e' + i"
      :x1="e.x1" :y1="e.y1" :x2="e.x2" :y2="e.y2"
      :class="['edge', { hi: e.hi }]"
    />
    <!-- etiquetas de las aristas (las letras) -->
    <text
      v-for="(e, i) in layout.edges"
      :key="'t' + i"
      :x="(e.x1 + e.x2) / 2"
      :y="(e.y1 + e.y2) / 2"
      :class="['edge-label', { hi: e.hi }]"
      text-anchor="middle"
      dy="-3"
    >{{ e.label }}</text>
    <!-- nodos -->
    <circle
      v-for="n in layout.nodes"
      :key="'n' + n.id"
      :cx="n.x" :cy="n.y" r="13"
      :class="['node', { word: n.isWord, hi: n.hi }]"
    />
  </svg>
</template>

<style scoped>
.trie-svg {
  display: block;
}
.edge {
  stroke: #cbd5e1;
  stroke-width: 2;
}
.edge.hi {
  stroke: var(--highlight);
  stroke-width: 3;
}
.edge-label {
  fill: #475569;
  font-size: 13px;
  font-family: "Consolas", monospace;
}
.edge-label.hi {
  fill: #b8860b;
  font-weight: 700;
}
.node {
  fill: #ffffff;
  stroke: #94a3b8;
  stroke-width: 2;
}
.node.word {
  fill: var(--node-word);
  stroke: var(--node-word);
}
.node.hi {
  stroke: var(--highlight);
  stroke-width: 3;
}
</style>
