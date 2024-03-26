import {computed, ComputedRef, inject, onBeforeUnmount, Ref, ref, UnwrapRef, watch} from "vue";
import dayjs, {Dayjs} from "dayjs";
import {Department, Facility, Holiday, Ticket, TicketUser, User} from "@/api";
import {round} from "@/utils/math";
import {GLOBAL_STATE_KEY} from "@/composable/globalState";
import {dayBetween, ganttDateToYMDDate, getNumberOfBusinessDays} from "@/coreFunctions/manHourCalculation";
import {Api} from "@/api/axios";
import {FacilityStatus, FacilityType} from "@/const/common";
import {DisplayType} from "@/composable/ganttFacilityMenu";
import {globalFilterGetter, globalFilterMutation} from "@/utils/globalFilterState";
import {pileUpLabelFormat} from "@/utils/filters";


export type PileUps = {
    departmentId: number
    labels: number[]
    styles: any[],
    display: boolean,
    assignedUser: AssignedPileUp
    unAssignedPileUp: UnAssignedPileUp
    noOrdersReceivedPileUp: NoOrdersReceivedPileUp
}
type AssignedPileUp = {
    labels: number[],
    styles: any[],
    display: boolean,
    users: PileUpByPerson[]
}
type UnAssignedPileUp = {
    labels: number[],
    styles: any[],
    display: boolean,
    facilities: PileUpByFacility[]
}
type NoOrdersReceivedPileUp = {
    labels: number[],
    styles: any[],
    display: boolean,
    facilities: PileUpByFacility[]
}

interface PileUpRow {
    labels: number[],
    styles: any[],
}

type PileUpByPerson = {
    user: User,
    labels: number[],
    styles: any[],
    hasError: boolean
}
type PileUpByFacility = {
    facilityId: number,
    labels: number[],
    styles: any[],
    hasError: boolean
}

