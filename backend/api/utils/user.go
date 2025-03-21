package utils

import (
	"github.com/kenkonno/gantt-chart-proto/backend/api/constants"
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
