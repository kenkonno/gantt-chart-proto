import dayjs, {Dayjs} from "dayjs";
import {GanttBarObject, MileStone} from "@infectoone/vue-ganttastic";
import {computed, inject, ref, toValue} from "vue";
import {GanttGroup, OperationSetting, Ticket, TicketUser, User} from "@/api";
import {Api} from "@/api/axios";
import {useGanttGroupTable} from "@/composable/ganttGroup";
import {useTicketTable} from "@/composable/ticket";
import {useTicketUserTable} from "@/composable/ticketUser";
import {useFacility} from "@/composable/facility";
import {changeSort} from "@/utils/sort";
import {GLOBAL_ACTION_KEY, GLOBAL_STATE_KEY} from "@/composable/globalState";
import {
    addBusinessDays,
    endOfDay,
    ganttDateToYMDDate,
    getEndDateByRequiredBusinessDay,
    getNumberOfBusinessDays, YMDDateToGanttEndDate, YMDDateToGanttStartDate
} from "@/coreFunctions/manHourCalculation";
import {DAYJS_FORMAT} from "@/utils/day";
import {ApiMode, DEFAULT_PROCESS_COLOR} from "@/const/common";
import {DisplayType} from "@/composable/ganttFacilityMenu";
import {GLOBAL_DEPARTMENT_USER_FILTER_KEY} from "@/composable/departmentUserFilter";
import {allowed} from "@/composable/role";
import {useMilestoneTable} from "@/composable/milestone";
import {getUserInfo} from "@/composable/auth";
import Swal from "sweetalert2";
import {globalFilterGetter, globalFilterMutation} from "@/utils/globalFilterState";

export type GanttChartGroup = {
    ganttGroup: GanttGroup // TODO: 結局ganttRowにganttGroup設定しているから設計としては微妙っぽい。
    rows: GanttRow[]
    unitId: number
}

export type GanttRow = {
    bar: GanttBarObject
    ganttGroup?: GanttGroup
    ticket?: Ticket
    ticketUsers?: TicketUser[]
}
export type UnitGroupInfo = {
    unitId: number
    isOpen: boolean
}
const BAR_NORMAL_COLOR = "rgb(147 206 255)"
const BAR_COMPLETE_COLOR = "rgb(200 200 200)"

function getScheduledOperatingHours(operationSettings: OperationSetting[], row: GanttRow) {
    return operationSettings.filter(operationSetting => {
        return operationSetting.facility_id === row.ganttGroup?.facility_id &&
            operationSetting.unit_id === row.ganttGroup?.unit_id
    }).reduce((accumulateValue, currentValue) => {
        const workHours = currentValue.workHours.find(v => v.process_id === row.ticket?.process_id)
        return accumulateValue + workHours!.work_hour!
    }, 0);
}

