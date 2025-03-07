package constants

const (
	RoleAdmin   = "admin"
	RoleManager = "manager"
	RoleWorker  = "worker"
	RoleViewer  = "viewer"
	RoleGuest   = "guest"
)

var RoleDisplayNames = map[string]string{
	"管理者":    "admin",
	"マネージャー": "manager",
	"作業者":    "worker",
	"閲覧者":    "viewer",
	"ゲスト":    "guest",
}

const (
	GuestID = int32(-1) // GUESTログインで使用するユーザーID
)
