
export const initScroll = (v: { xPosition: number }, refElement: HTMLDivElement | undefined) => {
    if (v != null && v.xPosition != 0) {
        const to = v.xPosition - 30 * 3 // 3マス分くらいは余裕を持たせる
        if (to > 90) {
            refElement?.scrollTo(to, 0)
        } else {
            refElement?.scrollTo(0, 0)
        }
    }
}
