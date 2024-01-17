import {computed, ComputedRef, inject, Ref, ref, UnwrapRef, watch} from "vue";
import dayjs, {Dayjs} from "dayjs";
import {Department, Holiday, Ticket, TicketUser, User} from "@/api";
import {round} from "@/utils/math";
import {GLOBAL_STATE_KEY} from "@/composable/globalState";
import {dayBetween, ganttDateToYMDDate, getNumberOfBusinessDays} from "@/coreFunctions/manHourCalculation";
import {Api} from "@/api/axios";
import {FacilityStatus} from "@/const/common";
import {DisplayType} from "@/composable/ganttFacilityMenu";

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
const PILEUP_MERGE_COLOR = "rgb(68 141 255)"

function initPileUps(
    endDate: string, startDate: string,
    userList: User[], pileUpsByPerson: Ref<UnwrapRef<PileUpByPerson[]>>, departmentList: Department[],
    pileUpsByDepartment: Ref<UnwrapRef<PileUpByDepartment[]>>,
    defaultPileUpsByPerson: PileUpByPerson[] = [], defaultPileUpsByDepartment: PileUpByDepartment[] = [],
    globalStartDate?: string,) {

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

    // defaultがある場合はindexを計算してつけなおす。大小関係がるのでマージする必要がある。
    // マージし始めるindexを計算する
    if (globalStartDate != null) {
        const mergeStartIndex = dayjs(startDate).diff(dayjs(globalStartDate), 'day')
        // 人の積み上げ
        pileUpsByPerson.value.forEach(row => {
            const target = defaultPileUpsByPerson.find(v => v.user.id === row.user.id)!
            row.labels.forEach((v, index) => {
                const targetIndex = mergeStartIndex + index
                if (0 <= targetIndex && targetIndex < target.labels.length && target.labels[mergeStartIndex + index] != 0) {
                    row.labels[index] += target.labels[mergeStartIndex + index]
                    row.styles[index] = {color: PILEUP_MERGE_COLOR}
                }
            })
        })
        // 部署の積み上げ
        pileUpsByDepartment.value.forEach(row => {
            const target = defaultPileUpsByDepartment.find(v => v.departmentId === row.departmentId)!
            row.users.forEach((v, index) => {
                const targetIndex = mergeStartIndex + index
                if (0 <= targetIndex && targetIndex < target.users.length && target.users[mergeStartIndex + index].length > 0) {
                    row.users[index].push(...target.users[mergeStartIndex + index])
                }
            })
        })
    }

}

