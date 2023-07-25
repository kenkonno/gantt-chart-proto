import {Api} from "@/api/axios";
import {Post@Upper@sRequest, @Upper@} from "@/api";
import {ref} from "vue";
import {toast} from "vue3-toastify";


// ユーザー一覧。特別ref系は必要ない。
export async function use@Upper@Table() {
    const list = ref<@Upper@[]>([])
    const refresh = async () => {
        const resp = await Api.get@Upper@s()
        list.value.splice(0, list.value.length)
        list.value.push(...resp.data.list)
    }
    await refresh()
    return {list, refresh}
}

// ユーザー追加・更新。
export async function use@Upper@(@Lower@Id?: number) {

    const @Lower@ = ref<@Upper@>({
@DefaultMapping@
    })
    if (@Lower@Id !== undefined) {
        const {data} = await Api.get@Upper@sId(@Lower@Id)
        if (data.@Lower@ != undefined) {
@ResponseMapping@
        }
    }

    return {@Lower@}

}

export async function post@Upper@(@Lower@: @Upper@, emit: any) {
    const req: Post@Upper@sRequest = {
        @Lower@: @Lower@
    }
    await Api.post@Upper@s(req).then(() => {
        toast("成功しました。")
    }).finally(() => {
        emit('closeEditModal')
    })
}

export async function post@Upper@ById(@Lower@: @Upper@, emit: any) {
    const req: Post@Upper@sRequest = {
        @Lower@: @Lower@
    }
    if (@Lower@.id != null) {
        await Api.post@Upper@sId(@Lower@.id, req).then(() => {
            toast("成功しました。")
        }).finally(() => {
            emit('closeEditModal')
        })
    }
}

export async function delete@Upper@ById(id: number, emit: any) {
    await Api.delete@Upper@sId(id).then(() => {
        toast("成功しました。")
    }).finally(() => {
        emit('closeEditModal')
    })
}

