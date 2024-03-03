import {Api} from "@/api/axios";
import {PostTicketMemoIdRequest, PostTicketsRequest, Ticket} from "@/api";
import {ref} from "vue";
import {toast} from "vue3-toastify";
import {Emit} from "@/const/common";


// ユーザー一覧。特別ref系は必要ない。
export async function useTicketTable() {
    const list = ref<Ticket[]>([])
    const refresh = async (ganttGroupIds: number[]) => {
        const resp = await Api.getTickets(ganttGroupIds)
        list.value.splice(0, list.value.length)
        list.value.push(...resp.data.list)
    }
    return {list, refresh}
}

// ユーザー追加・更新。
export async function useTicket(ticketId?: number) {

    const ticket = ref<Ticket>({
        id: null,
        gantt_group_id: 0,
        process_id: 0,
        department_id: 0,
        limit_date: "",
        estimate: 0,
        days_after: 0,
        start_date: "",
        end_date: "",
        progress_percent: 0,
        order: 0,
        created_at: undefined,
        updated_at: undefined
    })
    if (ticketId !== undefined) {
        const {data} = await Api.getTicketsId(ticketId)
        if (data.ticket != undefined) {
            ticket.value.id = data.ticket.id
            ticket.value.gantt_group_id = data.ticket.gantt_group_id
            ticket.value.process_id = data.ticket.process_id
            ticket.value.department_id = data.ticket.department_id
            ticket.value.limit_date = data.ticket.limit_date
            ticket.value.estimate = data.ticket.estimate
            ticket.value.days_after = data.ticket.days_after
            ticket.value.start_date = data.ticket.start_date
            ticket.value.end_date = data.ticket.end_date
            ticket.value.progress_percent = data.ticket.progress_percent
            ticket.value.order = data.ticket.order
            ticket.value.created_at = data.ticket.created_at
            ticket.value.updated_at = data.ticket.updated_at
        }
    }

    return {ticket}

}

export async function postTicket(ticket: Ticket, emit: Emit) {
    const req: PostTicketsRequest = {
        ticket: ticket
    }
    await Api.postTickets(req).then(() => {
        toast("成功しました。")
    }).finally(() => {
        emit('closeEditModal')
    })
}

export async function postTicketById(ticket: Ticket, emit: Emit) {
    const req: PostTicketsRequest = {
        ticket: ticket
    }
    if (ticket.id != null) {
        await Api.postTicketsId(ticket.id, req).then(() => {
            toast("成功しました。")
        }).finally(() => {
            emit('closeEditModal')
        })
    }
}

export async function postTicketMemoById(ticketId: number, memo: string, updatedAt: number) {
    const req: PostTicketMemoIdRequest = {
        memo: memo,
        updated_at: updatedAt,
    }
    const {data} = await Api.postTicketMemoId(ticketId, req)
    toast("成功しました。")
    return data
}


export async function deleteTicketById(id: number, emit: Emit) {
    await Api.deleteTicketsId(id).then(() => {
        toast("成功しました。")
    }).finally(() => {
        emit('closeEditModal')
    })
}



