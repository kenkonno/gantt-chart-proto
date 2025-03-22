package upload_users_csv_file

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/constants"
	"github.com/kenkonno/gantt-chart-proto/backend/api/middleware"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/api/utils"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"github.com/kenkonno/gantt-chart-proto/backend/repository/common"
	"github.com/samber/lo"
	"net/http"
	"strconv"
	"strings"
	"time"
)

/*
*
初のCSVアップロードなのでメモ
Requestはフレームワーク（自分の）ルール上必要なので作っているが意味ない。
実態はopenapiに定義してるoctetstreamのstring/binary
*/
func PostUploadUsersCsvFileInvoke(c *gin.Context) (openapi_models.PostUploadUsersCsvFileResponse, error) {

	userRep := repository.NewUserRepository(middleware.GetRepositoryMode(c)...)

	r := csv.NewReader(c.Request.Body)
	rows, err := r.ReadAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, "CSVの読み込みに失敗しました。")
		fmt.Println(err)
		return openapi_models.PostUploadUsersCsvFileResponse{}, err
	}

	column := map[string]int{
		"DepartmentId":        0,
		"LastName":            1,
		"FirstName":           2,
		"Role":                3,
		"Email":               4,
		"Password":            5,
		"EmploymentStartDate": 6,
		"EmploymentEndDate":   7,
	}

	var newUsers []db.User
	for i, row := range rows {
		// 先頭はヘッダー
		if i == 0 {
			continue
		}

		for ii, col := range row {
			row[ii] = strings.TrimRight(strings.TrimLeft(col, " "), " ")
		}

		departmentId := row[column["DepartmentId"]]
		lastName := row[column["LastName"]]
		firstName := row[column["FirstName"]]
		role := row[column["Role"]]
		email := row[column["Email"]]
		password := row[column["Password"]]
		employmentStartDate := row[column["EmploymentStartDate"]]
		employmentEndDate := row[column["EmploymentEndDate"]]

		user, err := validateRow(departmentId, lastName, firstName, role, email, password, employmentStartDate, employmentEndDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return openapi_models.PostUploadUsersCsvFileResponse{}, err
		}
		newUsers = append(newUsers, user)
	}

	oldUsers := userRep.FindAll()

	allUsers := append(newUsers, oldUsers...)

	for _, v := range allUsers {
		if len(lo.Filter(allUsers, func(item db.User, index int) bool {
			return item.Email == v.Email
		})) > 1 {
			c.JSON(http.StatusBadRequest, fmt.Sprintf("メールアドレスが重複しています。 %s", v.Email))
			return openapi_models.PostUploadUsersCsvFileResponse{}, errors.New("duplicated email address")
		}
	}

	for _, v := range newUsers {
		userRep.Upsert(v)
	}

	return openapi_models.PostUploadUsersCsvFileResponse{}, nil

}

func validateRow(departmentId string, lastName string, firstName string, role string, email string, password string, employmentStartDate string, employmentEndDate string) (db.User, error) {

	// 本当はここにrepository生成するのはNGだけどフレームワークの作りが甘かったので妥協
	departmentRep := common.NewDepartmentRepository()

	id, err := strconv.Atoi(departmentId)
	if err != nil {
		return db.User{}, err
	}
	department := departmentRep.Find(int32(id))
	if department.Id == nil {
		return db.User{}, errors.New(fmt.Sprintf("存在しない部署IDです。%d", id))
	}

	if lastName == "" {
		return db.User{}, errors.New(fmt.Sprintf("名は必須です"))
	}
	if firstName == "" {
		return db.User{}, errors.New(fmt.Sprintf("性は必須です"))
	}
	if !utils.ValidatePassword(password) {
		return db.User{}, errors.New(fmt.Sprintf(constants.E0001))
	}
	hashedPassword, err := utils.CryptPassword(password)
	if err != nil {
		return db.User{}, errors.New(fmt.Sprintf("パスワードの暗号化に失敗しました。"))
	}

	if utils.GetDisplayNameRole(role) == "" {
		return db.User{}, errors.New(fmt.Sprintf("存在しないRoleです。%s", role))
	}

	if utils.ValidateEmail(email) != nil {
		return db.User{}, errors.New(fmt.Sprintf("不正なメールアドレス形式です。%s", email))
	}

	esd := utils.GetTimeByYMDString(employmentStartDate)
	if esd == nil {
		return db.User{}, errors.New(fmt.Sprintf("在籍期間(開始)の形式が正しくありません。%s", employmentStartDate))
	}

	var eed *time.Time
	if employmentEndDate != "" {
		eed = utils.GetTimeByYMDString(employmentEndDate)
		if eed == nil {
			return db.User{}, errors.New(fmt.Sprintf("在籍期間(終了)の形式が正しくありません。%s", employmentEndDate))
		}
	}

	return db.User{
		Id:                  nil,
		DepartmentId:        *department.Id,
		LimitOfOperation:    8,
		LastName:            lastName,
		FirstName:           firstName,
		Password:            hashedPassword,
		Email:               strings.ToLower(email),
		Role:                role,
		PasswordReset:       false,
		EmploymentStartDate: *esd,
		EmploymentEndDate:   eed,
		CreatedAt:           time.Time{},
		UpdatedAt:           0,
	}, nil
}
