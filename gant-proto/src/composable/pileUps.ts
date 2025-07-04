import {onBeforeUnmount, Ref, ref, StyleValue, UnwrapRef, watch} from "vue";
import dayjs, {Dayjs} from "dayjs";
import {Department, Facility, Holiday, Ticket, TicketUser, User} from "@/api";
import {dayBetween, ganttDateToYMDDate} from "@/coreFunctions/manHourCalculation";
import {Api} from "@/api/axios";
import {FacilityStatus, FacilityType} from "@/const/common";
import {DisplayType} from "@/composable/ganttFacilityMenu";
import {globalFilterGetter, globalFilterMutation} from "@/utils/globalFilterState";
import {pileUpLabelFormat} from "@/utils/filters";
import {allowed} from "@/composable/role";


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

const PILEUP_DANGER_COLOR = "rgb(255 89 89)"
const PILEUP_MERGE_COLOR = "rgb(68 141 255)"
const PILEUP_HOLIDAY_COLOR = "rgb(243 238 226)"

function initPileUps(
    endDate: string, startDate: string, isALlMode: boolean,
    userList: User[], pileUps: Ref<UnwrapRef<PileUps[]>>, departmentList: Department[], facilityList: Facility[], defaultValidUserIndexMap: Map<number, number[]>,
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
            styles: Array(duration).fill(0).map(() => ({})),
            display: filter ? filter.display : false,
            assignedUser: {
                labels: Array(duration).fill(0),
                styles: Array(duration).fill(0).map(() => ({})),
                display: filter ? filter.assignedUser : false,
                users: users.map(user => {
                    return {
                        user: user,
                        labels: Array(duration).fill(0),
                        styles: Array(duration).fill(0).map(() => ({})),
                        hasError: false,
                    }
                })
            },
            unAssignedPileUp: {
                labels: Array(duration).fill(0),
                styles: Array(duration).fill(0).map(() => ({})),
                display: filter ? filter.unAssignedPileUp : false,
                facilities: orderedFacilities.map(facility => {
                    return {
                        facilityId: facility.id!,
                        labels: Array(duration).fill(0),
                        styles: Array(duration).fill(0).map(() => ({})),
                        hasError: false,
                    }
                })
            },
            noOrdersReceivedPileUp: {
                labels: Array(duration).fill(0),
                styles: Array(duration).fill(0).map(() => ({})),
                display: filter ? filter.noOrdersReceivedPileUp : false,
                facilities: noOrderedFacilities.map(facility => {
                    return {
                        facilityId: facility.id!,
                        labels: Array(duration).fill(0),
                        styles: Array(duration).fill(0).map(() => ({})),
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
        pileUps.value.forEach((row) => {
            // 部署のマージ
            const target = defaultPileUps.find(v => v.departmentId === row.departmentId);
            if (target == undefined) return console.warn("departmentId is not exists", row.departmentId)
            mergeAndUpdate(target, row, mergeStartIndex, isALlMode, false);

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
        // 在籍確認Mapを詰める
        const newMap = new Map<number, number[]>();
        let currentKey = 0;
        for (const [, value] of Array.from(defaultValidUserIndexMap.entries()).slice(mergeStartIndex)) {
            newMap.set(currentKey++, value);
        }
        return newMap
    }

    function mergeAndUpdate(target: PileUpRow, row: PileUpRow, mergeStartIndex: number, isALlMode: boolean, mergedColor = true) {
        row.labels.forEach((v, index) => {
            const targetIndex = mergeStartIndex + index;
            if (0 <= targetIndex && targetIndex < target.labels.length) {
                if (target.labels[targetIndex] != 0) {
                    row.labels[index] += target.labels[targetIndex];
                    if (mergedColor) {
                        row.styles[index] = {color: PILEUP_MERGE_COLOR}
                    }
                }
                // TODO: エラーをスタイルで管理しているので微妙なコード。
                const hasError = target.styles[targetIndex].color == PILEUP_DANGER_COLOR
                if (isALlMode) {
                    row.styles[index] = target.styles[targetIndex]
                } else {
                    if (hasError) {
                        row.styles[index] = {color: PILEUP_DANGER_COLOR}
                    }
                }
                // 在籍期間のグレーアウトを設定する
                if (target.styles[targetIndex]["background-color"] != undefined) {
                    row.styles[index]["background-color"] = target.styles[targetIndex]["background-color"]
                }
            }
        });
    }

}

// NOTE: displayType は単なるストリングなので computedRefで渡さないと watch が機能しない
export const usePileUps = (
    startDate: string, endDate: string, isAllMode: boolean,
    facility: Facility,
    tickets: Ref<Ticket[]>, ticketUsers: Ref<TicketUser[]>,
    displayType: Ref<"day" | "week" | "hour" | "month">,
    holidays: Ref<Holiday[]>, departmentList: Department[], userList: User[], facilityList: Facility[],
    defaultValidUserIndexMap: Map<number, number[]>,
    defaultPileUps?: PileUps[],
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
        console.log("############## function refreshPileUps", refreshPileUpByPersonExclusive)
        if (refreshPileUpByPersonExclusive) {
            return
        } else {
            refreshPileUpByPersonExclusive = true
        }
        // 初期表示時（pileUps未作成）の状態での呼び出し防止
        if (pileUps.value.length > 0) saveFilter()
        const validUserIndexMap = initPileUps(endDate, startDate, isAllMode, userList, pileUps, departmentList, filteredFacilityList, defaultValidUserIndexMap, defaultPileUps, globalStartDate);

        // 全てのチケットから積み上げを更新する
        tickets.value.forEach(ticket => {
            setWorkHour(
                pileUps.value,
                facility,
                ticket,
                ticketUsers.value.filter(v => v.ticket_id === ticket.id),
                startDate, endDate, holidays.value, userList, validUserIndexMap!
            )
        })

        if (displayType.value === "week") {
            aggregatePileUpsByWeek(startDate, endDate, pileUps.value, holidays.value)
        }

        // 在籍期間による非表示の処理を行う。
        const excludeUserIds: number[] = []
        if (facility) {
            // 案件ビューの場合はその案件の期間内であれば表示
            excludeUserIds.push(...userList.filter(user => {
                const userEmployStart = new Date(user.employment_start_date);
                const userEmployEnd = user.employment_end_date ? new Date(user.employment_end_date) : Infinity; // If null, set to Infinity
                const facilityTermStart = new Date(facility.term_from);
                const facilityTermEnd = new Date(facility.term_to);
                return userEmployStart > facilityTermEnd || userEmployEnd < facilityTermStart;
            }).map(v => v.id!))
        } else {
            // 全体ビューであれば、在籍期間が過去であれば非表示
            const currentDate = new Date()
            currentDate.setHours(0, 0, 0, 0)
            excludeUserIds.push(...userList.filter(user => {
                const userEmployEnd = user.employment_end_date ? new Date(user.employment_end_date) : Infinity; // If null, set to Infinity
                return userEmployEnd < currentDate;
            }).map(v => v.id!))
        }

        pileUps.value.forEach(pileUp => {
            pileUp.assignedUser.users = pileUp.assignedUser.users.filter(userPileUp => {
                return !excludeUserIds.includes(userPileUp.user.id!)
            })
        })


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
     * @param defaultValidUserIndexMap
     */
    const setWorkHour = (pileUps: PileUps[],
                         facility: Facility,
                         ticket: Ticket,
                         ticketUsers: TicketUser[],
                         facilityStartDate: string,
                         facilityEndDate: string,
                         holidays: Holiday[], userList: User[], defaultValidUserIndexMap: Map<number, number[]>) => {

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

        // NOTE: APIのget_default_pile_ups.go と必ず合わせること
        const ticketUserIds = ticketUsers.map(v => v.user_id)
        let workPerDay: number
        if (ticketUserIds.length > 0) {
            const numberOfWorkerByDay = validIndexes.reduce((sum, key) => {
                return sum + (defaultValidUserIndexMap.get(key)?.filter(v => ticketUserIds.includes(v)).length || 0)
            }, 0)
            workPerDay = ticket.estimate / numberOfWorkerByDay
        } else if (ticket.number_of_worker != undefined && ticket.number_of_worker > 0) {
            workPerDay = ticket.estimate / validIndexes.length / ticket.number_of_worker
        } else {
            return
        }

        /**
         * 集計とエラーの仕様を整理する
         * ・担当有
         *  ・ユーザー：ユーザーは1より大きい、部署、アサイン済み
         *   ・部署：部署人数計
         *・確定
         *   ・アサイン済み：部署人数計
         * ・担当無
         *  ・確定
         *   ・未アサインの設備：エラー無
         *   ・部署：部署人数系
         *   ・未アサイン：部署人数系
         * ・未確定（担当の有無は関係ない）
         *   ・未確定の設備：エラー無
         *   ・部署：部署人数系
         *   ・未確定：部署人数計
         */
        if (ticketUserIds.length > 0) {
            // アサイン済みの場合 [担当者] [積み上げ、アサイン済み]に計上する
            validIndexes.forEach(validIndex => {
                ticketUserIds.forEach(userId => {
                    // アサイン済みの場合はユーザーから対象のpileUpsを指定する（部署が指定されていないケースがあるため）
                    const user = userList.find(user => user.id === userId)
                    if (user == undefined) return console.warn("user is not exists", userId)
                    const targetPileUp = pileUps.find(pileUp => pileUp.departmentId === user.department_id)!
                    if (targetPileUp == undefined) return console.warn("departmentId is not exists", user.department_id)

                    // 設備の期間外の場合は処理を中断
                    if (targetPileUp.labels.length <= validIndex) {
                        return
                    }

                    // 在籍期間外の場合は対象外とする。
                    if (!(defaultValidUserIndexMap.get(validIndex)?.includes(userId))) {
                        return
                    }

                    // 部署人数合計 validUserMapで既にRoleは絞り込まれているので特に対応無
                    const numberOfDepartmentUsers = userList.filter(v => v.department_id === user.department_id).filter(v => defaultValidUserIndexMap.get(validIndex)?.includes(v.id!)).length

                    // 部署への積み上げ(共通)
                    targetPileUp.labels[validIndex] += workPerDay
                    if (pileUpLabelFormat(targetPileUp.labels[validIndex]) > numberOfDepartmentUsers) {
                        targetPileUp.styles[validIndex] = {color: PILEUP_DANGER_COLOR}
                    }

                    if (facility.type == FacilityType.Ordered) {
                        // アサイン済みへの積み上げ
                        targetPileUp.assignedUser.labels[validIndex] += workPerDay
                        if (pileUpLabelFormat(targetPileUp.assignedUser.labels[validIndex]) > numberOfDepartmentUsers) {
                            targetPileUp.assignedUser.styles[validIndex] = {color: PILEUP_DANGER_COLOR}
                        }
                        // ユーザーの足し上げ処理
                        const targetUserPileUp = targetPileUp.assignedUser.users.find(v => v.user.id == userId)!
                        targetUserPileUp.labels[validIndex] += workPerDay
                        if (pileUpLabelFormat(targetUserPileUp.labels[validIndex]) > 1) {
                            targetUserPileUp.styles[validIndex] = {color: PILEUP_DANGER_COLOR}
                        }
                    } else {
                        // 未確定の場合 への積み上げ
                        targetPileUp.noOrdersReceivedPileUp.labels[validIndex] += workPerDay
                        if (pileUpLabelFormat(targetPileUp.noOrdersReceivedPileUp.labels[validIndex]) > numberOfDepartmentUsers) {
                            targetPileUp.noOrdersReceivedPileUp.styles[validIndex] = {color: PILEUP_DANGER_COLOR}
                        }

                        // 未確定の場合の設備への積み上げ
                        const targetFacilityPileUp = targetPileUp.noOrdersReceivedPileUp.facilities.find(v => v.facilityId === facility.id)!
                        targetFacilityPileUp.labels[validIndex] += workPerDay
                    }
                })
            })
        } else {
            // 未アサインの場合は部署の設定がされていないと計上しない
            const targetPileUp = pileUps.find(pileUp => ticket.department_id != undefined && pileUp.departmentId === ticket.department_id)!
            if (targetPileUp == undefined) return console.warn("departmentId is not exists", ticket.department_id)

            // 未アサインの場合 [未アサインのその設備] [積み上げ、未アサイン積み上げ]に計上する
            validIndexes.forEach(validIndex => {
                // 設備の期間外の場合は処理を中断
                if (targetPileUp.labels.length <= validIndex) {
                    return
                }

                const numberOfDepartmentUsers = userList.filter(v => v.department_id === targetPileUp.departmentId).filter(v => defaultValidUserIndexMap.get(validIndex)?.includes(v.id!)).length

                // 部署への積み上げ（共通）
                targetPileUp.labels[validIndex] += workPerDay
                if (pileUpLabelFormat(targetPileUp.labels[validIndex]) > numberOfDepartmentUsers) {
                    targetPileUp.styles[validIndex] = {color: PILEUP_DANGER_COLOR}
                }
                if (facility.type === FacilityType.Ordered) {
                    // 未アサインへ積み上げ
                    targetPileUp.unAssignedPileUp.labels[validIndex] += workPerDay
                    if (pileUpLabelFormat(targetPileUp.unAssignedPileUp.labels[validIndex]) > numberOfDepartmentUsers) {
                        targetPileUp.unAssignedPileUp.styles[validIndex] = {color: PILEUP_DANGER_COLOR}
                    }

                    // 未アサインの設備へ積み上げ
                    const targetFacilityPileUp = targetPileUp.unAssignedPileUp.facilities.find(v => v.facilityId === facility.id)!
                    targetFacilityPileUp.labels[validIndex] += workPerDay
                } else {
                    // 未確定の場合 への積み上げ
                    targetPileUp.noOrdersReceivedPileUp.labels[validIndex] += workPerDay
                    if (pileUpLabelFormat(targetPileUp.noOrdersReceivedPileUp.labels[validIndex]) > numberOfDepartmentUsers) {
                        targetPileUp.noOrdersReceivedPileUp.styles[validIndex] = {color: PILEUP_DANGER_COLOR}
                    }
                    // 設備への積み上げ
                    const targetFacilityPileUp = targetPileUp.noOrdersReceivedPileUp.facilities.find(v => v.facilityId === facility.id)!
                    targetFacilityPileUp.labels[validIndex] += workPerDay
                }
            })
        }
        // 旧資産をメモとして残しておくgithubつかえやってはなしではあるが。
        // // 確定の場合
        // if (facility.type === FacilityType.Ordered) {
        //     const ticketUserIds = ticketUsers.map(v => v.user_id)
        //     // アサイン済み
        //     if (ticketUserIds.length > 0) {
        //         // 総稼働日数を計算
        //         // 予定工数 / 総稼働日数で分配する
        //
        //         // ユーザーに紐づく部署を設定
        //         ticketUserIds.forEach(userId => {
        //             const user = userList.find(user => user.id === userId)
        //             if (user == undefined) return console.warn("user is not exists", userId)
        //             const pileUp = pileUps.find(pileUp => pileUp.departmentId === user.department_id)!
        //             if (pileUp == undefined) return console.warn("departmentId is not exists", user.department_id)
        //             rows.push(pileUp.assignedUser.users.find(v => v.user.id === user.id)!)
        //             summaryRows.push([pileUp, pileUp.assignedUser])
        //             summaryLimit.push(userList.filter(v => v.department_id === user.department_id).length)
        //         })
        //     } else {
        //         workHour = ticket.estimate / numberOfBusinessDays / (ticket.number_of_worker ?? 1)
        //         // チケットに紐づく部署を設定
        //         const pileUp = pileUps.find(pileUp => pileUp.departmentId === ticket.department_id)
        //         if (pileUp == undefined) return console.warn("departmentId is not exists", ticket.department_id)
        //         rows.push(pileUp.unAssignedPileUp.facilities.find(v => v.facilityId === facility.id)!)
        //         // TODO: サマリーの計上先にするかヒアリング。絵面的には含めないっぽい。含める場合はpileUp事態をsummaryRowsに追加する
        //         summaryRows.push([pileUp, pileUp.unAssignedPileUp])
        //         summaryLimit.push(userList.filter(v => v.department_id === pileUp.departmentId).length, 9999) // TODO: あり得ない上限にしてエラーにならないようにしている。
        //         rowErrorFunc = facilityErrorFunc
        //     }
        // } else if (facility.type === FacilityType.Prepared) {
        //     workHour = ticket.estimate / numberOfBusinessDays / (ticket.number_of_worker ?? 1)
        //     const pileUp = pileUps.find(pileUp => pileUp.departmentId === ticket.department_id)
        //     if (pileUp == undefined) return console.warn("departmentId is not exists", ticket.department_id)
        //     rows.push(pileUp.noOrdersReceivedPileUp.facilities.find(v => v.facilityId === facility.id)!)
        //     summaryRows.push([pileUp, pileUp.noOrdersReceivedPileUp])
        //     summaryLimit.push(userList.filter(v => v.department_id === pileUp.departmentId).length)
        //     rowErrorFunc = facilityErrorFunc
        // }
        // allocateWorkingHours(rows, summaryRows, rowErrorFunc, summaryLimit, validIndexes, estimate, workHour)
    }
    // /**
    //  * 作業時間の割り当てを実施する。
    //  * summaryRowsはrowsのi番目の結果と同じ作業時間を割り当てる
    //  * @param rows
    //  * @param summaryRows
    //  * @param rowErrorFunc
    //  * @param summaryLimit
    //  * @param validIndexes
    //  * @param estimate
    //  * @param workHour
    //  */
    // const allocateWorkingHours = (rows: PileUpRow[], summaryRows: PileUpRow[][], rowErrorFunc: (v: number) => boolean, summaryLimit: number[], validIndexes: number[], estimate: number, workHour: number) => {
    //     validIndexes.forEach((validIndex) => {
    //         rows.forEach((v, i) => {
    //             if (estimate < 0) {
    //                 return
    //             }
    //             let workHourResult = 0
    //             if (estimate - workHour < 0) {
    //                 workHourResult = estimate
    //             } else {
    //                 workHourResult = workHour
    //             }
    //             v.labels[validIndex] += workHourResult
    //             if (rowErrorFunc(v.labels[validIndex])) {
    //                 v.styles[validIndex] = {color: PILEUP_DANGER_COLOR}
    //             }
    //             // ユーザーごとに割当先が異なることに注意
    //             summaryRows[i].forEach(pileUpRow => {
    //                 pileUpRow.labels[validIndex] += workHourResult
    //                 if (pileUpLabelFormat(pileUpRow.labels[validIndex]) > summaryLimit[i]) {
    //                     pileUpRow.styles[validIndex] = {color: PILEUP_DANGER_COLOR}
    //                 }
    //             })
    //             estimate -= workHour
    //         })
    //     })
    // }

    // 日付からどのindexに該当するか取得する
    const getIndexByDate = (facilityStartDate: Dayjs, date: Dayjs) => {
        return date.diff(facilityStartDate, 'days')
    }
    // TODO: 週表示, 祝日が案件ごとに違うので見る案件によっては結果が変わってしまう。
    const aggregatePileUpsByWeek = (term_from: string, term_to: string, pileUps: PileUps[], holidays: Holiday[]) => {
        // 日毎で計算された結果を週ごとに集約する。
        // 集約元のIndexと集約先のIndexをマッピングする
        // const dayjsHolidays = holidays.map(v => dayjs(v.date))
        const dayjsEndDate = dayjs(term_to)
        let currentDate = dayjs(term_from)
        let currentStartOfWeek = dayjs(term_from).startOf("week")
        let currentWeekIndex = 0
        let currentDayIndex = 0
        const indexMap: { [index: number]: number[] } = {}
        const holidayMap: { [index: number]: number } = {}

        const padding = (v: number) => {
            return (v + "").padStart(2, '0');
        }
        const dayjsToString = (v: Dayjs) => {
            return v.year() + "-" + padding(v.month() + 1) + "-" + padding(v.date())
        }

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
            // NOTE:パフォーマンス対応のために文字列での比較を行っている。そもそもHolidayをDate型で定義すべきだった。
            if (holidays.find(v => v.date.substring(0, 10) == dayjsToString(currentDate))) holidayMap[currentWeekIndex]++
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
            const numOfWorkdays = (indexMap[key].length - holidayMap[key])
            // その週の稼働日が0の場合は計算しない(NaN, Infinity対策)
            if (numOfWorkdays <= 0) {
                result.labels.push(0)
            } else {
                result.labels.push(workHour / (8 * numOfWorkdays))
            }
            hasError = indexMap[key].some(vv => row.styles[vv].color == PILEUP_DANGER_COLOR)
            hasMerge = indexMap[key].some(vv => row.styles[vv].color == PILEUP_MERGE_COLOR)
            // エラー > マージ > デフォルト となるようにスタイルを決定する。
            const style: StyleValue = hasError ? {color: PILEUP_DANGER_COLOR} : hasMerge ? {color: PILEUP_MERGE_COLOR} : {}

            // 全ての日が在籍期間外の場合は休みとする。
            if (indexMap[key].every(v => row.styles[v]["background-color"] != undefined)) {
                style["background-color"] = PILEUP_HOLIDAY_COLOR
            }

            result.styles.push(style)
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
    facilityTypes: string[]) => {
    // const {facilityList, userList, departmentList, facilityTypes} = inject(GLOBAL_STATE_KEY)!
    // const filteredFacilityList = facilityList.filter(v => v.status === FacilityStatus.Enabled)
    //
    // const {data: allTickets} = await Api.getAllTickets(facilityTypes)
    // const {data: allTicketUsers} = await Api.getTicketUsers(allTickets.list.map(v => v.id!))
    // const {data: allData} = await Api.getPileUps(excludeFacilityId, facilityTypes)
    //
    // // 全案件の最小
    // const startDate: string = filteredFacilityList.slice().sort((a, b) => {
    //     return a.term_from > b.term_from ? 1 : -1
    // }).shift()!.term_from.substring(0, 10)
    // // 全案件の最大
    // const endDate: string = filteredFacilityList.slice().sort((a, b) => {
    //     return a.term_to > b.term_to ? 1 : -1
    // }).pop()!.term_to.substring(0, 10)
    //
    // const defaultPileUps = ref<PileUps[]>([])
    // initPileUps(endDate, startDate, isAllMode, userList, defaultPileUps, departmentList, filteredFacilityList);
    //
    // for (const facility of filteredFacilityList) {
    //     // 対象外の案件はスキップ
    //     if (facility.id === excludeFacilityId) continue
    //     const targetData = allData.list.find(v => v.facilityId === facility.id)
    //     // 対象データなしの場合もスキップ
    //     if (targetData == undefined) continue
    //
    //     const innerHolidays: Holiday[] = []
    //     innerHolidays.push(...targetData.holidays)
    //
    //     const holidays = computed(() => {
    //         return innerHolidays
    //     })
    //     const tickets = computed(() => {
    //         return allTickets.list.filter(v => targetData.ganttGroups.map(vv => vv.id!).includes(v.gantt_group_id))
    //     })
    //     const ticketUsers = computed(() => {
    //         return allTicketUsers.list.filter(v => tickets.value.map(vv => vv.id).includes(v.ticket_id))
    //     })
    //
    //     const {pileUps} = usePileUps(startDate, endDate, isAllMode,
    //         facility,
    //         tickets, ticketUsers,
    //         computed(() => displayType),
    //         holidays, departmentList, userList, filteredFacilityList, defaultPileUps.value, startDate)
    //     defaultPileUps.value.length = 0
    //     defaultPileUps.value.push(...pileUps.value)
    // }
    // return {
    //     globalStartDate: startDate,
    //     defaultPileUps,
    // }

    // TODO: どうなんだろうこれ。APIでは呼び出せない権限なのでUIガワで対応する。
    if (allowed("VIEW_PILEUPS")) {
        const {data} = await Api.getDefaultPileUps(excludeFacilityId, isAllMode, facilityTypes)

        // styleの適応を実施する
        return {
            globalStartDate: data.globalStartDate,
            defaultPileUps: data.defaultPileUps,
            defaultValidUserIndexMap: new Map(data.defaultValidUserIndexes.map(v => [v.ValidIndex, v.UserIds])),
        }
    } else {
        // styleの適応を実施する
        return {
            globalStartDate: "",
            defaultPileUps: [],
            defaultValidUserIndexMap: new Map(),
        }
    }
}

