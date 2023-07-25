import {ref} from "vue";

export function useModalWithId() {
    const id = ref<number>(0)
    const modalIsOpen = ref<boolean>(false)

    const openEditModal = (targetId: number) => {
        id.value = targetId
        modalIsOpen.value = true
    }
    const closeEditModal = () => {
        id.value = 0
        modalIsOpen.value = false
    }

    return {modalIsOpen, id, openEditModal, closeEditModal}
}