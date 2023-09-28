import {computed, ComputedRef, inject, ref, watch} from "vue";
import dayjs, {Dayjs} from "dayjs";
import {Holiday, Ticket, TicketUser, User} from "@/api";
import {round} from "@/utils/math";
import {GLOBAL_STATE_KEY} from "@/composable/globalState";
import {DisplayType, GanttRow} from "@/composable/ganttFacility";
import {dayBetween, ganttDateToYMDDate, getNumberOfBusinessDays} from "@/coreFunctions/manHourCalculation";

type PileUpByPerson = {
    user: User,
    labels: number[],
    // styles: StyleValue[]
    styles: any[],
    hasError: boolean
}
type PileUpByDepartment = {
    departmentId: number,
    users: number[][],
    styles: any[],
    hasError: boolean
}

type PileUpFilter = {
    departmentId: number,
    displayUsers: boolean,
}

type DisplayPileUp = {
    labels: string[],
    // styles?: StyleValue[], TODO: なぜか StyleValue[] だと再起的な何とかでlintエラーになるので一旦any
    styles?: any[]
    hasError: boolean
}

const PILEUP_DANGER_COLOR = "rgb(255 89 89)"

export const usePielUps = (
    startDate: string, endDate: string,
    tickets: ComputedRef<Ticket[]>, ticketUsers: ComputedRef<TicketUser[]>,
    displayType: ComputedRef<DisplayType>, holidays: ComputedRef<Holiday[]>) => {
    // ガントチャート描画用の日付だとフォーマットが違うので変換しておく
    startDate = ganttDateToYMDDate(startDate)
    endDate = ganttDateToYMDDate(endDate)

    const {userList, departmentList} = inject(GLOBAL_STATE_KEY)!

    const pileUpsByPerson = ref<PileUpByPerson[]>([])
    const pileUpsByDepartment = ref<PileUpByDepartment[]>([])
    const pileUpFilters = ref<PileUpFilter[]>(departmentList.map(v => {
        return <PileUpFilter>{
            departmentId: v.id,
            displayUsers: false
        }
    }))
    let refreshPileUpByPersonExclusive = false
    watch([displayType, tickets, ticketUsers, holidays], () => {
        refreshPileUps() // FIXME: watchでやるべきなのかどうかめちゃ悩む。これがMなのか？
    }, {
        deep: true
    })

    // 人単位・部署単位ともに更新する
    const refreshPileUps = () => {
        if (refreshPileUpByPersonExclusive) {
            return
        } else {
            refreshPileUpByPersonExclusive = true
        }
        pileUpsByPerson.value.length = 0
        pileUpsByDepartment.value.length = 0
        // 開始日・終了日は表示中の設備にする
        const duration = dayjs(endDate).diff(dayjs(startDate), 'day')
        // 全ユーザー分初期化する
        userList.forEach(v => {
            pileUpsByPerson.value.push({
                labels: Array(duration).fill(0),
                user: v,
                styles: Array(duration).fill({}),
                hasError: false,
            })
        })
        // 部署ごとの初期化
        departmentList.forEach(v => {
            pileUpsByDepartment.value.push({
                departmentId: v.id!,
                users: Array(duration).fill([]),
                styles: Array(duration).fill({}),
                hasError: false,
            })
        })
        // fillは同じオブジェクトを参照するため上書きする。
        pileUpsByDepartment.value.forEach(v => {
            v.users.forEach((vv, ii) => {
                v.users[ii] = []
            })
        })

        // 全てのチケットから積み上げを更新する
        tickets.value.forEach(ticket => {
            setWorkHour(
                pileUpsByPerson.value, pileUpsByDepartment.value,
                ticket,
                ticketUsers.value.filter(v => v.ticket_id === ticket.id),
                startDate, holidays.value
            )
        })
        // 稼働上限のスタイルを適応する
        pileUpsByPerson.value.forEach(v => {
            v.styles = v.labels.map(workHour => {
                if (v.user.limit_of_operation < workHour) {
                    v.hasError = v.hasError || true
                    return {color: PILEUP_DANGER_COLOR}
                } else {
                    return {}
                }
            })
        })
        if (displayType.value === "week") {
            aggregatePileUpsByWeek(startDate, endDate, pileUpsByPerson.value, pileUpsByDepartment.value)
        }
        syncDepartmentStyleByPerson();
        refreshPileUpByPersonExclusive = false
    }

    /**
     * 人の積み上げを行う。
     * 工数 / 営業日 / 人数 で均等に分配する
     * 期間が十分でない場合は最後の人に全て割り当てる。
     * 期間が十分で余りが出る場合は稼働予定が少ない順、チケットに割り当て順で積み上げる。
     * @param pileUpsByPerson
     * @param pileUpByDepartment
     * @param ticket
     * @param ticketUsers
     * @param facilityStartDate
     * @param holidays
     */
    const setWorkHour = (pileUpsByPerson: PileUpByPerson[], pileUpByDepartment: PileUpByDepartment[],
                         ticket: Ticket,
                         ticketUsers: TicketUser[],
                         facilityStartDate: string,
                         holidays: Holiday[]) => {
        // validation
        if (ticket.start_date == null || ticket.end_date == null || ticket.estimate == null ||
            (ticketUsers == null || ticketUsers.length <= 0)) {
            return
        }
        const dayjsFacilityStartDate = dayjs(facilityStartDate)
        const dayjsStartDate = dayjs(ticket.start_date)
        const dayjsEndDate = dayjs(ticket.end_date)
        const startIndex = getIndexByDate(dayjsFacilityStartDate, dayjsStartDate)
        const endIndex = getIndexByDate(dayjsFacilityStartDate, dayjsEndDate)
        // 営業日の取得
        const numberOfBusinessDays = getNumberOfBusinessDays(dayjsStartDate, dayjsEndDate, holidays)
        const workHour = ticket.estimate / numberOfBusinessDays / ticketUsers.length
        let estimate = ticket.estimate
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
        const ticketUserIds = ticketUsers.map(v => v.user_id)
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
        const departmentTarget = pileUpByDepartment.find(v => v.departmentId === ticket.department_id)
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
    const aggregatePileUpsByWeek = (term_from: string, term_to: string, pileUpsByPerson: PileUpByPerson[], pileUpsByDepartment: PileUpByDepartment[]) => {
        // 日毎で計算された結果を週ごとに集約する。
        // 集約元のIndexと集約先のIndexをマッピングする
        const dayjsEndDate = dayjs(term_to)
        let currentDate = dayjs(term_from)
        let currentStartOfWeek = dayjs(term_from).startOf("week")
        let currentWeekIndex = 0
        let currentDayIndex = 0
        const indexMap: { [index: number]: number[] } = {}
        while (currentDate.isBefore(dayjsEndDate)) {
            // 週が異なっていれば次の週とする
            if (!currentStartOfWeek.isSame(currentDate.startOf("week"))) {
                currentWeekIndex++
                currentStartOfWeek = currentDate.startOf("week")
            }
            // 週のマップがなければ初期化
            if (indexMap[currentWeekIndex] == null) {
                indexMap[currentWeekIndex] = []
            }
            // その週に紐づく日付のindexを追加する
            indexMap[currentWeekIndex].push(currentDayIndex)
            // シーケンスを進める
            currentDate = currentDate.add(1, "day")
            currentDayIndex++
        }
        pileUpsByPerson.forEach(v => {
            const labelsByDay = v.labels.concat()
            v.labels.length = 0
            v.styles.length = 0
            let hasError = false
            for (const key in indexMap) {
                v.labels.push(indexMap[key].reduce((p, c) => {
                    return p + labelsByDay[c]
                }, 0))
                hasError = indexMap[key].some(vv => labelsByDay[vv] > v.user.limit_of_operation)
                v.styles.push(hasError ? {color: PILEUP_DANGER_COLOR} : {})
                v.hasError = v.hasError || hasError // 一度trueになるとtrueになり続けるやつ
            }
        })
        pileUpsByDepartment.forEach(v => {
            // ユーザーを集約する
            const usersByDay = v.users.concat()
            v.users.length = 0
            for (const key in indexMap) {
                const r: number[] = []
                indexMap[key].forEach(v => {
                    r.push(...usersByDay[v])
                })
                // unique
                v.users.push(Array.from(new Set(r)))
            }
        })

    }

    // 山積みの並び順通りに配列を返す
    const displayPileUps = computed(() => {
        const result: DisplayPileUp[] = []
        pileUpFilters.value.forEach(f => {
            // 部署の追加、ユーザー数を追加する。
            const v = pileUpsByDepartment.value.find(v => v.departmentId === f.departmentId)!
            result.push(
                {
                    labels: v.users.map(vv => vv.length === 0 ? '' : vv.length.toString()),
                    hasError: v.hasError,
                    styles: v.styles
                })
            if (f.displayUsers) {
                const v = pileUpsByPerson.value.filter(v => v.user.department_id === f.departmentId)
                v.forEach(user => {
                    result.push(
                        {
                            labels: user.labels.map(vv => vv === 0 ? '' : round(vv).toString()),
                            styles: user.styles,
                            hasError: user.hasError,
                        }
                    )
                })
            }
        })
        return result
    })

    // 人の積み上げから部署の積み上げを同期する
    function syncDepartmentStyleByPerson() {
        pileUpsByDepartment.value.forEach(v => {
            // 部署に紐づく人間がエラーを持っているか？
            v.hasError = pileUpsByPerson.value.filter(vv => vv.user.department_id === v.departmentId).some(vv => vv.hasError)
            // 部署に紐づく担当者のエラーの複合
            pileUpsByPerson.value.filter(vv => vv.user.department_id === v.departmentId).forEach(vv => {
                vv.styles.forEach((vvv, index) => {
                    // FIXME: colorがあるときにするのはいまいち。運用的には期待値通りにはなる。
                    if (Object.keys(vvv).includes("color")) v.styles[index] = {color: PILEUP_DANGER_COLOR}
                })
            })
        })
    }

    refreshPileUps()

    return {
        pileUpFilters,
        pileUpsByDepartment,
        pileUpsByPerson,
        displayPileUps,
        refreshPileUps,
    }
}