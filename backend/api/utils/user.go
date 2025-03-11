package utils

import (
	"github.com/kenkonno/gantt-chart-proto/backend/api/constants"
	"golang.org/x/crypto/bcrypt"
	"net/mail"
	"regexp"
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
	var result time.Time
	var err error
	if result, err = time.Parse("2006-01-02", ymd); err != nil {
		return nil
	}
	return &result
}

func ValidateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	return err
}
