import {Api} from "@/api/axios";
import {Post@Upper@sRequest, @Upper@} from "@/api";
import {ref} from "vue";
import {toast} from "vue3-toastify";


// ユーザー一覧。特別ref系は必要ない。
export async function use@Upper@Table() {
    const resp = await Api.get@Upper@s()
    return resp.data.list
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

export async function post@Upper@(@Lower@: @Upper@) {
    const req: Post@Upper@sRequest = {
        @Lower@: @Lower@
    }
    await Api.post@Upper@s(req).then(() => {
        toast("成功しました。")
    })
}

export async function post@Upper@ById(@Lower@: @Upper@) {
    const req: Post@Upper@sRequest = {
        @Lower@: @Lower@
    }
    if (@Lower@.id != null) {
        await Api.post@Upper@sId(@Lower@.id, req).then(() => {
            toast("成功しました。")
        })
    }
}

export async function delete@Upper@ById(id: number) {
    await Api.delete@Upper@sId(id).then(() => {
        toast("成功しました。")
    })
}