export async function useGanttFacility() {
    // injectはsetupと同期的に呼び出す必要あり
    const {currentFacilityId} = inject(GLOBAL_STATE_KEY)!
    const {getScheduleAlert} = inject(GLOBAL_ACTION_KEY)!
    const {selectedDepartment, selectedUser} = inject(GLOBAL_DEPARTMENT_USER_FILTER_KEY)!
    const {userList, processList, departmentList, holidayList, operationSettingMap, unitMap} = inject(GLOBAL_STATE_KEY)!
    const {facility} = await useFacility(currentFacilityId)
    const chartStart = ref(dayjs(facility.value.term_from).format(DAYJS_FORMAT))
    const chartEnd = ref(dayjs(facility.value.term_to).format(DAYJS_FORMAT))
    const milestones = ref<MileStone[]>([])

    const getUnitName = computed(() => (id: number) => {
        return unitMap[currentFacilityId].find(v => v.id === id)?.name
    })
    const getDepartmentName = computed(() => (id: number) => {
        return departmentList.find(v => v.id === id)?.name
    })
    const getProcessName = (id: number) => {
        return processList.find(v => v.id === id)?.name
    }
    const getProcessColor = (id?: number | null) => {
        if (id == null) {
            return DEFAULT_PROCESS_COLOR
        }
        return processList.find(v => v.id === id)?.color
    }
    const getOperationList = computed(() => {
        return operationSettingMap[currentFacilityId]
    })
    const getHolidaysForGantt = (displayType: DisplayType) => {
        // TODO: 稼働日対応（祝日マスタと稼働日マスタを組み合わせる）
        if (displayType === "day") {
            return holidayList.map(v => new Date(v.date))
        } else {
            return []
        }
    }
    const getHolidays = computed(() => {
        // TODO: 稼働日対応（祝日マスタと稼働日マスタを組み合わせる）
        return holidayList
    })
    const getGanttChartWidth = (displayType: DisplayType) => {
        // 1日30pxとして計算する
        return (dayjs(facility.value.term_to).diff(dayjs(facility.value.term_from), displayType) + 1) * 30 + "px"
    }

    // 積み上げに渡すよう
    const getTickets = computed(() => {
        return ticketList
    })

    // ユニットのグルーピング化情報
    const unitGroupInfo = ref<UnitGroupInfo[]>([])
    unitMap[currentFacilityId].forEach(unit => {
        const savedUnitGroupInfo = globalFilterGetter.getUnitGroupInfo()
        unitGroupInfo.value.push({
            unitId: unit.id!,
            isOpen: savedUnitGroupInfo?.[unit.id!] ?? true,
        })
    })
    const isOpenUnit = (unitId: number) => {
        return unitGroupInfo.value.find(v => v.unitId == unitId)?.isOpen ?? false
    }
    const toggleUnitOpen = (unitId: number, forceStatus?: boolean) => {
        const target = unitGroupInfo.value.find(v => v.unitId == unitId)
        if (target) {
            target.isOpen = !target.isOpen
            if (forceStatus != null) {
                target.isOpen = forceStatus
            }

            // 変更を保存する
            const savedUnitGroupInfo = globalFilterGetter.getUnitGroupInfo() || {}
            savedUnitGroupInfo[unitId] = target.isOpen

            // updateUnitGroupInfo関数を使用して保存
            globalFilterMutation.updateUnitGroupInfo(savedUnitGroupInfo)
        }
        refreshBars()
    }
    const isAllOpenUnit = computed(() => {
        return unitGroupInfo.value.every(v => v.isOpen)
    })
    const toggleAllUnitOpen = () => {
        const dst = !isAllOpenUnit.value
        return unitGroupInfo.value.forEach(v => toggleUnitOpen(v.unitId, dst))
    }


    // ganttGroupIdに紐づくUnitの開閉状況によりクラスを返却する
    const getUnitCollapseClass = (bar: GanttBarObject[]) => {
        if (bar.length == 0 ) return ""
        if (bar[0].ganttGroupId == undefined) return "" // 追加行のためクラス不要
        const ganttGroupId = bar[0].ganttGroupId
        const unitId = ganttGroupList.value.find(v => v.id === ganttGroupId)?.unit_id
        if (unitId) {
            return isOpenUnit(unitId) ? "" : "unit-closed"
        }
        return ""
    }

    const refreshMilestones = async () => {
        const {list} = await useMilestoneTable(currentFacilityId)
        milestones.value.length = 0
        milestones.value.push(<MileStone>{
            date: new Date(facility.value.shipment_due_date + " 00:00:00"),
            description: "出荷期日"
        })
        milestones.value.push(
            ...list.value.map(v => {
                return <MileStone>{date: new Date(v.date), description: v.description}
            })
        )
    }
    await refreshMilestones()

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
            let filteredTicketList = ticketList.value.filter(ticket => ticket.gantt_group_id === ganttGroup.id)
                .sort((a, b) => a.order < b.order ? -1 : 1)
            if (selectedDepartment.value != undefined) {
                filteredTicketList = filteredTicketList.filter(ticket => ticket.department_id == selectedDepartment.value)
            }
            if (selectedUser.value != undefined) {
                // 選択されたユーザーだけに絞り込む
                const targetTicketIds = ticketUserList.value.filter(ticketUser => ticketUser.user_id == selectedUser.value).map(v => v.ticket_id)
                filteredTicketList = filteredTicketList.filter(ticket => {
                    return targetTicketIds.includes(ticket.id!)
                })
            }

            // push処理を実行すべきかどうかを判断する条件
            // 1件以上あればpushする。部署が選択されていなければ空でもプッシュする。
            const shouldPush = (selectedDepartment.value === undefined && selectedUser.value === undefined) || filteredTicketList.length > 0;

            if (shouldPush) {
                ganttChartGroup.value.push(
                    {
                        ganttGroup: ganttGroup,
                        rows: filteredTicketList
                            .map(ticket => ticketToGanttRow(ticket, ticketUserList.value, ganttGroup)),
                        unitId: ganttGroup.unit_id,
                    }
                );
            }
        })
    }
    const hasFilter = computed(() => {
        return selectedDepartment.value != undefined || selectedUser.value
    })


    // computedだとドラッグ関連が上手くいかない為やはりオブジェクトとして双方向に同期をとるようにする。
    const bars = ref<GanttBarObject[][]>([])
    // ガントチャート描画用のオブジェクトの生成
    const refreshBars = async () => {
        const prodGanttBarObjects = await getProdTicketsIfSimulation()
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
            // まとめて表示機能を実装する. まとめられている場合はユニットの情報をすべて１つの要素に詰め込める。
            if (!isOpenUnit(unit.unitId)) {
                const result: GanttBarObject[] = []
                unit.rows.forEach(task => {
                    const prodBar = prodGanttBarObjects.find(v => v.ganttBarConfig.id == task.ticket?.id?.toString())
                    if (prodBar != undefined) {
                        prodBar.ganttBarConfig.id = "simulate_" + prodBar.ganttBarConfig.id
                        result.push(prodBar)
                    }
                    result.push(<GanttBarObject>{
                        beginDate: dayjs(task.ticket!.start_date!).format(DAYJS_FORMAT),
                        endDate: endOfDay(task.ticket!.end_date!),
                        ganttGroupId: unit.ganttGroup.id,
                        ganttBarConfig: {
                            dragLimitLeft: 0,
                            dragLimitRight: 0,
                            hasHandles: allowed('UPDATE_TICKET'),
                            immobile: !allowed('UPDATE_TICKET'), // TODO: なぜか判定が逆転している
                            id: task.ticket?.id!.toString(),
                            label: getProcessName(task.ticket?.process_id == null ? -1 : task.ticket?.process_id),
                            style: {backgroundColor: getProcessColor(task.ticket?.process_id)},
                            progress: task.ticket?.progress_percent,
                            progressColor: BAR_COMPLETE_COLOR
                        }
                    })
                })
                bars.value.push(result)
            } else {
                bars.value.push(...unit.rows.map(task => {
                    const result: GanttBarObject[] = []
                    const prodBar = prodGanttBarObjects.find(v => v.ganttBarConfig.id == task.ticket?.id?.toString())
                    if (prodBar != undefined) {
                        prodBar.ganttBarConfig.id = "simulate_" + prodBar.ganttBarConfig.id
                        result.push(prodBar)
                    }
                    result.push(<GanttBarObject>{
                        beginDate: dayjs(task.ticket!.start_date!).format(DAYJS_FORMAT),
                        endDate: endOfDay(task.ticket!.end_date!),
                        ganttGroupId: unit.ganttGroup.id,
                        ganttBarConfig: {
                            dragLimitLeft: 0,
                            dragLimitRight: 0,
                            hasHandles: allowed('UPDATE_TICKET'),
                            immobile: !allowed('UPDATE_TICKET'), // TODO: なぜか判定が逆転している
                            id: task.ticket?.id!.toString(),
                            label: getProcessName(task.ticket?.process_id == null ? -1 : task.ticket?.process_id),
                            style: {backgroundColor: getProcessColor(task.ticket?.process_id)},
                            progress: task.ticket?.progress_percent,
                            progressColor: BAR_COMPLETE_COLOR
                        }
                    })
                    return result
                }))
            }
            // ボタン用の空行を追加する
            if (!hasFilter.value && allowed('UPDATE_TICKET') && isOpenUnit(unit.ganttGroup.unit_id)) {
                bars.value.push([emptyRow])
            }
        })
        console.log(bars.value)
    }

    // シミュレーション中の場合本番のチケットを習得する
    const getProdTicketsIfSimulation = async () => {
        const result: GanttBarObject[] = []
        const isSimulateUser = getUserInfo().isSimulateUser
        if (isSimulateUser == true) {
            const {data: prodTickets} = await Api.getTickets(ganttGroupIds, ApiMode.prod)
            result.push(...
                prodTickets.list.map(v => {
                    return <GanttBarObject>{
                        beginDate: dayjs(v.start_date!).format(DAYJS_FORMAT),
                        endDate: endOfDay(v.end_date!),
                        ganttGroupId: v.gantt_group_id,
                        ganttBarConfig: {
                            dragLimitLeft: 0,
                            dragLimitRight: 0,
                            hasHandles: false,
                            immobile: true, // TODO: なぜか判定が逆転している
                            id: v.id!.toString(),
                            label: getProcessName(v.process_id == null ? -1 : v.process_id),
                            style: {backgroundColor: getProcessColor(v.process_id), opacity: 0.5},
                            progress: v.progress_percent,
                            progressColor: BAR_COMPLETE_COLOR
                        }
                    }
                })
            )
        }
        return result
    }

    refreshLocalGantt()
    await refreshBars()

    const isUpdateOrder = ref(false)

    const updateOrder = async (ganttRows: GanttRow[], index: number, direction: number) => {
        isUpdateOrder.value = true
        try {
            // ソート前にガントチャート上のIndexを検索しておく。 TODO: シミュレーション対応
            const barIndex = bars.value.findIndex(v => v[0].ganttBarConfig.id === ganttRows[index].bar.ganttBarConfig.id)
            const sorted = changeSort(ganttRows, index, direction)
            const newTickets: Ticket[] = []
            // 変更がった場合はガントチャートのオブジェクトと同期をとる
            if (sorted) {
                // bars.valueは全体から見たIndexを指定する必要があった。
                changeSort(bars.value, barIndex, direction)
                for (const v of ganttRows) {
                    const clone = Object.assign({}, v.ticket)
                    clone.order = ganttRows.indexOf(v)
                    newTickets.push(clone)
                }
                const tickets = newTickets.map(v => modifyTicketDateTimes(v))
                const {data} = await Api.postBulkUpdateTickets({tickets: tickets})
                await refreshTickets(data.tickets)
            }
        } finally {
            isUpdateOrder.value = false
        }
    }

    // 部署の変更。
    const getUserListByDepartmentId = (departmentId?: number, startDate?: string, endDate?: string) => {
        let filteredUsers = (departmentId == null) ? userList : userList.filter(v => v.department_id === departmentId);

        if (startDate && endDate) {
            filteredUsers = filterEmploymentDuration(filteredUsers, startDate, endDate)
        }

        return filteredUsers;
    }

    const filterEmploymentDuration = (userList: User[], startDate?: string | null, endDate?: string | null) => {
        if (startDate && endDate) {
            const start = new Date(startDate).getTime();
            const end = new Date(endDate).getTime();

            userList = userList.filter(user => {
                const userStart = new Date(user.employment_start_date.substring(0, 10)).getTime();
                const userEnd = user.employment_end_date ? new Date(user.employment_end_date.substring(0, 10)).getTime() : Infinity;
                // 指定された日付範囲が、利用者の雇用期間と被っている場合にのみ、その利用者を含めます。
                return (userStart <= end && userEnd >= start);
            });
        }
        return userList;
    }

    const deleteTicket = async (ticket: Ticket) => {
        // TODO: confirmを作る
        await Api.deleteTicketsId(ticket.id!)
        // チケット情報をリフレッシュする
        await ticketRefresh(ganttGroupIds)
        refreshLocalGantt()
        await refreshBars()
    }

    // DBへのストア及びローカルのガントに情報を反映する
    const updateTicket = async (ticket: Ticket) => {
        const reqTicket = Object.assign({}, ticket)
        modifyTicketDateTimes(reqTicket)
        try {
            console.log("########### UPDATE_TICKET START CALL API")
            const {data} = await Api.postTicketsId(ticket.id!, {ticket: reqTicket})
            console.log("########### UPDATE_TICKET UPDATE FINISH")
            return data.ticket
        } catch (e) {
            console.log(e)
        }
    }

    const modifyTicketDateTimes = (ticket: Ticket) => {
        if (ticket.start_date) {
            ticket.start_date = ticket.start_date + "T00:00:00.00000+09:00"
        } else {
            ticket.start_date = undefined
        }
        if (ticket.end_date) {
            ticket.end_date = ticket.end_date + "T00:00:00.00000+09:00"
        } else {
            ticket.end_date = undefined
        }
        if (ticket.limit_date != undefined) {
            ticket.limit_date = ticket.limit_date + "T00:00:00.00000+09:00"
        }
        return ticket
    }

    /**
     * 与えられたチケットに関連付けられたユーザーが雇用期間内にいるかどうかを検証します
     *
     * @param {Ticket} ticket - チケットの詳細が格納されたチケットオブジェクト
     * @returns {boolean} - チケットに関連付けられたユーザーが雇用期間内にいる場合はtrueを返し、そうでない場合はfalseを返します
     */
    const validateEmploymentTicketUser = async (ticket: Ticket) => {
        const ticketId = ticket.id!
        const ticketUserIds = ticketUserList.value.filter(v => v.ticket_id == ticketId).map(v => v.user_id)
        const ticketUsers = userList.filter(v => ticketUserIds.includes(v.id!))
        const afterUserList = filterEmploymentDuration(ticketUsers, ticket.start_date, ticket.end_date)
        const invalidUsers = ticketUsers.filter(v => !afterUserList.find(vv => vv.id == v.id))
        let result = true
        if (invalidUsers.length > 0) {
            const userNames = invalidUsers.map(v => v.lastName + v.firstName).join("\n")
            const swalResult = await Swal.fire({
                title: '在籍期間外になる担当者が存在します。',
                text: userNames,
                icon: 'warning',
                showCancelButton: true,
                confirmButtonColor: '#3085d6',
                cancelButtonColor: '#d33',
                confirmButtonText: '更新',
            })
            result = swalResult.isConfirmed
        }
        return {result, afterUserList}
    }

    // DBからのチケット更新をフロントに反映させる。
    const refreshTicket = async (ticket: Ticket) => {
        const target = ticketList.value.find(v => v.id === ticket.id)
        if (target == undefined) {
            console.error("チケットが存在しません。")
            return
        }
        target.days_after = ticket.days_after
        target.department_id = ticket.department_id
        target.end_date = ticket.end_date
        target.estimate = ticket.estimate
        target.limit_date = ticket.limit_date
        target.order = ticket.order
        target.process_id = ticket.process_id
        target.progress_percent = ticket.progress_percent
        target.start_date = ticket.start_date
        target.updated_at = ticket.updated_at
        target.number_of_worker = ticket.number_of_worker
        refreshLocalGantt()
        await refreshBars()
    }
    const refreshTickets = async (ticket: Ticket[]) => {
        ticket.forEach(ticket => {
            const target = ticketList.value.find(v => v.id === ticket.id)
            if (target == undefined) {
                console.error("チケットが存在しません。")
                return
            }
            target.days_after = ticket.days_after
            target.department_id = ticket.department_id
            target.end_date = ticket.end_date
            target.estimate = ticket.estimate
            target.limit_date = ticket.limit_date
            target.order = ticket.order
            target.process_id = ticket.process_id
            target.progress_percent = ticket.progress_percent
            target.start_date = ticket.start_date
            target.updated_at = ticket.updated_at
            target.number_of_worker = ticket.number_of_worker
        })
        refreshLocalGantt()
        await refreshBars()
    }


    /**
     * 忘れそうなのでメモ。
     * TicketMemoを更新してUpdatedAtを更新しただけだと競合が起きるので、APIから再度取り直す。
     * @param ticketId
     */
    const refreshTicketMemo = async (ticketId: number) => {
        const target = ticketList.value.find(v => v.id === ticketId)
        if (target == undefined) {
            console.error("チケットが存在しません。")
            return
        }
        const {data} = await Api.getTicketsId(ticketId)
        await refreshTicket(data.ticket!)
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
            updated_at: 0,
            number_of_worker: 1
        }
        const {data} = await Api.postTickets({ticket: newTicket})
        ticketList.value.push(data.ticket!)
        refreshLocalGantt()
        await refreshBars()
    }
    // DBへのストア及びローカルのガントに情報を反映する
    const ticketUserUpdate = async (ticket: Ticket, userIds: number[]) => {
        // ticketUsersのcreatedAtを検索する
        const ticketUser = ticketUserList.value.find(v => v.ticket_id == ticket.id)
        let createdAt = undefined
        if (ticketUser != undefined) {
            createdAt = ticketUser.created_at
        }
        try {
            const {data} = await Api.postTicketUsers({ticketId: ticket.id!, userIds: userIds, createdAt: createdAt})
            return data.ticketUsers
        } catch (e) {
            console.warn(e)
        }
    }

    const adjustBar = async (bar: GanttBarObject) => {
        // id をもとに VueJs 側のデータも更新する
        const targetGanttGroup = ganttChartGroup.value.find(v => v.ganttGroup.id === bar.ganttGroupId)
        if (!targetGanttGroup) {
            console.error(`Unit ID: ${bar.ganttGroupId} が現在のガントに存在しません。`, ganttChartGroup)
        } else {
            const targetTicket = targetGanttGroup.rows.find(v => v.ticket?.id!.toString() === bar.ganttBarConfig.id)
            if (!targetTicket) {
                console.error(`ID: ${bar.ganttBarConfig.id} が現在のガントに存在しません。`, ganttChartGroup)
            } else {
                const clone = Object.assign({}, targetTicket.ticket)
                clone.start_date = ganttDateToYMDDate(bar.beginDate)
                clone.end_date = ganttDateToYMDDate(bar.endDate)
                const {result, afterUserList} = await validateEmploymentTicketUser(clone)
                if (!result) {
                    bar.beginDate = YMDDateToGanttStartDate(targetTicket.ticket!.start_date!)
                    bar.endDate = YMDDateToGanttEndDate(targetTicket.ticket!.end_date!)
                    return
                } else {
                    await mutation.setTicketUser(clone, afterUserList.map(v => v.id!))
                }
                // NOTE:少し気持ち悪いが setTicketUserでTicketの更新も行うのでコメントアウト。
                // 担当の変更が起きるticketの人数を同時に更新する必要があるため。
                // const newTicket = await updateTicket(clone)
                // await refreshTicket(newTicket!)
            }
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
    const setScheduleByPersonDay = async (rows: GanttRow[]) => {
        const holidays = toValue(getHolidays)
        const operationSettings = toValue(getOperationList)
        const newTickets: Ticket[] = []
        let prevEndDate: Dayjs | null = null

        for (let index = 0; index < rows.length; index++) {
            const row = rows[index]
            const clone = Object.assign({}, row.ticket)

            let startDate: Dayjs
            // 先頭行の場合
            if (index === 0) {
                // 開始日・工数・工程・人数が未設定の場合はスキップする
                if (!clone?.start_date || !clone?.estimate ||
                    !clone?.number_of_worker || !clone?.process_id
                ) {
                    console.log("###### SKIP due to first row")
                    continue;
                }
                startDate = dayjs(clone!.start_date)
            } else {
                // 先頭行以降は (開始日 or days_after)・工数・工程・人数が未設定の場合はスキップする
                if (!clone?.estimate || !clone?.number_of_worker || !clone?.process_id) {
                    console.log("###### SKIP due to after first row pt1")
                    continue;
                } else {
                    // 開始日も日後も設定がなければスキップ
                    if (!clone?.start_date && !clone?.days_after) {
                        console.log("###### SKIP due to after first row pt2")
                        continue;
                    } else {
                        // 日後が0でなければ優先
                        if (clone?.days_after != null) {
                            startDate = addBusinessDays(prevEndDate!, clone.days_after, holidays)
                        } else {
                            // 開始日をそのまま利用する
                            startDate = dayjs(clone!.start_date)
                        }
                    }
                }
            }
            console.log("########### 工数重視の実行", startDate.format(DAYJS_FORMAT))
            // 工程から総稼働予定時間を取得する
            const scheduledOperatingHours = getScheduledOperatingHours(operationSettings, row) * clone.number_of_worker
            // 必要日数。少数が出たら最小の整数にする。例：二人で5人日の場合 3日必要 5/2 = 2.5 つまり少数が出たら整数に足す
            const numberOfDaysRequired = Math.ceil(clone.estimate! / scheduledOperatingHours)
            // 開始日を設定する
            clone.start_date = ganttDateToYMDDate(startDate.format(DAYJS_FORMAT))
            // 終了日を決定する、祝日が含まれている場合終了日をずらす。
            clone.end_date = ganttDateToYMDDate(
                addBusinessDays(startDate, numberOfDaysRequired, holidays, true)
                    .add(-1, 'minute').format(DAYJS_FORMAT)
            )
            const newTicket = await updateTicket(clone)
            newTickets.push(newTicket!)
            // 前回の終了日を設定する
            prevEndDate = dayjs(newTicket!.end_date)
            console.log("########### 工数重視の完了", startDate.format(DAYJS_FORMAT))
        }
        await refreshTickets(newTickets)
    }

    /**
     * リスケ（日付）重視
     * 現在設定されている日付に日後を適応させる。
     *
     * @param rows
     */
    const setScheduleByFromTo = async (rows: GanttRow[]) => {
        const holidays = toValue(getHolidays)
        let prevEndDate: Dayjs | null = null
        const newTickets: Ticket[] = []
        for (let index = 0; index < rows.length; index++) {
            const row = rows[index]
            const clone = Object.assign({}, row.ticket)
            if (clone && clone?.start_date && clone?.end_date) {
                // 日後の処理
                if (clone.days_after != null && prevEndDate != null) {
                    // 現在設定されている営業日を計算する
                    let dayjsStartDate = dayjs(clone.start_date)
                    let dayjsEndDate = dayjs(clone.end_date)
                    const numberOfRequiredBusinessDays = getNumberOfBusinessDays(dayjsStartDate, dayjsEndDate, holidays)
                    // 開始日を日後分ずらした日付に設定する
                    dayjsStartDate = addBusinessDays(prevEndDate, clone.days_after, holidays)
                    // 終了日を設定されている営業日分確保する
                    dayjsEndDate = getEndDateByRequiredBusinessDay(dayjsStartDate, numberOfRequiredBusinessDays, holidays)
                    // チケットに反映させる
                    clone.start_date = ganttDateToYMDDate(dayjsStartDate.format(DAYJS_FORMAT))
                    clone.end_date = ganttDateToYMDDate(dayjsEndDate.format(DAYJS_FORMAT))
                }
                const newTicket = await updateTicket(clone)
                newTickets.push(newTicket!)
            }
            prevEndDate = dayjs(clone.end_date)
        }
        await refreshTickets(newTickets)
    }

    const mutation = {
        setProcessId: async (processId: string, ticket?: Ticket) => {
            console.log("############# processId", processId)
            const clone = Object.assign({}, ticket)
            clone.process_id = Number(processId)
            const newTicket = await updateTicket(clone)
            await refreshTicket(newTicket!)
            await getScheduleAlert()

        }, setDepartmentId: async (departmentId: string, ticket?: Ticket) => {
            const clone = Object.assign({}, ticket)
            if (departmentId == "") {
                clone.department_id = undefined
            } else {
                clone.department_id = Number(departmentId)
            }
            clone.number_of_worker = 1
            const newTicket = await updateTicket(clone)

            // 担当者をすべて外す。
            const ticketUsers = ticketUserList.value.filter(v => v.ticket_id === ticket?.id)
            if (ticketUsers.length > 0) {
                ticketUsers.length = 0
            }
            // 部署替えの場合既存のユーザー一覧を空に更新する
            const newTicketUsers = await ticketUserUpdate(clone, [])
            const newTicketUserList = ticketUserList.value.filter(v => v.ticket_id !== ticket!.id!)
            newTicketUserList.push(...newTicketUsers!)
            ticketUserList.value.length = 0
            ticketUserList.value.push(...newTicketUserList)
            await refreshTicket(newTicket!)

        }, setNumberOfWorker: async (numberOfWorker: number, ticket?: Ticket) => {
            const clone = Object.assign({}, ticket)
            clone.number_of_worker = numberOfWorker
            const newTicket = await updateTicket(clone)
            await refreshTicket(newTicket!)
            await getScheduleAlert()
        }, setEstimate: async (estimate: number, ticket?: Ticket) => {
            const clone = Object.assign({}, ticket)
            clone.estimate = estimate
            const newTicket = await updateTicket(clone)
            await refreshTicket(newTicket!)
            await getScheduleAlert()
        }, setDaysAfter: async (daysAfter: number, ticket?: Ticket) => {
            const clone = Object.assign({}, ticket)
            clone.days_after = daysAfter
            const newTicket = await updateTicket(clone)
            await refreshTicket(newTicket!)
            await getScheduleAlert()
        }, setStartDate: async (startDate: string, ticket?: Ticket) => {
            const clone = Object.assign({}, ticket)
            clone.start_date = startDate
            const {result, afterUserList} = await validateEmploymentTicketUser(clone)
            if (!result) {
                await refreshTicket(ticket!)
                return
            } else {
                await mutation.setTicketUser(clone, afterUserList.map(v => v.id!))
            }
            // const newTicket = await updateTicket(clone)
            // await refreshTicket(newTicket!)
            // await getScheduleAlert()
        }, setEndDate: async (endDate: string, ticket?: Ticket) => {
            const clone = Object.assign({}, ticket)
            clone.end_date = endDate
            const {result, afterUserList} = await validateEmploymentTicketUser(clone)
            if (!result) {
                await refreshTicket(ticket!)
                return
            } else {
                await mutation.setTicketUser(clone, afterUserList.map(v => v.id!))
            }
            // const newTicket = await updateTicket(clone)
            // await refreshTicket(newTicket!)
            // await getScheduleAlert()
        }, setProgressPercent: async (progressPercent: number, ticket?: Ticket) => {
            const clone = Object.assign({}, ticket)
            clone.progress_percent = progressPercent
            const newTicket = await updateTicket(clone)
            await refreshTicket(newTicket!)
            await getScheduleAlert()
        }, setTicketUser: async (ticket: Ticket, value: number[]) => {
            const clone = Object.assign({}, ticket)
            clone.number_of_worker = Number(value.length === 0 ? 1 : value.length)
            const newTicket = await updateTicket(clone)

            const newTicketUsers = await ticketUserUpdate(clone, value)
            const newTicketUserList = ticketUserList.value.filter(v => v.ticket_id !== ticket!.id!)
            newTicketUserList.push(...newTicketUsers!)
            ticketUserList.value.length = 0
            ticketUserList.value.push(...newTicketUserList)
            await refreshTicket(newTicket!)
            await getScheduleAlert()
        }
    }

    const getUnitIdByTicketId = (ticketId: number) => {
        const ticket = ticketList.value.find(v => v.id == ticketId)!
        const ganttGroup = ganttGroupList.value.find(v => v.id === ticket.gantt_group_id)!
        return ganttGroup.unit_id
    }

    return {
        bars,
        chartEnd,
        chartStart,
        facility,
        getHolidaysForGantt,
        ganttChartGroup,
        getDepartmentName,
        getGanttChartWidth,
        getUnitName,
        getOperationList,
        getHolidays,
        getTickets,
        ticketUserList,
        addNewTicket,
        adjustBar,
        deleteTicket,
        getUserListByDepartmentId,
        setScheduleByFromTo,
        setScheduleByPersonDay,
        ticketUserUpdate,
        updateOrder,
        updateTicket,
        refreshTicketMemo,
        hasFilter,
        milestones,
        mutation,
        getUnitIdByTicketId,
        isUpdateOrder,
        isOpenUnit,
        toggleUnitOpen,
        getUnitCollapseClass,
        isAllOpenUnit,
        toggleAllUnitOpen,
    }
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
            beginDate: dayjs(ticket.start_date!).format(DAYJS_FORMAT),
            endDate: dayjs(ticket.end_date!).format(DAYJS_FORMAT),
            ganttBarConfig: {
                hasHandles: true,
                id: ticket.id!.toString(), // TODO: まあIDが数字でもいっか
                immobile: false,
                label: "", // TODO: コメントをラベルにする
            },
        },
        ticket: ticket,
        ticketUsers: ticketUserList.filter(v => v.ticket_id === ticket.id),
        ganttGroup: ganttGroup
    }
    return result
}


