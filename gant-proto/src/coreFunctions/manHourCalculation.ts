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