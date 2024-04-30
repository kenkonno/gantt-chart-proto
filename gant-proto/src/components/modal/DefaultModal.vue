<template>
  <Teleport to="body">
    <div id="overlay" class="overlay-event overlay-on">
      <div class="flex">
          <div class="modal" tabindex="-1">
            <div class="modal-dialog" :class="size()">
              <div class="modal-content">
                <div class="modal-header">
                  <h5 class="modal-title">{{ title }}</h5>
                  <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"
                          @click="$emit('closeEditModal')"></button>
                </div>
                <div class="modal-body">
                  <slot/>
                </div>
              </div>
            </div>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">

interface DefaultModal {
  title: string,
  size?: string
}

const props = defineProps<DefaultModal>()
defineEmits(['closeEditModal'])

const size = () => {
  if (props.size === "half") {
    return "half"
  }
  return "big"
}
</script>

<style scoped>
.modal {
  position: fixed;
  z-index: 999;
  width: 100%;
  display: block;
}

.half {
  width: 50%;
}

#overlay {
  position: fixed;
  top: 0;
  z-index: 998;
  width: 100vw;
  height: 100vh;
  visibility: hidden;
  opacity: 0;
  background: rgba(0, 0, 0, 0.6);
  transition: all 0.5s ease-out;
}

.modal-dialog {
  max-width: 90%;
  max-height: 90%;
  overflow-y: scroll;
}

.modal-content {
}

#overlay.overlay-on {
  visibility: visible;
  opacity: 1;
}

.flex {
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
}

</style>