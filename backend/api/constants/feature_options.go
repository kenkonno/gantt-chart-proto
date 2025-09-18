package constants

// Feature は、アプリケーションの各機能を示す識別子です。
type Feature string

// Equals は文字列との比較を行います
func (f Feature) Equals(s string) bool {
	return string(f) == s
}

// アプリケーションで利用可能な機能フラグの定数
const (
	// --- スケジュール管理 ---

	// ScheduleSimulation は「シミュレーション機能」を示します。
	ScheduleSimulation Feature = "ScheduleSimulation" // TODO: 未実装
	// UnitExpandCollapse は「ユニット開く、縮小」機能を示します。
	UnitExpandCollapse Feature = "UnitExpandCollapse" // TODO: 未実装
	// UnitCopy は「ユニットコピー」機能を示します。
	UnitCopy Feature = "UnitCopy" // TODO: 未実装
	// MultiSelectFilter は「部署、担当者フィルター複数選択」機能を示します。
	MultiSelectFilter Feature = "MultiSelectFilter" // TODO: 未実装
	// ProjectListFreeText は「案件一覧での自由入力欄」機能を示します。
	ProjectListFreeText Feature = "ProjectListFreeText" // TODO: 未実装
	// ProjectListNameSort は「案件一覧での案件名でのソート」機能を示します。
	ProjectListNameSort Feature = "ProjectListNameSort"

	// --- 進捗管理 ---

	// ProgressInput は「進捗入力」機能を示します。
	ProgressInput Feature = "ProgressInput" // TODO: 未実装
	// DelayNotification は「遅延通知」機能を示します。
	DelayNotification Feature = "DelayNotification" // TODO: 未実装

	// --- 負荷管理 ---

	// ResourceStackingView は「山積み表示」機能を示します。
	ResourceStackingView Feature = "ResourceStackingView" // TODO: 未実装
	// ResourceStackingGraph は「山積みグラフ」機能を示します。
	ResourceStackingGraph Feature = "ResourceStackingGraph" // TODO: 現状リリースが面倒なのでオプション対応をする
	// WorkloadWeighting は「負荷の重みづけ機能」を示します。
	WorkloadWeighting Feature = "WorkloadWeighting"
)