export type PileUpFilter = {
    departmentId: number
    display: boolean
    assignedUser: boolean
    unAssignedPileUp: boolean
    noOrdersReceivedPileUp: boolean
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
    endDate: string, startDate: string, isALlMode: boolean,
    userList: User[], pileUps: Ref<UnwrapRef<PileUps[]>>, departmentList: Department[], facilityList: Facility[],
    defaultPileUps: PileUps[] = [], globalStartDate?: string,) {

    pileUps.value.length = 0
    // 開始日・終了日は表示中の案件にする
    const duration = dayjs(endDate).diff(dayjs(startDate), 'day')
    const orderedFacilities = facilityList.filter(v => v.type === FacilityType.Ordered)
    const noOrderedFacilities = facilityList.filter(v => v.type === FacilityType.Prepared)
    // 部署でループをし、紐づくユーザーと全ての案件分初期化する
    departmentList.forEach(v => {
        // 部署に紐づくユーザー
        const users = userList.filter(user => user.department_id == v.id!)
        const savedPileUpsFilter = globalFilterGetter.getPileUpsFilter()
        const filter = savedPileUpsFilter.find(vv => vv.departmentId === v.id)
        pileUps.value.push({
            departmentId: v.id!,
            labels: Array(duration).fill(0),
            styles: Array(duration).fill(0),
            display: filter ? filter.display : false,
            assignedUser: {
                labels: Array(duration).fill(0),
                styles: Array(duration).fill(0),
                display: filter ? filter.assignedUser : false,
                users: users.map(user => {
                    return {
                        user: user,
                        labels: Array(duration).fill(0),
                        styles: Array(duration).fill(0),
                        hasError: false,
                    }
                })
            },
            unAssignedPileUp: {
                labels: Array(duration).fill(0),
                styles: Array(duration).fill(0),
                display: filter ? filter.unAssignedPileUp : false,
                facilities: orderedFacilities.map(facility => {
                    return {
                        facilityId: facility.id!,
                        labels: Array(duration).fill(0),
                        styles: Array(duration).fill(0),
                        hasError: false,
                    }
                })
            },
            noOrdersReceivedPileUp: {
                labels: Array(duration).fill(0),
                styles: Array(duration).fill(0),
                display: filter ? filter.noOrdersReceivedPileUp : false,
                facilities: noOrderedFacilities.map(facility => {
                    return {
                        facilityId: facility.id!,
                        labels: Array(duration).fill(0),
                        styles: Array(duration).fill(0),
                        hasError: false,
                    }
                })
            }
        })
    })

    // defaultがある場合はindexを計算してつけなおす。大小関係があるのでマージする必要がある。
    // マージし始めるindexを計算する
    if (globalStartDate != null) {
        const mergeStartIndex = dayjs(startDate).diff(dayjs(globalStartDate), 'day');
        pileUps.value.forEach((row, i) => {
            // 部署のマージ
            const target = defaultPileUps.find(v => v.departmentId === row.departmentId);
            if (target == undefined) return console.warn("departmentId is not exists", row.departmentId)
            mergeAndUpdate(target, row, mergeStartIndex, isALlMode);

            // アサイン済み（計）のマージ
            mergeAndUpdate(target.assignedUser, row.assignedUser, mergeStartIndex, isALlMode)
            // アサイン済みユーザーのマージ
            row.assignedUser.users.forEach(v => {
                const vv = target.assignedUser.users.find(vvv => vvv.user.id === v.user.id)!
                mergeAndUpdate(vv, v, mergeStartIndex, isALlMode)
            })

            // 未アサイン済みのマージ
            mergeAndUpdate(target.unAssignedPileUp, row.unAssignedPileUp, mergeStartIndex, isALlMode)
            // 未アサインの案件マージ
            row.unAssignedPileUp.facilities.forEach(v => {
                const vv = target.unAssignedPileUp.facilities.find(vvv => vvv.facilityId === v.facilityId)
                if (vv == undefined) return console.warn("facilityId is not exists", v.facilityId)
                mergeAndUpdate(vv, v, mergeStartIndex, isALlMode)
            })

            // 未確定のマージ
            mergeAndUpdate(target.noOrdersReceivedPileUp, row.noOrdersReceivedPileUp, mergeStartIndex, isALlMode)
            // 未確定の案件マージ
            row.noOrdersReceivedPileUp.facilities.forEach(v => {
                const vv = target.noOrdersReceivedPileUp.facilities.find(vvv => vvv.facilityId === v.facilityId)
                if (vv == undefined) return console.warn("facilityId is not exists", v.facilityId)
                mergeAndUpdate(vv, v, mergeStartIndex, isALlMode)
            })

        });
    }

    function mergeAndUpdate(target: PileUpRow, row: PileUpRow, mergeStartIndex: number, isALlMode: boolean) {
        row.labels.forEach((v, index) => {
            const targetIndex = mergeStartIndex + index;
            if (0 <= targetIndex && targetIndex < target.labels.length && target.labels[mergeStartIndex + index] != 0) {
                row.labels[index] += target.labels[targetIndex];
                // TODO: エラーをスタイルで管理しているので微妙なコード。
                const hasError = target.styles[targetIndex].color == PILEUP_DANGER_COLOR
                if(isALlMode) {
                    row.styles[index] = target.styles[targetIndex]
                } else {
                    if(hasError) {
                        row.styles[index] = {color: PILEUP_DANGER_COLOR}
                    } else {
                        row.styles[index] = {color: PILEUP_MERGE_COLOR}
                    }
                }
            }
        });
    }

}

