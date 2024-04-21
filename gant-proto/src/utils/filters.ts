import dayjs from "dayjs";
import {DisplayType} from "@/composable/ganttFacilityMenu";

export function dateFormat(date: string) {
    return dayjs(date).format("YYYY-MM-DD HH:mm:ss")
}

export function dateFormatYMD(date: string) {
    return dayjs(date).format("YYYY-MM-DD")
}

export function unixTimeFormat(unixTime: number) {
    return dayjs.unix(unixTime).format("YYYY-MM-DD HH:mm:ss")
}

export function pileUpsLabelFormat(labels: number[], displayType: DisplayType) {
    if (displayType === "day") {
        return labels.map(v => v === 0 ? "" : pileUpLabelFormat(v))
    } else {
        return labels.map(v => v === 0 ? "" : Math.round(v * 10) / 10)
    }
}

export function pileUpLabelFormat(v: number) {
    return Math.round(v * 10 / 8) / 10
}

export function progressFormat(v: number) {
    const result = v.toFixed(1)
    return result == "NaN" ? "0.0" : result
}


