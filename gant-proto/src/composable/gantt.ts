import dayjs, {Dayjs} from "dayjs";
import {GanttBarObject} from "@infectoone/vue-ganttastic";
import {computed, inject, ref} from "vue";
import {GanttGroup, Holiday, OperationSetting, Ticket, TicketUser, User} from "@/api";
import {Api} from "@/api/axios";
import {useGanttGroupTable} from "@/composable/ganttGroup";
import {useTicketTable} from "@/composable/ticket";
import {useTicketUserTable} from "@/composable/ticketUser";
import {useFacility} from "@/composable/facility";
import {changeSort} from "@/utils/sort";
import {GLOBAL_STATE_KEY} from "@/composable/globalState";

const format = ref("DD.MM.YYYY HH:mm")

type GanttChartGroup = {
    ganttGroup: GanttGroup // TODO: 結局ganttRowにganttGroup設定しているから設計としては微妙っぽい。
    rows: GanttRow[]
}

type GanttRow = {
    bar: GanttBarObject
    ganttGroup?: GanttGroup
    ticket?: Ticket
    ticketUsers?: TicketUser[]
}

type Header = {
    name: string,
    visible: boolean
}

type PileUpByPerson = {
    user: User,
    labels: number[],
}
type PileUpByDepartment = {
    departmentId: number,
    users: number[][],
}

const BAR_NORMAL_COLOR = "rgb(147 206 255)"
const BAR_COMPLETE_COLOR = "rgb(76 255 18)"

// const BAR_DANGER_COLOR = "rgb(255 89 89);"


function getScheduledOperatingHours(operationSettings: OperationSetting[], row: GanttRow) {
    return operationSettings.filter(operationSetting => {
        return operationSetting.facility_id === row.ganttGroup?.facility_id &&
            operationSetting.unit_id === row.ganttGroup?.unit_id
    }).reduce((accumulateValue, currentValue) => {
        const workHours = currentValue.workHours.find(v => v.process_id === row.ticket?.process_id)
        return accumulateValue + workHours!.work_hour!
    }, 0);
}

