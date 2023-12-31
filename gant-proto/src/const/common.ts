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
    Ordered: "受注済み",
    Prepared: "非受注",
}
