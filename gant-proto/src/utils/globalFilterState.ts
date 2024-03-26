/**
 * Global State
 * facilityに紐づかないものは直接
 * facilityに紐づくものは連想はれいつとして持つ。
 */
import {AggregationAxis, DisplayType, GanttFacilityHeader} from "@/composable/ganttFacilityMenu";
import {Header} from "@/composable/ganttAllMenu";
import {PileUpFilter} from "@/composable/pileUps";
import {FacilityType} from "@/const/common";

const LOCAL_STORAGE_KEY = "koteikanri"
const VERSION = 1.1 // フィルタの項目が変わったときに変える

/**
 * WebStorageを用いて各種フィルタ系の値を保持するようにする。
 *
 * 対象
 *
 * ・案件状況
 *   確定、未確定のチェックボックス
 *
 * ・ガントチャートフィルタ
 *   部署フィルタ
 *   担当者フィルタ
 *
 * ・案件ビューのヘッダー
 *
 * ・全体ビューのヘッダー
 *
 * ・PileUpsの部署を開くかどうか
 *
 *
 * 設計むずいなー
 * 基本コンポーネントの destory で保存して 作るときに初期値をこれで持ってくるって形がシンプルではある。
 */

// 初期値
type GlobalFilterState = {
    orderStatus: string[],
    ganttFacilityMenu: GanttFacilityHeader[],
    ganttAllMenu: Header[],
    pileUpsFilter: PileUpFilter[],
    viewType: DisplayType,
    aggregationAxis: AggregationAxis,
    version: number
}


// 初期値
const state: GlobalFilterState = {
    orderStatus: [FacilityType.Ordered],
    ganttFacilityMenu: [
        {name: "ユニット", visible: true},
        {name: "工程", visible: true},
        {name: "部署", visible: true},
        {name: "担当者", visible: true},
        {name: "人数", visible: false},
        {name: "工数(h)", visible: true},
        {name: "日後", visible: false},
        {name: "開始日", visible: false},
        {name: "終了日", visible: false},
        {name: "進捗", visible: true},
        {name: "操作", visible: false},
    ],
    ganttAllMenu: [
        {name: "案件名", visible: true},
        {name: "担当者", visible: false},
        {name: "開始日", visible: true},
        {name: "終了日", visible: true},
        {name: "工数(h)", visible: true},
        {name: "進捗", visible: true},
    ],
    pileUpsFilter: [],
    viewType: "day",
    aggregationAxis: "facility",
    version: VERSION

}

export const initStateValue = async () => {
    // ローカルストレージから取得
    const savedState = localStorage.getItem(LOCAL_STORAGE_KEY);
    if (savedState) {
        const parsedState = JSON.parse(savedState);
        // バージョンが異なる場合は初期化する
        if (parsedState.version != VERSION) {
            localStorage.setItem(LOCAL_STORAGE_KEY, JSON.stringify(state))
            return
        }
        state.orderStatus = parsedState.orderStatus;
        state.ganttFacilityMenu = getFacilityMenu(parsedState.ganttFacilityMenu);
        state.ganttAllMenu = parsedState.ganttAllMenu;
        state.pileUpsFilter = parsedState.pileUpsFilter;
        state.viewType = parsedState.viewType;
        state.aggregationAxis = parsedState.aggregationAxis;
        state.version = parsedState.version;
    } else {
        localStorage.setItem(LOCAL_STORAGE_KEY, JSON.stringify(state))
    }
}

/**
 * デフォルト値をもとにsavedFacilityMenuの値で鵜輪がいたものを返却する。
 * @param savedFacilityMenu
 */
const getFacilityMenu = (savedFacilityMenu: any): GanttFacilityHeader[] => {
    return state.ganttFacilityMenu.map((v, i) => {
        const savedValue = savedFacilityMenu.find((vv: any) => vv.name === v.name)
        if (savedValue != undefined) {
            v.visible = savedValue.visible
        }
        return {name: v.name, visible: v.visible}
    })
}

export const globalFilterGetter = {
    getOrderStatus: () => state.orderStatus,
    getGanttFacilityMenu: () => state.ganttFacilityMenu,
    getGanttAllMenu: () => state.ganttAllMenu,
    getPileUpsFilter: () => state.pileUpsFilter,
    getViewType: () => state.viewType,
    getAggregationAxis: () => state.aggregationAxis,
}

type StateKey = 'orderStatus' | 'ganttFacilityMenu' | 'ganttAllMenu' | 'pileUpsFilter' | 'viewType' | 'aggregationAxis';

function updateState(key: StateKey, value: any) {
    const storage = localStorage.getItem(LOCAL_STORAGE_KEY)
    if (storage == null) return
    const savedState = JSON.parse(storage) as GlobalFilterState;
    // NOTE: 型推論が完璧ではないためこのような分岐を利用しないと型エラーとなる。
    if ( key == "aggregationAxis") {
        savedState[key] = value;
        state[key] = value;
    } else {
        savedState[key] = value;
        state[key] = value;

    }
    localStorage.setItem(LOCAL_STORAGE_KEY, JSON.stringify(savedState))
}

export const globalFilterMutation = {
    updateOrderStatus: (orderStatus: string[]) => {
        updateState('orderStatus', orderStatus);
    },
    updateGanttFacilityMenu: (ganttFacilityMenu: GanttFacilityHeader[]) => {
        updateState('ganttFacilityMenu', ganttFacilityMenu);
    },
    updateGanttAllMenu: (header: Header[]) => {
        updateState('ganttAllMenu', header);
    },
    updatePileUpsFilter: (pileUpsFilter: PileUpFilter[]) => {
        updateState('pileUpsFilter', pileUpsFilter);
    },
    updateViewTypeFilter: (viewType: DisplayType) => {
        updateState('viewType', viewType);
    },
    updateAggregationAxisFilter: (aggregationAxis: AggregationAxis) => {
        updateState('aggregationAxis', aggregationAxis);
    },
}