<template>
  <div class="row p-1">
    <div class="col-2">共有リンク</div>
    <div class="col-8"><input type="text" :value="sharedLink" readonly class="form-control"/></div>
    <div class="col-2"><input type="button" value="コピー" class="btn btn-primary" @click="copyToClipBoard(sharedLink)"></div>
  </div>
  <div class="row p-1">
    <div class="col-2">ゲスト用URL</div>
    <div class="col-8"><input type="text" :value="guestLink" readonly class="form-control"/></div>
    <div class="col-2"><input type="button" value="コピー" class="btn btn-primary" @click="copyToClipBoard(guestLink)"></div>
  </div>
  <div class="row">
    <small>※ゲスト用URLを知っているすべての人がこの設備のスケジュールを閲覧できます。</small>
  </div>
</template>

<script setup lang="ts">
import {computed, inject} from "vue";
import {GLOBAL_STATE_KEY} from "@/composable/globalState";
import {useFacilitySharedLink} from "@/composable/facilitySharedLink";
import Swal from "sweetalert2"

const {currentFacilityId} = inject(GLOBAL_STATE_KEY)!
const {facilitySharedLink} = await useFacilitySharedLink(currentFacilityId)

const sharedLink = computed(() => {
  return `${window.location.protocol}//${window.location.hostname}/?facilityId=${currentFacilityId}`
})
const guestLink = computed(() => {
  return `${sharedLink.value}&uuid=${facilitySharedLink.value.uuid}`
})

const copyToClipBoard = (text: string) => {
  navigator.clipboard.writeText(text).then(() => {
    Swal.fire({
      title: 'コピーしました',
      icon: 'success',
      showCancelButton: false,
      confirmButtonColor: '#3085d6',
      cancelButtonColor: '#d33',
      confirmButtonText: '閉じる',
    })
  })
}


</script>