export const usePileUps = (
    startDate: string, endDate: string, isAllMode: boolean,
    facility: Facility,
    tickets: ComputedRef<Ticket[]>, ticketUsers: ComputedRef<TicketUser[]>,
    displayType: ComputedRef<DisplayType>,
    holidays: ComputedRef<Holiday[]>, departmentList: Department[], userList: User[], facilityList: Facility[],
    defaultPileUps: PileUps[] = [],
    globalStartDate?: string,
) => {
    // ガントチャート描画用の日付だとフォーマットが違うので変換しておく
    startDate = ganttDateToYMDDate(startDate)
    endDate = ganttDateToYMDDate(endDate)
    const filteredFacilityList = facilityList.filter(v => v.status === FacilityStatus.Enabled)

    const pileUps = ref<PileUps[]>([])

    const saveFilter = () => {
        globalFilterMutation.updatePileUpsFilter(
            pileUps.value.map(v => {
                return {
                    assignedUser: v.assignedUser.display,
                    departmentId: v.departmentId,
                    display: v.display,
                    noOrdersReceivedPileUp: v.noOrdersReceivedPileUp.display,
                    unAssignedPileUp: v.unAssignedPileUp.display
                }
            })
        )
    }
    onBeforeUnmount(() => {
        saveFilter()
        window.removeEventListener("beforeunload", saveFilter)
    })
    window.addEventListener("beforeunload", saveFilter)

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
        // 初期表示時（pileUps未作成）の状態での呼び出し防止
        if (pileUps.value.length > 0) saveFilter()
        initPileUps(endDate, startDate, isAllMode, userList, pileUps, departmentList, filteredFacilityList, defaultPileUps, globalStartDate);

        // 全てのチケットから積み上げを更新する
        tickets.value.forEach(ticket => {
            setWorkHour(
                pileUps.value,
                facility,
                ticket,
                ticketUsers.value.filter(v => v.ticket_id === ticket.id),
                startDate, endDate, holidays.value, userList
            )
        })

        if (displayType.value === "week") {
            aggregatePileUpsByWeek(startDate, endDate, pileUps.value, holidays.value)
        }
        refreshPileUpByPersonExclusive = false
    }

    /**
     * １つのチケットから山積みを更新する。
     * ユーザーまたは案件の計算を行った後、サマリーへ足し上げる
     *
     * [確定かつアサイン済みの場合]
     * ユーザーとアサイン済みと部署が対象
     * 担当者の人数と営業日で均等に割り当てる。
     *
     * [確定かつ未アサインの場合]
     * 未アサインと案件が対象
     * 部署が設定されている場合のみ計上する。
     * 営業日・工数・人数で均等に割り当てる。
     *
     * [未確定の場合]
     * 未確定と案件が対象
     * 営業日・工数・人数で均等に割り当てる。
     *
     * @param pileUps
     * @param facility
     * @param ticket
     * @param ticketUsers
     * @param facilityStartDate
     * @param facilityEndDate
     * @param holidays
     * @param userList
     */
    const setWorkHour = (pileUps: PileUps[],
                         facility: Facility,
                         ticket: Ticket,
                         ticketUsers: TicketUser[],
                         facilityStartDate: string,
                         facilityEndDate: string,
                         holidays: Holiday[], userList: User[]) => {

        // TODO: validationは変更する
        if (ticket.start_date == null || ticket.end_date == null || ticket.estimate == null) {
            return
        }

        // 常に共通なもの
        const dayjsFacilityStartDate = dayjs(facilityStartDate)
        const dayjsStartDate = dayjs(ticket.start_date)
        const maxDate = dayjs(facilityEndDate).add(-1, 'days')
        const dayjsEndDate = maxDate.isBefore(dayjs(ticket.end_date)) ? maxDate : dayjs(ticket.end_date)
        const startIndex = getIndexByDate(dayjsFacilityStartDate, dayjsStartDate)
        const endIndex = getIndexByDate(dayjsFacilityStartDate, dayjsEndDate)
        const numberOfBusinessDays = getNumberOfBusinessDays(dayjsStartDate, dayjsEndDate, holidays)
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
        const estimate = ticket.estimate

        // 場合によって変化するもの
        let workHour = 0
        // 対象の取得
        const rows: PileUpRow[] = []
        const summaryRows: PileUpRow[][] = []
        const summaryLimit: number[] = []
        const userErrorFunc = (v: number) => {
            return pileUpLabelFormat(v) > 1
        }
        const facilityErrorFunc = (v: number) => {
            return false
        }
        let rowErrorFunc = userErrorFunc

        // 確定の場合
        if (facility.type === FacilityType.Ordered) {
            const ticketUserIds = ticketUsers.map(v => v.user_id)
            // アサイン済み
            if (ticketUserIds.length > 0) {
                workHour = ticket.estimate / numberOfBusinessDays / ticketUsers.length
                // ユーザーに紐づく部署を設定
                ticketUserIds.forEach(userId => {
                    const user = userList.find(user => user.id === userId)
                    if (user == undefined) return console.warn("user is not exists", userId)
                    const pileUp = pileUps.find(pileUp => pileUp.departmentId === user.department_id)!
                    if (pileUp == undefined) return console.warn("departmentId is not exists", user.department_id)
                    rows.push(pileUp.assignedUser.users.find(v => v.user.id === user.id)!)
                    summaryRows.push([pileUp, pileUp.assignedUser])
                    summaryLimit.push(userList.filter(v => v.department_id === user.department_id).length)
                })
            } else {
                workHour = ticket.estimate / numberOfBusinessDays
                // チケットに紐づく部署を設定
                const pileUp = pileUps.find(pileUp => pileUp.departmentId === ticket.department_id)
                if (pileUp == undefined) return console.warn("departmentId is not exists", ticket.department_id)
                rows.push(pileUp.unAssignedPileUp.facilities.find(v => v.facilityId === facility.id)!)
                // TODO: サマリーの計上先にするかヒアリング。絵面的には含めないっぽい。含める場合はpileUp事態をsummaryRowsに追加する
                summaryRows.push([pileUp.unAssignedPileUp])
                summaryLimit.push(9999) // TODO: あり得ない上限にしてエラーにならないようにしている。
                rowErrorFunc = facilityErrorFunc
            }
        } else if (facility.type === FacilityType.Prepared) {
            workHour = ticket.estimate / numberOfBusinessDays
            const pileUp = pileUps.find(pileUp => pileUp.departmentId === ticket.department_id)
            if (pileUp == undefined) return console.warn("departmentId is not exists", ticket.department_id)
            rows.push(pileUp.noOrdersReceivedPileUp.facilities.find(v => v.facilityId === facility.id)!)
            summaryRows.push([pileUp, pileUp.noOrdersReceivedPileUp])
            summaryLimit.push(userList.filter(v => v.department_id === pileUp.departmentId).length)
            rowErrorFunc = facilityErrorFunc
        }
        allocateWorkingHours(rows, summaryRows, rowErrorFunc, summaryLimit, validIndexes, estimate, workHour)
    }
    /**
     * 作業時間の割り当てを実施する。
     * summaryRowsはrowsのi番目の結果と同じ作業時間を割り当てる
     * @param rows
     * @param summaryRows
     * @param rowErrorFunc
     * @param summaryLimit
     * @param validIndexes
     * @param estimate
     * @param workHour
     */
    const allocateWorkingHours = (rows: PileUpRow[], summaryRows: PileUpRow[][], rowErrorFunc: (v: number) => boolean, summaryLimit: number[], validIndexes: number[], estimate: number, workHour: number) => {
        validIndexes.forEach((validIndex, index) => {
            rows.forEach((v, i) => {
                if (estimate < 0) {
                    return
                }
                let workHourResult = 0
                if (estimate - workHour < 0) {
                    workHourResult = estimate
                } else {
                    workHourResult = workHour
                }
                v.labels[validIndex] += workHourResult
                if (rowErrorFunc(v.labels[validIndex])) {
                    v.styles[validIndex] = {color: PILEUP_DANGER_COLOR}
                }
                // ユーザーごとに割当先が異なることに注意
                summaryRows[i].forEach(pileUpRow => {
                    pileUpRow.labels[validIndex] += workHourResult
                    if (pileUpLabelFormat(pileUpRow.labels[validIndex]) > summaryLimit[i]) {
                        pileUpRow.styles[validIndex] = {color: PILEUP_DANGER_COLOR}
                    }
                })
                estimate -= workHour
            })
        })
    }

    // 日付からどのindexに該当するか取得する
    const getIndexByDate = (facilityStartDate: Dayjs, date: Dayjs) => {
        return date.diff(facilityStartDate, 'days')
    }
    // TODO: 週表示, 祝日が案件ごとに違うので見る案件によっては結果が変わってしまう。
    const aggregatePileUpsByWeek = (term_from: string, term_to: string, pileUps: PileUps[], holidays: Holiday[]) => {
        // 日毎で計算された結果を週ごとに集約する。
        // 集約元のIndexと集約先のIndexをマッピングする
        const dayjsHolidays = holidays.map(v => dayjs(v.date))
        const dayjsEndDate = dayjs(term_to)
        let currentDate = dayjs(term_from)
        let currentStartOfWeek = dayjs(term_from).startOf("week")
        let currentWeekIndex = 0
        let currentDayIndex = 0
        const indexMap: { [index: number]: number[] } = {}
        const holidayMap: { [index: number]: number } = {}
        while (currentDate.isBefore(dayjsEndDate)) {
            // 週が異なっていれば次の週とする
            if (!currentStartOfWeek.isSame(currentDate.startOf("week"))) {
                currentWeekIndex++
                currentStartOfWeek = currentDate.startOf("week")
            }
            // 週のマップがなければ初期化
            if (indexMap[currentWeekIndex] == null) {
                indexMap[currentWeekIndex] = []
                holidayMap[currentWeekIndex] = 0
            }
            // その週に紐づく日付のindexを追加する
            indexMap[currentWeekIndex].push(currentDayIndex)
            // 祝日だった場合祝日の数を増やす
            if (dayjsHolidays.find(v => v.isSame(currentDate))) holidayMap[currentWeekIndex]++
            // シーケンスを進める
            currentDate = currentDate.add(1, "day")
            currentDayIndex++
        }

        // 元データから週データを計算する
        const result: PileUps[] = []
        pileUps.forEach(v => {
            const {labels: departmentLabels, styles: departmentStyles} = aggregateToWeek(v, indexMap, holidayMap)
            const {
                labels: assignedUserLabels,
                styles: assignedUserStyles
            } = aggregateToWeek(v.assignedUser, indexMap, holidayMap)
            const {
                labels: noOrdersReceivedLabels,
                styles: noOrdersReceivedStyles
            } = aggregateToWeek(v.noOrdersReceivedPileUp, indexMap, holidayMap)
            const {
                labels: unAssignedPileUpLabels,
                styles: unAssignedPileUpStyles
            } = aggregateToWeek(v.unAssignedPileUp, indexMap, holidayMap)
            result.push(
                {
                    departmentId: v.departmentId,
                    labels: departmentLabels,
                    styles: departmentStyles,
                    display: v.display,
                    assignedUser: {
                        users: v.assignedUser.users.map(pileUpByPerson => {
                            const {
                                labels: userLabels,
                                styles: userStyles
                            } = aggregateToWeek(pileUpByPerson, indexMap, holidayMap)
                            return {hasError: false, labels: userLabels, styles: userStyles, user: pileUpByPerson.user}
                        }),
                        labels: assignedUserLabels,
                        styles: assignedUserStyles,
                        display: v.assignedUser.display,
                    },
                    noOrdersReceivedPileUp: {
                        facilities: v.noOrdersReceivedPileUp.facilities.map(pileUpByFacility => {
                            const {
                                labels: userLabels,
                                styles: userStyles
                            } = aggregateToWeek(pileUpByFacility, indexMap, holidayMap)
                            return {
                                hasError: false,
                                labels: userLabels,
                                styles: userStyles,
                                facilityId: pileUpByFacility.facilityId
                            }
                        }),
                        labels: noOrdersReceivedLabels,
                        styles: noOrdersReceivedStyles,
                        display: v.assignedUser.display,
                    },
                    unAssignedPileUp: {
                        facilities: v.unAssignedPileUp.facilities.map(pileUpByFacility => {
                            const {
                                labels: userLabels,
                                styles: userStyles
                            } = aggregateToWeek(pileUpByFacility, indexMap, holidayMap)
                            return {
                                hasError: false,
                                labels: userLabels,
                                styles: userStyles,
                                facilityId: pileUpByFacility.facilityId
                            }
                        }),
                        labels: unAssignedPileUpLabels,
                        styles: unAssignedPileUpStyles,
                        display: v.assignedUser.display,
                    },
                }
            )
        })

        // 結果を週に差し替える
        pileUps.length = 0
        pileUps.push(...result)
    }
    // 週に集約する関数
    const aggregateToWeek = (row: PileUpRow, indexMap: { [index: number]: number[] }, holidayMap: {
        [index: number]: number
    }) => {
        const result: PileUpRow = {labels: [], styles: []}
        let hasError = false
        let hasMerge = false
        for (const key in indexMap) {
            // 稼働時間を集約
            const workHour = indexMap[key].reduce((p, c) => {
                return p + row.labels[c]
            }, 0)
            // 週表示の場合は 営業日 * 8 を 1とする
            result.labels.push(workHour / (8 * (indexMap[key].length - holidayMap[key])))
            hasError = indexMap[key].some(vv => row.styles[vv].color == PILEUP_DANGER_COLOR)
            hasMerge = indexMap[key].some(vv => row.styles[vv].color == PILEUP_MERGE_COLOR)
            // エラー > マージ > デフォルト となるようにスタイルを決定する。
            result.styles.push(hasError ? {color: PILEUP_DANGER_COLOR} : hasMerge ? {color: PILEUP_MERGE_COLOR} : {})
        }
        return {labels: result.labels, styles: result.styles}
    }

    refreshPileUps()

    return {
        // pileUpFilters,
        // pileUpsByDepartment,
        // pileUpsByPerson,
        // displayPileUps,
        pileUps,
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
    isAllMode: boolean,
) => {
    const {facilityList, userList, departmentList, facilityTypes} = inject(GLOBAL_STATE_KEY)!
    const filteredFacilityList = facilityList.filter(v => v.status === FacilityStatus.Enabled)

    const {data: allTickets} = await Api.getAllTickets(facilityTypes)
    const {data: allTicketUsers} = await Api.getTicketUsers(allTickets.list.map(v => v.id!))
    const {data: allData} = await Api.getPileUps(excludeFacilityId, facilityTypes)

    // 全案件の最小
    const startDate: string = filteredFacilityList.slice().sort((a, b) => {
        return a.term_from > b.term_from ? 1 : -1
    }).shift()!.term_from.substring(0, 10)
    // 全案件の最大
    const endDate: string = filteredFacilityList.slice().sort((a, b) => {
        return a.term_to > b.term_to ? 1 : -1
    }).pop()!.term_to.substring(0, 10)

    const defaultPileUps = ref<PileUps[]>([])
    initPileUps(endDate, startDate, isAllMode, userList, defaultPileUps, departmentList, filteredFacilityList);

    for (const facility of filteredFacilityList) {
        // 対象外の案件はスキップ
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

        const {pileUps} = usePileUps(startDate, endDate, isAllMode,
            facility,
            tickets, ticketUsers,
            computed(() => displayType),
            holidays, departmentList, userList, filteredFacilityList, defaultPileUps.value, startDate)
        defaultPileUps.value.length = 0
        defaultPileUps.value.push(...pileUps.value)
    }
    return {
        globalStartDate: startDate,
        defaultPileUps,
    }
}

