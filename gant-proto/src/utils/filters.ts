import dayjs from "dayjs";

export function dateFormat(date: string) {
    return dayjs(date).format("YYYY-MM-DD HH:mm:ss")
}
export function dateFormatYMD(date: string) {
    return dayjs(date).format("YYYY-MM-DD")
}

export function unixTimeFormat(unixTime: number) {
    return dayjs.unix(unixTime).format("YYYY-MM-DD HH:mm:ss")
}