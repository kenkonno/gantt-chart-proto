export type Emit = (event: (string), ...args: any[]) => void

export const DEFAULT_PROCESS_COLOR = "rgb(66, 165, 246)"

export const FacilityStatus = {
    Enabled: "Enabled",
    Finished: "Finished",
    Disabled: "Disabled",
}
export const FacilityStatusMap = {
    Enabled: "有効",
    Finished: "完了",
    Disabled: "無効",
}
export const FacilityType = {
    Ordered: "Ordered",
    Prepared: "Prepared",
}
export const FacilityTypeMap = {
    Ordered: "確定",
    Prepared: "未確定",
}

export const RoleType = {
    Admin: "admin",
    Manager: "manager",
    Worker: "worker",
    WorkerWithPileUps: "worker_with_pile_ups",
    Viewer: "viewer",
    Guest: "guest",
}
export const RoleTypeMap = {
    admin: "管理者",
    manager: "マネージャー",
    worker: "作業者",
    worker_with_pile_ups: "作業者(山積み閲覧可)",
    viewer: "閲覧者",
}

export const ApiMode = {
    prod: "prod",
}

export const FacilityWorkScheduleType = {
    Holiday: "Holiday",
    WorkingDay: "WorkingDay",
}
export const FacilityWorkScheduleTypeMap = {
    Holiday: "休日",
    WorkingDay: "稼働日",
}

// 機能フラグの型定義
export type FeatureOption = string

// アプリケーションで利用可能な機能フラグの定数
export const FeatureOption = {
    // --- スケジュール管理 ---

    // シミュレーション機能
    ScheduleSimulation: "ScheduleSimulation" as FeatureOption,
    // ユニット開く、縮小機能
    UnitExpandCollapse: "UnitExpandCollapse" as FeatureOption,
    // ユニットコピー機能
    UnitCopy: "UnitCopy" as FeatureOption,
    // 部署、担当者フィルター複数選択機能
    MultiSelectFilter: "MultiSelectFilter" as FeatureOption,
    // 案件一覧での自由入力欄機能
    ProjectListFreeText: "ProjectListFreeText" as FeatureOption,
    // 案件一覧での案件名でのソート機能
    ProjectListNameSort: "ProjectListNameSort" as FeatureOption,

    // --- 進捗管理 ---

    // 進捗入力機能
    ProgressInput: "ProgressInput" as FeatureOption,
    // 遅延通知機能
    DelayNotification: "DelayNotification" as FeatureOption,

    // --- 負荷管理 ---

    // 山積み表示機能
    ResourceStackingView: "ResourceStackingView" as FeatureOption,
    // 山積みグラフ機能
    ResourceStackingGraph: "ResourceStackingGraph" as FeatureOption,
    // 負荷の重みづけ機能
    WorkloadWeighting: "WorkloadWeighting" as FeatureOption,
} as const