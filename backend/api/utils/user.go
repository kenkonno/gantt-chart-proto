package utils

import (
	"fmt"
	"github.com/kenkonno/gantt-chart-proto/backend/api/constants"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"golang.org/x/crypto/bcrypt"
	"net/mail"
	"regexp"
	"strings"
	"time"
)

func ValidatePassword(password string) bool {
	var (
		containsMin     = regexp.MustCompile(`[a-z]`).MatchString
		containsMax     = regexp.MustCompile(`[A-Z]`).MatchString
		containsNum     = regexp.MustCompile(`[0-9]`).MatchString
		containsSpecial = regexp.MustCompile(`[!@#\$%^&*()]`).MatchString
		lengthValid     = regexp.MustCompile(`.{8,}`).MatchString
	)

	return containsMin(password) && containsMax(password) && containsNum(password) && containsSpecial(password) && lengthValid(password)
}

func CryptPassword(password string) (string, error) {
	// パスワードをハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil

}

// GetDisplayNameRole Roleのから名称を返します。存在しない場合は空文字を返します。
func GetDisplayNameRole(value string) string {
	v, ok := constants.RoleDisplayNames[value]
	fmt.Println(v, value)
	if ok {
		return v
	}
	return ""
}

// GetTimeByYMDString YYYY-MM-DD の形式から time.Timeを返却する。パースできなければnilを返す。
func GetTimeByYMDString(ymd string) *time.Time {
	// スラッシュをハイフンに置き換え
	normalizedYmd := strings.ReplaceAll(ymd, "/", "-")

	// ハイフンで分割
	parts := strings.Split(normalizedYmd, "-")
	if len(parts) != 3 {
		return nil
	}

	// 月と日を0埋めして2桁にする
	if len(parts[1]) == 1 {
		parts[1] = "0" + parts[1]
	}
	if len(parts[2]) == 1 {
		parts[2] = "0" + parts[2]
	}

	// 正規化された日付文字列を作成
	normalizedYmd = parts[0] + "-" + parts[1] + "-" + parts[2]

	// パース
	var result time.Time
	var err error
	if result, err = time.Parse("2006-01-02", normalizedYmd); err != nil {
		return nil
	}
	return &result

}

func ValidateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	return err
}

func DetectWorkOutsideEmploymentPeriods(tickets []db.Ticket, EmploymentStartDate time.Time, EmploymentEndDate *time.Time) []db.Ticket {
	var outsideTickets []db.Ticket
	for _, ticket := range tickets {
		// チケットの開始日または終了日がnilの場合はスキップ（対象外）
		if ticket.StartDate == nil || ticket.EndDate == nil {
			continue
		}

		// 雇用終了日がnilの場合（まだ退職していない場合）
		if EmploymentEndDate == nil {
			// チケットの終了日が雇用開始日より前なら被らない
			if ticket.EndDate.Before(EmploymentStartDate) {
				outsideTickets = append(outsideTickets, ticket)
			}
		} else {
			// 雇用終了日が設定されている場合（退職した場合）
			// チケットの終了日が雇用開始日より前 または チケットの開始日が雇用終了日より後なら被らない
			if ticket.EndDate.Before(EmploymentStartDate) || ticket.StartDate.After(*EmploymentEndDate) {
				outsideTickets = append(outsideTickets, ticket)
			}
		}
	}
	return outsideTickets
}
