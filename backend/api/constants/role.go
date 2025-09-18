package constants

const (
	RoleAdmin             = "admin"
	RoleManager           = "manager"
	RoleWorker            = "worker"
	RoleWorkerWithPileUps = "worker_with_pile_ups"
	RoleViewer            = "viewer"
	RoleGuest             = "guest"
)

var RoleDisplayNames = map[string]string{
	"管理者":               "admin",
	"マネージャー":         "manager",
	"作業者":               "worker",
	"作業者(山積み閲覧可)": "worker_with_pile_ups",
	"閲覧者":               "viewer",
	"ゲスト":               "guest",
}

const (
	GuestID = int32(-1) // GUESTログインで使用するユーザーID
)