export const usePielUps = (
    startDate: string, endDate: string,
    tickets: ComputedRef<Ticket[]>, ticketUsers: ComputedRef<TicketUser[]>,
    displayType: ComputedRef<DisplayType>,
    holidays: ComputedRef<Holiday[]>, departmentList: Department[], userList: User[],
    selectedDepartment: Ref<number | undefined>,
    selectedUser: Ref<number | undefined>,
    defaultPileUpsByPerson: PileUpByPerson[] = [], defaultPileUpsByDepartment: PileUpByDepartment[] = [],
    globalStartDate?: string,
) => {
    // ガントチャート描画用の日付だとフォーマットが違うので変換しておく
    startDate = ganttDateToYMDDate(startDate)
    endDate = ganttDateToYMDDate(endDate)

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
        initPileUps(endDate, startDate, userList, pileUpsByPerson, departmentList, pileUpsByDepartment, defaultPileUpsByPerson, defaultPileUpsByDepartment, globalStartDate);

        // 全てのチケットから積み上げを更新する
        tickets.value.forEach(ticket => {
            setWorkHour(
                pileUpsByPerson.value, pileUpsByDepartment.value,
                ticket,
                ticketUsers.value.filter(v => v.ticket_id === ticket.id),
                startDate, holidays.value, userList
            )
        })
        // 稼働上限のスタイルを適応する
        pileUpsByPerson.value.forEach(v => {
            v.styles = v.labels.map((workHour, index) => {
                if (v.user.limit_of_operation < workHour) {
                    v.hasError = v.hasError || true
                    return {color: PILEUP_DANGER_COLOR}
                } else {
                    return v.styles[index]
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
     * @param userList
     */
    const setWorkHour = (pileUpsByPerson: PileUpByPerson[], pileUpByDepartment: PileUpByDepartment[],
                         ticket: Ticket,
                         ticketUsers: TicketUser[],
                         facilityStartDate: string,
                         holidays: Holiday[], userList: User[]) => {

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
                        return
                    }
                }
                if (estimate - workHour < 0) {
                    v.labels[validIndex] += estimate
                } else {
                    v.labels[validIndex] += workHour
                }
                estimate -= workHour
            })
            // このチケットに登場した人から部署の設定を行う。ユーザーから部署を引いて足す
            ticketUserIds.map( v => userList.find(user => user.id == v )).forEach(user => {
                if (user == undefined) return
                const departmentTarget = pileUpByDepartment.find(v => v.departmentId === user.department_id)
                if (departmentTarget != null && !departmentTarget.users[validIndex].includes(user.id!)) {
                    departmentTarget.users[validIndex].push(user.id!)
                }
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
            const stylesByDay = v.styles.concat()
            v.labels.length = 0
            v.styles.length = 0
            let hasError = false
            for (const key in indexMap) {
                v.labels.push(indexMap[key].reduce((p, c) => {
                    return p + labelsByDay[c]
                }, 0))
                hasError = indexMap[key].some(vv => labelsByDay[vv] > v.user.limit_of_operation)
                // 他設備データがあれば青にする
                const defaultStyle = indexMap[key].some(vv => Object.keys(stylesByDay[vv]).includes('color') && stylesByDay[vv]['color'] == PILEUP_MERGE_COLOR) ? {color: PILEUP_MERGE_COLOR} : {};
                v.styles.push(hasError ? {color: PILEUP_DANGER_COLOR} : defaultStyle)
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
                // FIXME: uniqueなんか間違ってね？
                v.users.push(Array.from(new Set(r)))
            }
        })

    }

    // 山積みの並び順通りに配列を返す
    const displayPileUps = computed(() => {
        // TODO: こことvue側が重複コードになっている
        let targetDepartmentId = selectedDepartment.value
        if (selectedUser.value != undefined) {
            targetDepartmentId = userList.find( v => v.id == selectedUser.value)?.department_id
        } else {
            targetDepartmentId = selectedDepartment.value
        }
        const filteredPileUpsByDepartment = pileUpsByDepartment.value.filter(v => {
            if (targetDepartmentId == undefined) {
                return true
            } else {
                return v.departmentId == targetDepartmentId
            }
        })
        const filteredPileUpsByPerson = pileUpsByPerson.value.filter(v => {
            if (selectedUser.value == undefined) {
                return true
            } else {
                return v.user.id == selectedUser.value
            }
        })
        const result: DisplayPileUp[] = []
        pileUpFilters.value.forEach(f => {
            // 部署の追加、ユーザー数を追加する。
            const v = filteredPileUpsByDepartment.find(v => v.departmentId === f.departmentId)!
            if ( v == undefined) {
                return
            }
            result.push(
                {
                    labels: v.users.map(vv => vv.length === 0 ? '' : vv.length.toString()),
                    hasError: v.hasError,
                    styles: v.styles
                })
            if (f.displayUsers) {
                const v = filteredPileUpsByPerson.filter(v => v.user.department_id === f.departmentId)
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
                    if (Object.keys(vvv).includes("color") && vvv["color"] === PILEUP_DANGER_COLOR) v.styles[index] = {color: PILEUP_DANGER_COLOR}
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
/**
 * FIXME: パフォーマンス的にこれはサーバーサイドでやるべきこと。設計をし直す。
 * @param excludeFacilityId
 * @param displayType
 */
export const getDefaultPileUps = async (
    excludeFacilityId: number,
    displayType: DisplayType,
    selectedDepartment: Ref<number | undefined>,
    selectedUser: Ref<number | undefined>,
) => {
    const {facilityList, userList, departmentList, facilityTypes} = inject(GLOBAL_STATE_KEY)!
    const filteredFacilityList = facilityList.filter(v => v.status === FacilityStatus.Enabled)

    const {data: allTickets} = await Api.getAllTickets(facilityTypes)
    const {data: allTicketUsers} = await Api.getTicketUsers(allTickets.list.map(v => v.id!))
    const {data: allData} = await Api.getPileUps(excludeFacilityId, facilityTypes)

    // 全設備の最小
    const startDate: string = filteredFacilityList.slice().sort((a, b) => {
        return a.term_from > b.term_from ? 1 : -1
    }).shift()!.term_from.substring(0, 10)
    // 全設備の最大
    const endDate: string = filteredFacilityList.slice().sort((a, b) => {
        return a.term_to > b.term_to ? 1 : -1
    }).pop()!.term_to.substring(0, 10)

    const defaultPileUpsByPerson = ref<PileUpByPerson[]>([])
    const defaultPileUpsByDepartment = ref<PileUpByDepartment[]>([])
    initPileUps(endDate, startDate, userList, defaultPileUpsByPerson, departmentList, defaultPileUpsByDepartment);

    for (const facility of filteredFacilityList) {
        // 対象外の設備はスキップ
        if (facility.id === excludeFacilityId) continue
        const targetData = allData.list.find(v => v.facilityId === facility.id)
        // 対象データなしの場合もスキップ
        if (targetData == undefined) continue

        const innerHolidays: Holiday[] = []
        innerHolidays.push(...targetData.holidays)

        const holidays = computed(() => {
            return innerHolidays
        })
        const tickets = computed(() => {
            return allTickets.list.filter(v => targetData.ganttGroups.map(vv => vv.id!).includes(v.gantt_group_id))
        })
        const ticketUsers = computed(() => {
            return allTicketUsers.list.filter(v => tickets.value.map(vv => vv.id).includes(v.ticket_id))
        })

        const {pileUpsByPerson, pileUpsByDepartment} = usePielUps(startDate, endDate,
            tickets, ticketUsers,
            computed(() => displayType),
            holidays, departmentList, userList, selectedDepartment, selectedUser)

        // defaultPileUpsに計上する, stylesとかは一旦無視でいいと思う。本チャンの方で着色するので。
        pileUpsByPerson.value.forEach(pileUp => {
            const target = defaultPileUpsByPerson.value.find(v => v.user === pileUp.user)!
            target.labels.forEach((v, i) => {
                target.labels[i] += pileUp.labels[i]
            })
        })
        pileUpsByDepartment.value.forEach(pileUp => {
            const target = defaultPileUpsByDepartment.value.find(v => v.departmentId === pileUp.departmentId)!
            target.users.forEach((v, i) => {
                const r: number[] = []
                r.push(...v, ...pileUp.users[i])
                target.users[i].length = 0
                target.users[i].push(...Array.from(new Set(r)))
            })
        })
    }
    return {
        globalStartDate: startDate,
        defaultPileUpsByPerson,
        defaultPileUpsByDepartment
    }
}

