<template>
  <div class="tiptap-toolbar clearfix">
    <button @click="editor.chain().focus().toggleHeading().run()">H</button>
    <button @click="editor.chain().focus().toggleBold().run()">B</button>
    <button @click="editor.chain().focus().toggleItalic().run()">I</button>
    <button @click="editor.chain().focus().toggleBulletList().run()">・</button>
    <button @click="editor.chain().focus().toggleOrderedList().run()">1.</button>
    <button @click="editor.chain().focus().toggleBlockquote().run()">></button>
    <button @click="editor.chain().focus().insertTable({ rows: 3, cols: 3, withHeaderRow: true }).run()"><span
        class="material-symbols-outlined">table</span></button>
    <button class="table-function" @click="editor.chain().focus().addColumnAfter().run()">
      列の追加
    </button>
    <button class="table-function" @click="editor.chain().focus().deleteColumn().run()">
      列の削除
    </button>
    <button class="table-function" @click="editor.chain().focus().addRowAfter().run()">
      行の追加
    </button>
    <button class="table-function" @click="editor.chain().focus().deleteRow().run()">
      行の削除
    </button>

  </div>
  <editor-content :editor="editor"/>
</template>


<script setup lang="ts">

import {useEditor, EditorContent} from '@tiptap/vue-3'
import StarterKit from '@tiptap/starter-kit'
import {TableCell} from "@tiptap/extension-table-cell";
import {TableRow} from "@tiptap/extension-table-row";
import {TableHeader} from "@tiptap/extension-table-header";
import {Table} from "@tiptap/extension-table";

const editor = useEditor({
  content: '<p>I’m running Tiptap with Vue.js. 🎉</p>',
  extensions: [
    StarterKit,
    Table.configure({resizable: true,}),
    TableRow,
    TableHeader,
    TableCell,
  ],
})
</script>

<style lang="scss">
.tiptap-toolbar {
  border: 1px solid #aaaaaa;
  padding: 2px;
  user-select: none;
  -moz-user-select: none; /* Firefox */
  -webkit-user-select: none; /* Chrome, Safari, and Opera */
  -ms-user-select: none; /* IE 10 and 11 */
  > button {
    float: left;
    width: 24px;
    height: 24px;
    border: 1px solid #aaaaaa;
    display: inline-block;
    line-height: 24px;
    text-align: center;
    cursor: pointer;
    margin-left: 3px;
    padding: 0;
  }
  > .table-function {
    margin: 0;
    width: auto;
    padding: 0 5px;
  }
}

.tiptap {
  height: 100%;
  border: 1px solid #aaaaaa;
  text-align: left;
  padding: 5px 10px;
  min-height: 500px;

  ul {
    li {
    }
  }

  p {
    margin: 0
  }

  blockquote {
    border-left: 2px solid #aaaaaa;
    padding-left: 1rem;
  }

  table {
    width: 100%;
    text-align: center;
    border-collapse: collapse;
    border-spacing: 0;
    border: solid 2px #aaaaaa;
  }
  th {
    border: solid 1px #aaaaaa;
    border-bottom: solid 2px #aaaaaa;
  }
  td {
    border-left: solid 1px #aaaaaa;
    text-align: left;
  }
  tr {
    border-top: dashed 1px #aaaaaa;
  }
  &.resize-cursor {
    cursor: w-resize;
  }
  .selectedCell {
    background-color: #f6f6f6;
  }

}
</style>
