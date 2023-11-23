import {ref} from "vue";

export function useModalWithId() {
    const id = ref<number>(0)
    const originalId = ref<number>(0)
    const modalIsOpen = ref<boolean>(false)

    const openEditModal = (targetId: number, orgId: number) => {
        id.value = targetId
        originalId.value = orgId
        modalIsOpen.value = true
    }
    const closeEditModal = () => {
        id.value = 0
        modalIsOpen.value = false
    }

    return {modalIsOpen, id, originalId, openEditModal, closeEditModal}
}