export async function useGantt() {
    // injectはsetupと同期的に呼び出す必要あり
    const {currentFacilityId} = inject(GLOBAL_STATE_KEY)!
    const {userList, processList, departmentList, holidayMap, operationSettingMap, unitMap} = inject(GLOBAL_STATE_KEY)!

    const GanttHeader = ref<Header[]>([
        {name: "ユニット", visible: true},
        {name: "工程", visible: true},
        {name: "部署", visible: false},
        {name: "担当者", visible: true},
        {name: "人数", visible: false},
        {name: "期日", visible: false},
        {name: "工数(h)", visible: true},
        {name: "日後", visible: false},
        {name: "開始日", visible: false},
        {name: "終了日", visible: false},
        {name: "進捗", visible: true},
        {name: "操作", visible: false},
    ])
    const {facility} = await useFacility(currentFacilityId)
    const chartStart = ref(dayjs(facility.value.term_from).format(format.value))
    const chartEnd = ref(dayjs(facility.value.term_to).format(format.value))
    const pileUpsByPerson = ref<PileUpByPerson[]>([])
    const pileUpsByDepartment = ref<PileUpByDepartment[]>([])

    const getUnitName = computed(() => (id: number) => {
        return unitMap[currentFacilityId].find(v => v.id === id)!.name
    })
    const getDepartmentName = computed(() => (id: number) => {
        return departmentList.find(v => v.id === id)!.name
    })
    const getProcessName = (id: number) => {
        return processList.find(v => v.id === id)!.name
    }
    const getOperationList = () => {
        return operationSettingMap[currentFacilityId]
    }
    const getHolidaysForGantt = () => {
        return holidayMap[currentFacilityId].map(v => new Date(v.date))
    }
    const getHolidays = () => {
        return holidayMap[currentFacilityId]
    }
    const getOperationSettings = () => {
        return operationSettingMap[currentFacilityId]
    }

    // とりあえず何も考えずにAPIからガントチャート表示に必要なオブジェクトを作る
    const {list: ganttGroupList, refresh: ganttGroupRefresh} = await useGanttGroupTable()
    const {list: ticketList, refresh: ticketRefresh} = await useTicketTable()
    const {list: ticketUserList, refresh: ticketUserRefresh} = await useTicketUserTable()
    await ganttGroupRefresh(currentFacilityId)
    const ganttGroupIds = ganttGroupList.value.map(v => v.id!)
    await ticketRefresh(ganttGroupIds)
    const ticketIds = ticketList.value.map(v => v.id!)
    await ticketUserRefresh(ticketIds)
    // ここからガントチャートに渡すオブジェクトを作成する
    const ganttChartGroup = ref<GanttChartGroup[]>([])
    const refreshLocalGantt = () => {
        ganttChartGroup.value.length = 0
        ganttGroupList.value.forEach(ganttGroup => {
            ganttChartGroup.value.push(
                {
                    ganttGroup: ganttGroup,
                    rows: ticketList.value.filter(ticket => ticket.gantt_group_id === ganttGroup.id)
                        .sort(v => v.order)
                        .map(ticket => ticketToGanttRow(ticket, ticketUserList.value, ganttGroup))
                }
            )
        })
    }
    // computedだとドラッグ関連が上手くいかない為やはりオブジェクトとして双方向に同期をとるようにする。
    const bars = ref<GanttBarObject[]>([])
    // ガントチャート描画用のオブジェクトの生成
    const refreshBars = () => {
        bars.value.length = 0;
        ganttChartGroup.value.forEach(unit => {
            const emptyRow: GanttBarObject = {
                beginDate: "",
                endDate: "",
                ganttBarConfig: {
                    id: Math.random().toString(),
                    style: {backgroundColor: BAR_NORMAL_COLOR}
                }
            }
            bars.value.push(...unit.rows.map(task => <GanttBarObject>{
                beginDate: dayjs(task.ticket!.start_date!).format(format.value),
                endDate: endOfDay(task.ticket!.end_date!),
                ganttGroupId: unit.ganttGroup.id,
                ganttBarConfig: {
                    hasHandles: true,
                    id: task.ticket?.id!.toString(),
                    label: getProcessName(task.ticket?.process_id == null ? -1 : task.ticket?.process_id),
                    style: {backgroundColor: BAR_NORMAL_COLOR},
                    progress: task.ticket?.progress_percent,
                    progressColor: BAR_COMPLETE_COLOR
                }
            }))
            // ボタン用の空行を追加する
            bars.value.push(emptyRow)
        })
    }
    refreshLocalGantt()
    refreshBars()

    const updateOrder = async (ganttRows: GanttRow[], index: number, direction: number) => {
        changeSort(ganttRows, index, direction)
        // ガントチャートのオブジェクトと同期をとる
        changeSort(bars.value, index, direction)
        for (const v of ganttRows) {
            v.ticket!.order = ganttRows.indexOf(v)
            await updateTicket(v.ticket!)
        }
    }

    // DBへのストア及びローカルのガントに情報を反映する
    const updateTicket = async (ticket: Ticket) => {
        const reqTicket = Object.assign({}, ticket)
        if (reqTicket.start_date != undefined) {
            reqTicket.start_date = ticket.start_date + "T00:00:00.00000+09:00"
        }
        if (reqTicket.end_date != undefined) {
            reqTicket.end_date = ticket.end_date + "T00:00:00.00000+09:00"
        }
        if (reqTicket.limit_date != undefined) {
            reqTicket.limit_date = ticket.limit_date + "T00:00:00.00000+09:00"
        }
        await Api.postTicketsId(ticket.id!, {ticket: reqTicket})
        reflectTicketToGantt(reqTicket)
        await refreshPileUps() // TODO: 山積みの更新の呼び出し
    }

    // DBへのストア及びローカルのガントに情報を反映する
    const addNewTicket = async (ganttGroupId: number) => {
        const order = ticketList.value.filter(v => v.gantt_group_id === ganttGroupId).length + 1
        const newTicket: Ticket = {
            created_at: undefined,
            days_after: undefined,
            department_id: undefined,
            end_date: undefined,
            estimate: undefined,
            gantt_group_id: ganttGroupId,
            id: undefined,
            limit_date: undefined,
            order: order,
            process_id: undefined,
            progress_percent: undefined,
            start_date: undefined,
            updated_at: 0
        }
        const {data} = await Api.postTickets({ticket: newTicket})
        ticketList.value.push(data.ticket!)
        refreshLocalGantt()
        refreshBars()
    }
    // DBへのストア及びローカルのガントに情報を反映する
    const ticketUserUpdate = async (ticket: Ticket, userIds: number[]) => {
        const {data} = await Api.postTicketUsers({ticketId: ticket.id!, userIds: userIds})
        // ticketUserList から TicketId をつけなおす
        const newTicketUserList = ticketUserList.value.filter(v => v.ticket_id !== ticket.id!)
        newTicketUserList.push(...data.ticketUsers)
        ticketUserList.value.length = 0
        ticketUserList.value.push(...newTicketUserList)

        // 人数を更新する
        ticket.number_of_worker = data.ticketUsers.length
        await updateTicket(ticket)
        // TODO: refreshLocalGanttは重くなるので性能対策が必要かも
        refreshLocalGantt()
        await refreshPileUps() // TODO: 山積みの更新の呼び出し
    }

    const adjustBar = (bar: GanttBarObject) => {
        // id をもとに VueJs 側のデータも更新する
        const targetGanttGroup = ganttChartGroup.value.find(v => v.ganttGroup.id === bar.ganttGroupId)
        if (!targetGanttGroup) {
            console.error(`Unit ID: ${bar.ganttGroupId} が現在のガントに存在しません。`, ganttChartGroup)
        } else {
            const targetTicket = targetGanttGroup.rows.find(v => v.ticket?.id!.toString() === bar.ganttBarConfig.id)
            if (!targetTicket) {
                console.error(`ID: ${bar.ganttBarConfig.id} が現在のガントに存在しません。`, ganttChartGroup)
            } else {
                targetTicket.ticket!.start_date = ganttDateToYMDDate(bar.beginDate)
                targetTicket.ticket!.end_date = ganttDateToYMDDate(bar.endDate)
                updateTicket(targetTicket.ticket!)
            }
        }
    }

    // モデル情報をガントチャート（bars）に反映させる
    const reflectTicketToGantt = (ticket: Ticket) => {
        const targetTicket = bars.value.find(v => v.ganttBarConfig.id! === ticket.id!.toString())
        if (!targetTicket) {
            console.error(`TicketID: ${ticket.id} が現在のガントに存在しません。`, bars)
        } else {
            // パフォーマンスのためにガントチャートに反映すべきものは特別にここで記述する
            targetTicket.beginDate = dayjs(ticket.start_date!).format(format.value)
            targetTicket.endDate = endOfDay(ticket.end_date!)
            if (ticket.progress_percent) {
                targetTicket.ganttBarConfig.progress = ticket.progress_percent
            }
            targetTicket.ganttBarConfig.label = getProcessName(ticket.process_id == null ? -1 : ticket.process_id)
        }
    }

    /**********************************************************
     *                   リスケジュール関連
     **********************************************************/
    /**
     * setScheduleByPersonDay
     * リスケ（工数h）重視
     * 稼働予定時間ベースでスケジュールを引き直す。
     *
     * @param rows
     * 色々考えたけど難しい（答えがない）のでマジでシンプルに１行ずつ処理するようにする。
     *
     * 基本的な仕様
     * ・人日重視で期間を変更する。
     * ・開始日・工数が設定されていない行は無視する
     * ・稼働の消化は担当者の設定順とする。
     * とりあえず作らないと話にならないので一旦作る
     */
    const setScheduleByPersonDay = (rows: GanttRow[]) => {
        const holidays = getHolidays()
        const operationSettings = getOperationList()
        let prevEndDate: Dayjs
        rows.forEach(row => {
            let startDate: Dayjs
            // 開始日・工数・担当者が未設定の場合はスキップする
            if (!row.ticket?.start_date || !row.ticket?.estimate || !row.ticketUsers?.length) {
                if (row.ticket?.end_date != null) {
                    prevEndDate = adjustStartDateByHolidays(dayjs(row.ticket.end_date).add(1, 'day'), holidays)
                }
                return;
            }
            if (row.ticket.days_after != null && row.ticket.days_after! >= 0 && prevEndDate != null) {
                startDate = prevEndDate.add(row.ticket.days_after!, 'day') // 現在日を最後に１日足しているため。
            } else {
                startDate = dayjs(row.ticket!.start_date)
            }
            // 開始日が祝日だった場合ずらす
            startDate = adjustStartDateByHolidays(startDate, holidays)

            // 工程から総稼働予定時間を取得する
            const scheduledOperatingHours = getScheduledOperatingHours(operationSettings, row) * row.ticketUsers.length
            // 必要日数。少数が出たら最小の整数にする。例：二人で5人日の場合 3日必要 5/2 = 2.5 つまり少数が出たら整数に足す
            const numberOfDaysRequired = Math.ceil(row.ticket.estimate! / scheduledOperatingHours)
            // 開始日を設定する
            row.ticket.start_date = ganttDateToYMDDate(startDate.format(format.value))
            // 終了日を決定する、祝日が含まれている場合終了日をずらす。
            row.ticket.end_date = ganttDateToYMDDate(
                getEndDateByRequiredBusinessDay(startDate, numberOfDaysRequired, holidays)
                    .add(-1, 'minute').format(format.value)
            )
            // ガントに反映する
            reflectTicketToGantt(row.ticket)
            updateTicket(row.ticket)
            // 前回の終了日を設定する
            prevEndDate = adjustStartDateByHolidays(dayjs(row.ticket.end_date).add(1, 'day'), holidays)
        })
    }

    /**
     * リスケ（日付）重視
     * 開始日・終了日重視で人数を設定する
     * 期間で満了できるように人数を割り当てる
     * ※担当者が入っていたとしても人数は必ず入っているので、人数列で参照する。
     *
     * @param rows
     */
    const setScheduleByFromTo = (rows: GanttRow[]) => {
        // 必要人数の算出。スケジュールを満了できる最少の人数を計算する。
        // 担当者が未定の時に人数を計算する
        // 必要日数 = 工数 / 労働時間
        // 必要人数 = 必要日数 / 期間日数 の数値は上に寄せる
        const holidays = getHolidays()
        const operationSettings = getOperationList()
        let prevEndDate: Dayjs
        rows.forEach(row => {
            if (row.ticket && row.ticketUsers?.length == 0 &&
                row.ticket?.estimate && row.ticket?.estimate > 0 &&
                row.ticket?.start_date && row.ticket?.end_date) {
                // 日後の処理
                if (row.ticket.days_after != null && row.ticket.days_after >= 0 && prevEndDate != null) {
                    // 現在設定されている営業日を計算する
                    let dayjsStartDate = dayjs(row.ticket.start_date)
                    let dayjsEndDate = dayjs(row.ticket.end_date)
                    const requiredNumberOfBusinessDays = getNumberOfBusinessDays(dayjsStartDate, dayjsEndDate, holidays)
                    // 開始日を日後分ずらした日付に設定する
                    dayjsStartDate = prevEndDate.add(row.ticket.days_after, 'days')
                    // 終了日を設定されている営業日分確保する
                    dayjsEndDate = getEndDateByRequiredBusinessDay(dayjsStartDate, requiredNumberOfBusinessDays, holidays)
                    // チケットに反映させる
                    row.ticket.start_date = ganttDateToYMDDate(dayjsStartDate.format(format.value))
                    row.ticket.end_date = ganttDateToYMDDate(dayjsEndDate.format(format.value))
                }

                // 工数(h)
                const estimate = row.ticket.estimate
                // 労働予定時間 ユニット & 工程 に紐づく稼働予定時間を取得
                const scheduledOperatingHours = getScheduledOperatingHours(operationSettings, row)
                // 必要日数
                const requiredNumberOfDays = estimate / scheduledOperatingHours
                // 日数（期間）
                const numberOfDays = getNumberOfDays(row, holidays)
                // 必要人数
                row.ticket.number_of_worker = Math.ceil(requiredNumberOfDays / numberOfDays)
                // 前回の終了日を設定する
                updateTicket(row.ticket)
            }
            prevEndDate = adjustStartDateByHolidays(dayjs(row.ticket?.end_date).add(1, 'day'), holidays)
        })
    }
    /**********************************************************
     *                         山積み関連
     **********************************************************/
    let refreshPileUpByPersonExclusive = false
    // 人単位・部署単位ともに更新する
    const refreshPileUps = async () => {
        if (refreshPileUpByPersonExclusive) {
            return
        } else {
            refreshPileUpByPersonExclusive = true
        }
        pileUpsByPerson.value.length = 0
        pileUpsByDepartment.value.length = 0
        // 開始日・終了日は表示中の設備にする
        const duration = dayjs(facility.value.term_to).diff(dayjs(facility.value.term_from), 'day')
        // 全ユーザー分初期化する
        userList.forEach(v => {
            pileUpsByPerson.value.push({labels: Array(duration).fill(0), user: v})
        })
        // 部署ごとの初期化
        departmentList.forEach(v => {
            pileUpsByDepartment.value.push({departmentId: v.id!, users: Array(duration).fill([])})
        })
        // fillは同じオブジェクトを参照するため上書きする。
        pileUpsByDepartment.value.forEach(v => {
            v.users.forEach((vv, ii) => {
                v.users[ii] = []
            })
        })

        // 全てのチケットから積み上げを更新する
        ganttChartGroup.value.forEach(unit => {
            unit.rows.forEach(row => {
                setWorkHour(pileUpsByPerson.value, pileUpsByDepartment.value, row, facility.value.term_from, getHolidays(), getOperationSettings())
            })
        })
        refreshPileUpByPersonExclusive = false
    }

    /**
     * 人の積み上げを行う。
     * 有効な期間から予定稼働時間をハメていく。
     * 期間が十分でない場合は最後の人に全て割り当てる。
     * 期間が十分で余りが出る場合は稼働予定が少ない順、チケットに割り当て順で積み上げる。
     * @param pileUpsByPerson
     * @param pileUpByDepartment
     * @param row
     * @param facilityStartDate
     * @param holidays
     * @param operationSettings
     */
    const setWorkHour = (pileUpsByPerson: PileUpByPerson[], pileUpByDepartment: PileUpByDepartment[],
                         row: GanttRow,
                         facilityStartDate: string,
                         holidays: Holiday[],
                         operationSettings: OperationSetting[]) => {
        // validation
        if (row.ticket?.start_date == null || row.ticket?.end_date == null || row.ticket.estimate == null ||
            (row.ticketUsers == null || row.ticketUsers?.length <= 0)) {
            return
        }
        const dayjsFacilityStartDate = dayjs(facilityStartDate)
        const dayjsStartDate = dayjs(row.ticket?.start_date)
        const dayjsEndDate = dayjs(row.ticket?.end_date)
        const startIndex = getIndexByDate(dayjsFacilityStartDate, dayjsStartDate)
        const endIndex = getIndexByDate(dayjsFacilityStartDate, dayjsEndDate)
        const workHour = getScheduledOperatingHours(operationSettings, row)
        let estimate = row.ticket?.estimate
        const holidayIndexes = holidays.filter(v => {
            return dayBetween(dayjs(v.date), dayjsStartDate, dayjsEndDate)
        }).map(v => getIndexByDate(dayjsFacilityStartDate, dayjs(v.date)))
        // 有効なIndexを設定する
        const validIndexes: number[] = []
        for (let i = 0; i + startIndex <= endIndex; i++) {
            validIndexes.push(i + startIndex)
        }
        // 祝日を削除する
        holidayIndexes.forEach(v => {
            const i = validIndexes.indexOf(v)
            if (i > -1) {
                validIndexes.splice(i, 1)
            }
        })
        const lastIndex = validIndexes[validIndexes.length - 1]
        // 対象の取得
        const ticketUserIds = row.ticketUsers?.map(v => v.user_id)
        const targets = pileUpsByPerson.filter(v => ticketUserIds.includes(v.user.id!))
        // 並び替える (末端の予定工数が少ない順, チケットにアサインされた順)
        targets.sort((a, b) => {
            if (a.labels[lastIndex] < b.labels[lastIndex]) return -1;
            if (a.labels[lastIndex] > b.labels[lastIndex]) return 1;
            if (ticketUserIds.indexOf(a.user.id!) < ticketUserIds.indexOf(b.user.id!)) return -1;
            if (ticketUserIds.indexOf(a.user.id!) > ticketUserIds.indexOf(b.user.id!)) return 1;
            return 0
        })
        // 対象部署を取得
        const departmentTarget = pileUpByDepartment.find(v => v.departmentId === row.ticket?.department_id)
        // 稼働時間を加算する。
        validIndexes.forEach((validIndex, index) => {
            targets.forEach((v, targetIndex) => {
                if (estimate < 0) {
                    return
                }
                // 最終行かつ最終の人で予定工数が余っていた場合は全て割り当てる。
                if (index === (validIndexes.length - 1) && targetIndex === (targets.length - 1)) {
                    if (estimate - workHour > 0) {
                        v.labels[validIndex] += estimate
                        if (departmentTarget != null && !departmentTarget.users[validIndex].includes(v.user.id!)) {
                            departmentTarget.users[validIndex].push(v.user.id!)
                        }
                        return
                    }
                }
                if (estimate - workHour < 0) {
                    v.labels[validIndex] += estimate
                } else {
                    v.labels[validIndex] += workHour
                }
                if (departmentTarget != null && !departmentTarget.users[validIndex].includes(v.user.id!)) {
                    departmentTarget.users[validIndex].push(v.user.id!)
                }
                estimate -= workHour
            })
        })
    }
    // 日付からどのindexに該当するか取得する
    const getIndexByDate = (facilityStartDate: Dayjs, date: Dayjs) => {
        return date.diff(facilityStartDate, 'days')
    }
    await refreshPileUps()


    return {
        chartStart,
        chartEnd,
        ganttChartGroup,
        bars,
        pileUpsByPerson,
        pileUpsByDepartment,
        format,
        GanttHeader,
        setScheduleByPersonDay,
        setScheduleByFromTo,
        adjustBar,
        addNewTicket,
        updateTicket,
        ticketUserUpdate,
        updateOrder,
        refreshPileUpByPerson: refreshPileUps,
        getUnitName,
        getDepartmentName,
        getOperationList,
        getHolidaysForGantt
    }
}


