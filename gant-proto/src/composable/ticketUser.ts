import {Api} from "@/api/axios";
import {PostTicketUsersRequest, TicketUser} from "@/api";
import {ref} from "vue";
import {toast} from "vue3-toastify";
import {Emit} from "@/const/common";


// ユーザー一覧。特別ref系は必要ない。
export async function useTicketUserTable() {
    const list = ref<TicketUser[]>([])
    const refresh = async (ticketIds: number[]) => {
        const resp = await Api.getTicketUsers(ticketIds)
        list.value.splice(0, list.value.length)
        list.value.push(...resp.data.list)
    }
    return {list, refresh}
}

// ユーザー追加・更新。
export async function useTicketUser(ticketUserId?: number) {

    const ticketUser = ref<TicketUser>({
        id: null,
        ticket_id: 0,
        user_id: 0,
        order: 0,
        created_at: undefined,
        updated_at: undefined
    })
    if (ticketUserId !== undefined) {
        const {data} = await Api.getTicketUsersId(ticketUserId)
        if (data.ticketUser != undefined) {
            ticketUser.value.id = data.ticketUser.id
            ticketUser.value.ticket_id = data.ticketUser.ticket_id
            ticketUser.value.user_id = data.ticketUser.user_id
            ticketUser.value.order = data.ticketUser.order
            ticketUser.value.created_at = data.ticketUser.created_at
            ticketUser.value.updated_at = data.ticketUser.updated_at
        }
    }

    return {ticketUser}

}

export async function postTicketUser(ticketId: number, userIds: number[], emit: Emit) {
    const req: PostTicketUsersRequest = {
        ticketId: ticketId,
        userIds: userIds

    }
    await Api.postTicketUsers(req).then(() => {
        toast("成功しました。")
    }).finally(() => {
        emit('closeEditModal')
    })
}


export async function deleteTicketUserById(id: number, emit: Emit) {
    await Api.deleteTicketUsersId(id).then(() => {
        toast("成功しました。")
    }).finally(() => {
        emit('closeEditModal')
    })
}



