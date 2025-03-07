<template>
  <div class="d-flex">
    <input type="file" accept=".csv" @change="onFileChange" style="display: none" ref="fileInput"/>
    <div class="align-middle">
      <p class="m-0">{{description}}</p>
      <p class="m-0"><a href="javascript:void(0)" @click="downloadSample()">サンプル</a></p>
    </div>
    <div class='dropzone'
         @dragover.prevent
         @drop='onDrop($event)'
         @click='onDropZoneClick'>
      <span v-if="!file"><input type="button" value="ファイルを選択するかドロップしてください。"> </span>
      <div v-if="file">
        {{ file.name }}
      </div>
    </div>

    <button class="btn btn-primary" @click="upload">アップロード</button>
  </div>
</template>

<script setup lang="ts">
import {ref, defineEmits} from 'vue'
import {User} from "@/api";

interface UploadFile {
  description: string,
  sample: string[][],
}

const props = defineProps<UploadFile>()

const file = ref<File | null>(null)
const fileInput = ref<HTMLInputElement | null>(null)

const onFileChange = (e: Event) => {
  if (e.target) {
    const target = e.target as HTMLInputElement
    file.value = target.files ? target.files[0] : null
  }
}

const downloadSample = () => {
  const csvContent = props.sample.map(e => e.join(",")).join("\n");

  const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
  const link = document.createElement('a');
  const url = URL.createObjectURL(blob);
  link.href = url;
  link.setAttribute('download', 'example.csv');
  document.body.appendChild(link);
  link.click();
}

const onDrop = (e: DragEvent) => {
  e.preventDefault();
  if (e.dataTransfer) {
    file.value = e.dataTransfer.files[0]
  }
}
const onDropZoneClick = () => {
  fileInput.value?.click();
}

const emit = defineEmits(['file-upload'])

const upload = () => {
  if (file.value) {
    emit('file-upload', file.value)
  }
}

</script>

<style scoped>
.dropzone {
  width: 300px;
  border: 2px dashed #ccc;
  display: flex;
  justify-content: center;
  align-items: center;
  cursor: pointer;
  vertical-align: middle;
}
</style>