/**
 * 祝日を除く期間の日数を返却する
 * @param row
 * @param holidays
 */
function getNumberOfDays(row: GanttRow, holidays: Holiday[]) {
    const dayjsStartDate = dayjs(row.ticket?.start_date)
    const dayjsEndDate = dayjs(row.ticket?.end_date)
    const numberOfHolidays = holidays.filter(v => {
        return dayBetween(dayjs(v.date), dayjsStartDate, dayjsEndDate)
    })
    const numberOfDays = dayjsEndDate.diff(dayjsStartDate, 'day') + 1
    return numberOfDays - numberOfHolidays.length
}


const ticketToGanttRow = (ticket: Ticket, ticketUserList: TicketUser[], ganttGroup: GanttGroup) => {
    if (ticket.start_date != undefined) {
        ticket.start_date = ticket.start_date.substring(0, 10)
    }
    if (ticket.end_date != undefined) {
        ticket.end_date = ticket.end_date.substring(0, 10)
    }
    if (ticket.limit_date != undefined) {
        ticket.limit_date = ticket.limit_date.substring(0, 10)
    }

    const result: GanttRow = {
        bar: {
            beginDate: dayjs(ticket.start_date!).format(format.value),
            endDate: dayjs(ticket.end_date!).format(format.value),
            ganttBarConfig: {
                hasHandles: true,
                id: ticket.id!.toString(), // TODO: まあIDが数字でもいっか
                label: "", // TODO: コメントをラベルにする
            },
        },
        ticket: ticket,
        ticketUsers: ticketUserList.filter(v => v.ticket_id === ticket.id),
        ganttGroup: ganttGroup
    }
    return result
}

