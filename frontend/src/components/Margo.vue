<script setup>
import { reactive, ref, onMounted, onUnmounted } from 'vue'
import { ConvertToMd } from '../../wailsjs/go/main/App'
import { BrowserOpenURL } from '../../wailsjs/runtime/runtime'
import openFileSvg from '../assets/open-file.svg?raw'
import searchSvg from '../assets/search.svg?raw'
import zoomInSvg from '../assets/zoom-in.svg?raw'
import zoomOutSvg from '../assets/zoom-out.svg?raw'
import infoSvg from '../assets/info.svg?raw'

const data = reactive({
  content: null
})

const loading = ref(false)
const footerExpanded = ref(false)
let collapseTimer = null

const COLLAPSE_DELAY = 8000

function handleLinkClick(e) {
  const anchor = e.target.closest('a')
  if (anchor && anchor.href) {
    e.preventDefault()
    BrowserOpenURL(anchor.href)
  }
}

/** Clears any pending collapse timer. */
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

onMounted(() => {
  document.querySelector('main')?.addEventListener('click', handleLinkClick)
})

onUnmounted(() => {
  document.querySelector('main')?.removeEventListener('click', handleLinkClick)
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
      footerExpanded.value = true
      scheduleCollapse()
    }
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
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
          <div class="item" @click="loadFile" v-html="openFileSvg" title="Open File"></div>
        </div>
        <div class="card">
          <div class="item" v-html="searchSvg" title="Search"></div>
        </div>
        <div class="card">
          <div class="item" v-html="zoomInSvg" title="Zoom In"></div>
        </div>
        <div class="card">
          <div class="item" v-html="zoomOutSvg" title="Zoom Out"></div>
        </div>
        <div class="card">
          <div class="item" v-html="infoSvg" title="About"></div>
        </div>
      </div>
      <div class="hamburger" :class="{ visible: !footerExpanded }" @click="expandFooter">☰</div>
    </div>
    <div id="no-content" v-else>
      <p v-if="loading">Loading markdown...</p>
      <button id="load-md" @click="loadFile" v-else>Load Markdown</button>
    </div>
  </main>
</template>
