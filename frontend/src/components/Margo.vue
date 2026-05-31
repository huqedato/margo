<script setup>
import { reactive, ref, onMounted, onUnmounted, nextTick } from 'vue'
import { OpenFileFromFrontend, GetStartupFile, GetAppInfo } from '../../wailsjs/go/main/App'
import { BrowserOpenURL } from '../../wailsjs/runtime/runtime'
import openFileSvg from '../assets/open-file.svg?raw'
import searchSvg from '../assets/search.svg?raw'
import zoomInSvg from '../assets/zoom-in.svg?raw'
import zoomOutSvg from '../assets/zoom-out.svg?raw'
import infoSvg from '../assets/info.svg?raw'
import Mark from 'mark.js';

const data = reactive({
  content: null,
  appInfo: {}
})
const searchInstance = ref(null)
const searchOpen = ref(false)
const loading = ref(false)
const footerExpanded = ref(false)
let collapseTimer = null

const COLLAPSE_DELAY = 10000


onMounted(() => {
  document.querySelector('main')?.addEventListener('click', handleLinkClick)
  window.addEventListener('keydown', handleKeydown, true)
  LoadInitialFile()
  LoadAppInfo()
})

onUnmounted(() => {
  document.querySelector('main')?.removeEventListener('click', handleLinkClick)
  window.removeEventListener('keydown', handleKeydown, true)
  clearCollapseTimer()
})

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

async function LoadAppInfo() {
  try {
    const info = await GetAppInfo()
    data.appInfo = info.info
    //console.dir(info)
  } catch (error) {
    console.error('Error loading app info:', error)
  }
}


async function LoadInitialFile() {
  try {
    const file = await GetStartupFile()
    if (file) {
      data.content = file.html
      await nextTick()
      const contentEl = document.getElementById('content')
      if (contentEl) {
        searchInstance.value = new Mark(contentEl)
      }
    }
  } catch (error) {
    //console.error('Error loading initial file:', error)
  }
}

async function loadFile() {
  clearCollapseTimer()
  loading.value = true
  try {
    data.content = null
    const file = await OpenFileFromFrontend()
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
  if (e.key === 'm') { e.preventDefault(); expandFooter() }
  if (ctrl && e.key.toLowerCase() === 'o') { e.preventDefault(); loadFile() }
  if (e.key === 'F3') { e.preventDefault(); openSearch() }
  if (e.key === 'F4') { e.preventDefault(); ZoomIn() }
  if (e.key === 'F5') { e.preventDefault(); ZoomOut() }
  if (e.key === 'F1') { e.preventDefault(); openAbout() }
  if (e.key === 'Escape') { e.preventDefault(); closeSearch(); closeAbout() }
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

function openAbout() {
  const aboutModal = document.getElementById('about-modal')
  if (aboutModal) {
    aboutModal.style.display = 'flex'
  }
}

function closeAbout() {
  const aboutModal = document.getElementById('about-modal')
  if (aboutModal) {
    aboutModal.style.display = 'none'
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
          <div class="item" v-html="infoSvg" title="About - F1" @click="openAbout"></div>
        </div>
      </div>
      <div class="hamburger" :class="{ visible: !footerExpanded }" @click="expandFooter">☰</div>
      <div id="search-panel">
        <span>Look for:</span>
        <input type="text" @input="find($event.target.value)" />
        <span class="close-btn" @click="closeSearch">×</span>
      </div>
      <div id="about-modal">
        <div>{{ data.appInfo.productName }}</div>
        <br>
        <small>{{ data.appInfo.description }}</small>
        <small>v{{ data.appInfo.productVersion }}</small>
        <small>Copyright © 2026 {{ data.appInfo.companyName }}</small>
        <small>License: {{ data.appInfo.license }}</small>
        <small><a href="https://github.com/huqedato/margo" target="_blank">{{ data.appInfo.productName }} on
            GitHub</a></small>
        <span class="close-btn" @click="closeAbout">×</span>
      </div>
    </div>
    <div id="no-content" v-else>
      <p v-if="loading">Loading markdown...</p>
      <button id="load-md" @click="loadFile" v-else>Load Markdown</button>
    </div>
  </main>
</template>
