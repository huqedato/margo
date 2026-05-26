<script setup>
import { reactive, ref, onMounted, onUnmounted } from 'vue'
import { ConvertToMd } from '../../wailsjs/go/main/App'
import { BrowserOpenURL } from '../../wailsjs/runtime/runtime'

const data = reactive({
  content: null
})

const loading = ref(false)

/**
 * Intercepts clicks on anchor tags inside the rendered markdown
 * and opens them in the host OS default browser instead of the webview.
 */
function handleLinkClick(e) {
  const anchor = e.target.closest('a')
  if (anchor && anchor.href) {
    e.preventDefault()
    BrowserOpenURL(anchor.href)
  }
}

onMounted(() => {
  // Event delegation on the main container — always present in DOM
  document.querySelector('main')?.addEventListener('click', handleLinkClick)
})

onUnmounted(() => {
  document.querySelector('main')?.removeEventListener('click', handleLinkClick)
})

async function loadFile() {
  loading.value = true
  try {
    data.content = null
    const file = await ConvertToMd()
    if (file) {
      //console.dir(file)
      data.content = file.html
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
    <div id="content" v-if="data.content" v-html="data.content"></div>
    <div id="no-content" v-else>
      <p v-if="loading">Loading...</p>
      <button id="load-md" @click="loadFile" v-else>Load Markdown</button>
    </div>
  </main>
</template>
