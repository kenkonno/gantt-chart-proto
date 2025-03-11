import dayjs, {Dayjs} from "dayjs";
import {Holiday} from "@/api";
import {DAYJS_FORMAT} from "@/utils/day";

/**
 * 開始日・終了日に祝日が含まれていれば終了日をその分ずらす
 * @param startDate
 * @param endDate
 * @param holidays
 */
export const getNumberOfBusinessDays = (startDate: Dayjs, endDate: Dayjs, holidays: Holiday[]) => {
    // 差がない場合は1営業日になる
    const result = endDate.diff(startDate, 'days') + 1
    // 開始日と終了日の内、祝日を除いた日数を返す
    const includes = holidays.filter(holiday => {
        const dayjsHoliday = dayjs(holiday.date)
        return dayBetween(dayjsHoliday, startDate, endDate)
    })
    return result - includes.length
}
export const ganttDateToYMDDate = (date: string) => {
    const e = date.split(" ")[0].split(".")
    return `${e[2]}-${e[1]}-${e[0]}`
}
export const YMDDateToGanttStartDate = (date: string) => {
    const e = date.split("-");
    return `${e[2]}.${e[1]}.${e[0]} 00:00`;
}
export const YMDDateToGanttEndDate = (date: string) => {
    const e = date.split("-");
    return `${e[2]}.${e[1]}.${e[0]} 23:59`;
}
export const endOfDay = (dateString: string) => {
    return dayjs(dateString).add(1, 'day').add(-1, 'minute').format(DAYJS_FORMAT)
}
export const getEndDateByRequiredBusinessDay = (startDate: Dayjs, requiredNumberOfBusinessDays: number, holidays: Holiday[]) => {
    let currentDate = startDate
    const dayjsHolidays = holidays.map(v => dayjs(v.date))
    // 1営業日だとしたらstartDateを返せばよい
    while (requiredNumberOfBusinessDays > 1) {
        currentDate = currentDate.add(1, 'day')
        if (dayjsHolidays.some(v => v.isSame(currentDate))) continue
        requiredNumberOfBusinessDays--
    }
    return currentDate.endOf('day')
}
export const dayBetween = (day: Dayjs, form: Dayjs, to: Dayjs) => {
    return (form.isSame(day) || to.isSame(day)) ||
        (day.isAfter(form) && day.isBefore(to))
}
/**
 * 開始日が祝日だった場合開始日をずらす
 * @param startDate
 * @param holidays
 */
export const adjustStartDateByHolidays = (startDate: Dayjs, holidays: Holiday[]) => {
    let result = startDate
    let endCheck = holidays.find(v => dayjs(v.date).isSame(result))
    while (endCheck) {
        result = result.add(1, "day")
        endCheck = holidays.find(v => dayjs(v.date).isSame(result))
    }
    return result
}

/**
 * 営業日分日にちを計算する
 * @param startDate
 * @param numberOfBusinessDays
 * @param holidays
 */
export const addBusinessDays = (startDate: Dayjs, numberOfBusinessDays: number, holidays: Holiday[], wantEndDate = false) => {
    let result = startDate
    // 0営業日の場合は開始を返す
    if(numberOfBusinessDays === 0 ) {
        return result
    }
    // 進める方向の決定
    let direction = 1
    if(numberOfBusinessDays < 0 ) {
        direction = -1
    }

    // 1日進める
    let recursiveLimit = 365
    while(numberOfBusinessDays !== 0 && recursiveLimit !== 0) {
        result = result.add(direction, "day")
        let isHoliday = false
        if(direction > 0 && wantEndDate) {
            // +方向の時は1minuteマイナスする
            isHoliday = holidays.find(v => dayjs(v.date).isSame(result.add(-1,'minute').startOf('day'))) != undefined
        } else {
            isHoliday = holidays.find(v => dayjs(v.date).isSame(result)) != undefined
        }
        // 祝日でなければ必要営業日を１減らす。末尾の時は祝日でも許可する
        if(!isHoliday) {
            numberOfBusinessDays -= direction
        } else {
            // 無限ループ用に再起回数の最大を制御
            recursiveLimit--
        }
    }

    return result
}