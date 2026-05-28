<script setup>
import { reactive, ref, onMounted, onUnmounted, nextTick } from 'vue'
import { ConvertToMd } from '../../wailsjs/go/main/App'
import { BrowserOpenURL } from '../../wailsjs/runtime/runtime'
import openFileSvg from '../assets/open-file.svg?raw'
import searchSvg from '../assets/search.svg?raw'
import zoomInSvg from '../assets/zoom-in.svg?raw'
import zoomOutSvg from '../assets/zoom-out.svg?raw'
import infoSvg from '../assets/info.svg?raw'
import Mark from 'mark.js';

const data = reactive({
  content: null
})


const searchInstance = ref(null)
const searchOpen = ref(false)
const loading = ref(false)
const footerExpanded = ref(false)
let collapseTimer = null

const COLLAPSE_DELAY = 10000

function handleLinkClick(e) {
  const anchor = e.target.closest('a')
  if (anchor && anchor.href) {
    e.preventDefault()
    BrowserOpenURL(anchor.href)
  }
}

function clearCollapseTimer() {
  if (collapseTimer !== null) {
    clearTimeout(collapseTimer)
    collapseTimer = null
  }
}

function scheduleCollapse() {
  clearCollapseTimer()
  collapseTimer = setTimeout(() => {
    footerExpanded.value = false
  }, COLLAPSE_DELAY)
}

function instantCollapse() {
  clearCollapseTimer()
  footerExpanded.value = false
}

function expandFooter() {
  footerExpanded.value = true
  scheduleCollapse()
}


// max zoom level is 2.0, min zoom level is 0.3
function ZoomIn() {
  const content = document.getElementById('content')
  if (content) {
    const currentZoom = parseFloat(content.style.zoom) || 1
    content.style.zoom = Math.min(currentZoom + 0.03, 2.0).toString()
  }
}


function ZoomOut() {
  const content = document.getElementById('content')
  if (content) {
    const currentZoom = parseFloat(content.style.zoom) || 1
    content.style.zoom = Math.max(currentZoom - 0.03, 0.3).toString()
  }
}


onMounted(() => {
  document.querySelector('main')?.addEventListener('click', handleLinkClick)
  window.addEventListener('keydown', handleKeydown, true)
})

onUnmounted(() => {
  document.querySelector('main')?.removeEventListener('click', handleLinkClick)
  window.removeEventListener('keydown', handleKeydown, true)
  clearCollapseTimer()
})

async function loadFile() {
  clearCollapseTimer()
  loading.value = true
  try {
    data.content = null
    const file = await ConvertToMd()
    if (file) {
      //console.dir(file)
      data.content = file.html
      await nextTick()
      const contentEl = document.getElementById('content')
      if (contentEl) {
        searchInstance.value = new Mark(contentEl)
      }
      footerExpanded.value = true
      scheduleCollapse()
    }
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

function handleKeydown(e) {
  const ctrl = e.ctrlKey || e.metaKey;
  if (ctrl && e.key.toLowerCase() === 'o') { e.preventDefault(); loadFile() }
  if (e.key === 'F3') { e.preventDefault(); openSearch() }
  if (e.key === 'F4') { e.preventDefault(); ZoomIn() }
  if (e.key === 'F5') { e.preventDefault(); ZoomOut() }
  if (e.key === 'F1') { e.preventDefault(); console.log('About Margo: A Markdown Viewer built with Wails and Vue.js') }
  if (e.key === 'Escape' && searchOpen.value) { e.preventDefault(); closeSearch() }
}

function find(term) {
  if (searchInstance.value) {
    searchInstance.value.unmark({
      done: () => {
        searchInstance.value.mark(term)
      }
    })
  }
}

function openSearch() {
  footerExpanded.value = false
  searchOpen.value = true
  const searchPanel = document.getElementById('search-panel')
  if (searchPanel) {
    searchPanel.style.display = 'flex'
    const input = searchPanel.querySelector('input')
    if (input) {
      input.focus()
    }
  }
}

function closeSearch() {
  searchOpen.value = false
  const searchPanel = document.getElementById('search-panel')
  if (searchPanel) {
    searchPanel.style.display = 'none'
    if (searchInstance.value) {
      searchInstance.value.unmark()
    }
  }
}

</script>

<template>
  <main>
    <div v-if="data.content">
      <div id="content" @click="instantCollapse" v-html="data.content"></div>
      <div id="footer" class="grid" :class="{ collapsed: !footerExpanded }" @click="scheduleCollapse"
        @mouseenter="scheduleCollapse" @mouseleave="scheduleCollapse" @mousemove="scheduleCollapse">
        <div class="card">
          <div class="item" @click="loadFile" v-html="openFileSvg" title="Open File - Ctrl+O"></div>
        </div>
        <div class="card">
          <div class="item" @click="openSearch" v-html="searchSvg" title="Search - F3"></div>
        </div>
        <div class="card">
          <div class="item" @click="ZoomIn" v-html="zoomInSvg" title="Zoom In - F4"></div>
        </div>
        <div class="card">
          <div class="item" @click="ZoomOut" v-html="zoomOutSvg" title="Zoom Out - F5"></div>
        </div>
        <div class="card">
          <div class="item" v-html="infoSvg" title="About - F1"></div>
        </div>
      </div>
      <div class="hamburger" :class="{ visible: !footerExpanded }" @click="expandFooter">☰</div>
      <div id="search-panel">
        <span>Look for:</span>
        <input type="text" @input="find($event.target.value)" />
        <span class="close-btn" @click="closeSearch">×</span>
      </div>
    </div>
    <div id="no-content" v-else>
      <p v-if="loading">Loading markdown...</p>
      <button id="load-md" @click="loadFile" v-else>Load Markdown</button>
    </div>
  </main>
</template>
