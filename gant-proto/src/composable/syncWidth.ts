import {nextTick, onMounted, onUnmounted, Ref, ref, StyleValue} from "vue";

export function useSyncWidthAndScroll(
    widthElement: Ref<HTMLDivElement | undefined>,
    parentElement: Ref<HTMLDivElement | undefined>,
    childElement: Ref<HTMLDivElement | undefined>
) {
    const syncWidth = ref<StyleValue>()!

    const resizeSyncWidth = () => {
        const parentWidth = widthElement.value?.clientWidth
        syncWidth.value = {width: parentWidth + "px", overflow: 'scroll'}
    }
    const forceScroll = () => {
        parentElement.value?.dispatchEvent(new Event('scroll'))
    }
    onMounted(() => {
        resizeSyncWidth()
        nextTick(resizeSyncWidth) // たまに上手くいかないので念のため
        parentElement.value?.addEventListener("scroll", (event) => {
            // @ts-expect-error よくわからなけどいったん抑制
            childElement.value?.scrollTo(event.srcElement.scrollLeft, 0)
        })
        childElement.value?.addEventListener("scroll", (event) => {
            // @ts-expect-error よくわからなけどいったん抑制
            parentElement.value?.scrollTo(event.srcElement.scrollLeft, 0)
        })
    })
    onUnmounted(() => {
        // @ts-expect-error よくわからなけどいったん抑制
        parentElement.value?.removeEventListener("scroll")
        // @ts-expect-error よくわからなけどいったん抑制
        childElement.value?.removeEventListener("scroll")
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