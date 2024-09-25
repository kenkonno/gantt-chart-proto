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
    Viewer: "viewer",
    Guest: "guest",
}
export const RoleTypeMap = {
    admin: "管理者",
    manager: "マネージャー",
    worker: "作業者",
    viewer: "閲覧者",
}

