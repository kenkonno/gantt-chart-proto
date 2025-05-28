import {ComputedRef, nextTick, onMounted, onUnmounted, Ref, ref, StyleValue} from "vue";


export function useSyncWidthAndScroll(
    widthElement: Ref<HTMLDivElement | undefined>,
    parentElement: Ref<HTMLDivElement | undefined>,
    childElement: Ref<HTMLDivElement | undefined>,
    showPileup: ComputedRef<boolean>,
) {
    const syncWidth = ref<StyleValue>()!

    const resizeSyncWidth = () => {
        console.log("################## RESIZE")
        const parentWidth = widthElement.value?.clientWidth
        syncWidth.value = {width: parentWidth + "px", overflow: 'scroll'}
    }

    const forceScroll = () => {
        parentElement.value?.dispatchEvent(new Event('scroll'))
    }

    // 親要素のスクロールイベントハンドラー
    const handleParentScroll = (event: Event) => {
        const target = event.target as HTMLElement

        // targetが存在し、scrollLeftプロパティを持っているか確認
        if (target && 'scrollLeft' in target && typeof target.scrollLeft === 'number') {
            // 子要素が存在するか確認
            if (childElement.value && typeof childElement.value.scrollTo === 'function') {
                childElement.value.scrollTo(target.scrollLeft, 0)
            }
        }
    }

    // 子要素のスクロールイベントハンドラー
    const handleChildScroll = (event: Event) => {

        if (!showPileup.value) { return }

        const target = event.target as HTMLElement

        // targetが存在し、scrollLeftプロパティを持っているか確認
        if (target && 'scrollLeft' in target && typeof target.scrollLeft === 'number') {
            // 親要素が存在するか確認
            if (parentElement.value && typeof parentElement.value.scrollTo === 'function') {
                parentElement.value.scrollTo(target.scrollLeft, 0)
            }
        }
    }

    onMounted(() => {
        resizeSyncWidth()
        nextTick(resizeSyncWidth) // たまに上手くいかないので念のため

        // イベントリスナーの登録
        parentElement.value?.addEventListener("scroll", handleParentScroll)
        childElement.value?.addEventListener("scroll", handleChildScroll)
    })

    onUnmounted(() => {
        // イベントリスナーの削除（正しい参照を渡す）
        parentElement.value?.removeEventListener("scroll", handleParentScroll)
        childElement.value?.removeEventListener("scroll", handleChildScroll)
    })

    return {
        syncWidth,
        resizeSyncWidth,
        forceScroll
    }
}

export function useSyncScrollY(
    parentElement: Ref<HTMLDivElement | undefined>,
    childElement: Ref<HTMLDivElement | undefined>
) {

    // TODO: gGanttGrid.$el 周りがハードコーディングになっている

    const forceScroll = () => {
        // @ts-expect-error よくわからなけどいったん抑制
        parentElement.value.$refs.ganttChart.dispatchEvent(new Event('scroll'))
    }
    onMounted(() => {
        // @ts-expect-error よくわからなけどいったん抑制
        parentElement.value.$refs.ganttChart.addEventListener("scroll", (event) => {
            // @ts-expect-error よくわからなけどいったん抑制
            childElement.value.$refs.gGanttWrapperRef.scrollTo(0, event.srcElement.scrollTop)

        })
        // @ts-expect-error よくわからなけどいったん抑制
        childElement.value.$refs.gGanttWrapperRef.addEventListener("scroll", (event) => {
            // @ts-expect-error よくわからなけどいったん抑制
            parentElement.value.$refs.ganttChart.scrollTo(0, event.srcElement.scrollTop)
        })
    })
    // 要素が消えるからeventも消さなくて良いみたい。
    // onUnmounted(() => {
    //     // @ts-expect-error よくわからなけどいったん抑制
    //     parentElement.value.$refs.ganttChart.removeEventListener("scroll")
    //     // @ts-expect-error よくわからなけどいったん抑制
    //     childElement.value.$refs.gGanttGrid.$el.removeEventListener("scroll")
    // })
    return {
        forceScroll
    }
}