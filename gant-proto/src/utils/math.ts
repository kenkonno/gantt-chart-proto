export function round(num: number, digits = 1) {
    const d = Math.pow(10, digits)
    return Math.round(num * d) / d
}