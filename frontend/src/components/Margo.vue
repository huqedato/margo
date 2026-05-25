<script setup>
import { reactive, ref } from 'vue'
import { Greet } from '../../wailsjs/go/main/App'
import { ResolveFilePaths } from '../../wailsjs/runtime/runtime'

const data = reactive({
  content: null
})

const loading = ref(false)

function greet() {
  Greet(data.name).then(result => {
    data.resultText = result
  })
}

function loadFile(event) {
  //do something with the file
  loading.value = true
}

</script>

<template>
  <main>
  <div id="content" v-if="data.content" v-html="data.content"></div>
  <div id="no-content" v-else>
    <p v-if="loading">Loading...</p>
    <button id="load-md" @click="loadFile" v-else>Load Markdown File</button>
    
  </div>
</main>
</template>
