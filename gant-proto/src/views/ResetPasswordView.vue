<template>
  <default-modal title="パスワード初期化" @close-edit-modal="closeModalProxy" :size="'half'">
    <async-user-edit :id="id" @close-edit-modal="closeModalProxy" :mode="'reset-password'"></async-user-edit>
  </default-modal>
</template>

<script setup lang="ts">
import AsyncUserEdit from "@/components/user/AsyncUserEdit.vue";
import DefaultModal from "@/components/modal/DefaultModal.vue";
import {useModalWithId} from "@/composable/modalWIthId";
import {provide} from "vue";
import {
  GLOBAL_ACTION_KEY,
  GLOBAL_GETTER_KEY,
  GLOBAL_MUTATION_KEY,
  GLOBAL_STATE_KEY,
  useGlobalState
} from "@/composable/globalState";
import {getUserInfo} from "@/composable/auth";
import {isNumber} from "@tiptap/vue-3";
import router from "@/router";

const {globalState, actions, mutations, getters} = await useGlobalState()
provide(GLOBAL_STATE_KEY, globalState.value)
provide(GLOBAL_ACTION_KEY, actions)
provide(GLOBAL_MUTATION_KEY, mutations)
provide(GLOBAL_GETTER_KEY, getters)

const {modalIsOpen, id, openEditModal, closeEditModal} = useModalWithId()

const {userInfo} = getUserInfo()
if (userInfo !== undefined && isNumber(userInfo?.id)) {
  id.value = userInfo?.id
}

const emit = defineEmits(["update"])
const closeModalProxy = async () => {
  closeEditModal()
  router.push("/all-view")
}

</script>