const ganttDateToYMDDate = (date: string) => {
    const e = date.split(" ")[0].split(".")
    return `${e[2]}-${e[1]}-${e[0]}`
}


const endOfDay = (dateString: string) => {
    return dayjs(dateString).add(1, 'day').add(-1, 'minute').format(format.value)
}
/**
 * 開始日・終了日に祝日が含まれていれば終了日をその分ずらす
 * @param startDate
 * @param endDate
 * @param holidays
 */
const getNumberOfBusinessDays = (startDate: Dayjs, endDate: Dayjs, holidays: Holiday[]) => {
    // 差がない場合は1営業日になる
    const result = endDate.diff(startDate, 'days') + 1
    // 開始日と終了日の内、祝日を除いた日数を返す
    const includes = holidays.filter(holiday => {
        const dayjsHoliday = dayjs(holiday.date)
        return dayBetween(dayjsHoliday, startDate, endDate)
    })
    return result + includes.length
}

const getEndDateByRequiredBusinessDay = (startDate: Dayjs, requiredNumberOfBusinessDays: number, holidays: Holiday[]) => {
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

const dayBetween = (day: Dayjs, form: Dayjs, to: Dayjs) => {
    return (form.isSame(day) || to.isSame(day)) ||
        (day.isAfter(form) && day.isBefore(to))
}

/**
 * 開始日が祝日だった場合開始日をずらす
 * @param startDate
 * @param holidays
 */
const adjustStartDateByHolidays = (startDate: Dayjs, holidays: Holiday[]) => {
    let result = startDate
    let endCheck = holidays.find(v => dayjs(v.date).isSame(result))
    while (endCheck) {
        result = result.add(1, "day")
        endCheck = holidays.find(v => dayjs(v.date).isSame(result))
    }
    return result
}