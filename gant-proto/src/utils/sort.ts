export function changeSort<T>(list: T[], index: number, direction: number) {
    // 一番上の時は何もしない
    if (index === 0 && direction === -1) {
        return
    }
    // 一番下の時も何もしない
    if (index === list.length && direction === 1) {
        return
    }

    const target = list.splice(index, 1)
    list.splice(index + direction, 0, target[0])
}