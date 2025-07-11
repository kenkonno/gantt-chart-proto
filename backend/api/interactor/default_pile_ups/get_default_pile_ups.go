package default_pile_ups

import (
	"github.com/gin-gonic/gin"
	"github.com/kenkonno/gantt-chart-proto/backend/api/constants"
	"github.com/kenkonno/gantt-chart-proto/backend/api/middleware"
	"github.com/kenkonno/gantt-chart-proto/backend/api/openapi_models"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"github.com/samber/lo"
	"golang.org/x/exp/slices"
	"math"
	"strconv"
	"time"
)

const errorStyle = "rgb(255 89 89)"
const grayOutStyle = "rgb(243 238 226)"

func GetDefaultPileUpsInvoke(c *gin.Context) (openapi_models.GetDefaultPileUpsResponse, error) {
	pileUpsRep := repository.NewPileUpsRepository(middleware.GetRepositoryMode(c)...)
	facilityRep := repository.NewFacilityRepository(middleware.GetRepositoryMode(c)...)
	departmentRep := repository.NewDepartmentRepository(middleware.GetRepositoryMode(c)...)
	userRep := repository.NewUserRepository(middleware.GetRepositoryMode(c)...)

	excludeFacilityId, err := strconv.Atoi(c.Query("facilityId"))
	if err != nil {
		panic(err)
	}
	isALlMode := c.Query("isAllMode") == "true"
	qFacilityTypes := c.QueryArray("facilityTypes")
	var facilityTypes []string
	if slices.Contains(qFacilityTypes, constants.FacilityTypeOrdered) {
		facilityTypes = append(facilityTypes, constants.FacilityTypeOrdered)
	}
	if slices.Contains(qFacilityTypes, constants.FacilityTypePrepared) {
		facilityTypes = append(facilityTypes, constants.FacilityTypePrepared)
	}

	departments := departmentRep.FindAll()
	facilities := facilityRep.FindAll(facilityTypes, []string{constants.FacilityStatusEnabled})
	facilityMap := lo.Associate(facilities, func(item db.Facility) (int32, db.Facility) {
		return *item.Id, item
	})

	users := userRep.FindAll()
	departmentUserMap := lo.GroupBy(users, func(item db.User) int32 {
		return item.DepartmentId
	})
	userMap := lo.Associate(users, func(item db.User) (int32, db.User) {
		return *item.Id, item
	})

	//////////////////////////////////////////////
	// 初期化をする
	//////////////////////////////////////////////
	var defaultPileUps []openapi_models.DefaultPileUp
	defaultPileUps, globalStartDate := getDefaultPileUps(departments, defaultPileUps, departmentUserMap, facilities)

	// クエリーで全ての設備の積み上げを取得する
	if isALlMode {
		excludeFacilityId = -1 // 全件取得モード（全体ビュー）の場合は除外するFacilityIdを存在しないものにする。
	}
	allFacilityPileUps := pileUpsRep.GetDefaultPileUps(int32(excludeFacilityId), facilityTypes)
	// validIndexに対する稼動可能なユーザーの一覧
	validUserIndexes := pileUpsRep.GetValidIndexUsers(globalStartDate)
	validUserMap := lo.Associate(validUserIndexes, func(item db.ValidIndexUser) (int32, []int32) {
		return item.ValidIndex, item.UserIds
	})

	// バックエンドのみの処理 在籍期間外をグレーアウトする
	applyNonEnrollmentStyle(defaultPileUps, validUserMap)

	// 積み上げ情報をマージして返却する
	for _, ticket := range allFacilityPileUps {
		// 一日当たりの稼働工数の算出（どこかで使用齟齬が起きたと思われる、操作していて自然な方に合わせる）
		var workPerDay float32
		if len(ticket.UserIds) > 0 {
			// アサイン済みの場合は一日当たりの稼働は総工数 / 総稼動可能日数
			workPerDay = float32(*ticket.Estimate) / ticket.NumberOfWorkerByDay
		} else {
			// 未アサインの場合は総工数 / 営業日 / 人数
			workPerDay = float32(*ticket.Estimate) / float32(ticket.NumberOfWorkDay) / float32(*ticket.NumberOfWorker)
		}
		// 対象部署の取得
		// 確定の場合
		// TODO: スタイルとエラーの適応。一旦は数字が合うのかフロントとの結合も大事なので保留とする
		// アサイン済みの場合は 8で割って 1より大きいとエラー
		// それ以外は部署の人数を超過したらエラー フロントも直す必要がある。
		if len(ticket.UserIds) > 0 {
			// アサイン済みの場合 [担当者] [積み上げ、アサイン済み]に計上する
			for _, validIndex := range ticket.ValidIndexes {
				for _, userId := range ticket.UserIds {
					// アサイン済みの場合はユーザーから対象のpileUpsを指定する（部署が指定されていないケースがあるため）
					targetPileUp, exists := lo.Find(defaultPileUps, func(item openapi_models.DefaultPileUp) bool {
						return item.DepartmentId == userMap[userId].DepartmentId
					})

					if !exists {
						continue
					}
					// 設備の期間外の場合は処理を中断。
					if len(targetPileUp.Labels) <= int(validIndex) {
						continue
					}
					// 在籍期間外の場合は対象外とする。
					if !lo.Contains(validUserMap[validIndex], userId) {
						continue
					}

					// 部署人数合計 validUserMapで既にRoleは絞り込まれているので特に対応無
					numberOfDepartmentUsers := float64(len(
						lo.Filter(departmentUserMap[targetPileUp.DepartmentId], func(item db.User, index int) bool {
							return lo.Contains(validUserMap[validIndex], *item.Id)
						}),
					))

					// 部署への積み上げ(共通)
					targetPileUp.Labels[validIndex] += workPerDay
					// エラー判定
					if pileUpLabelFormat(targetPileUp.Labels[validIndex]) > numberOfDepartmentUsers {
						applyErrorStyle(&targetPileUp.Styles[validIndex])
					}

					if facilityMap[ticket.FacilityId].Type == constants.FacilityTypeOrdered {
						// アサイン済みへの積み上げ
						targetPileUp.AssignedUser.Labels[validIndex] += workPerDay
						if pileUpLabelFormat(targetPileUp.AssignedUser.Labels[validIndex]) > numberOfDepartmentUsers {
							applyErrorStyle(&targetPileUp.AssignedUser.Styles[validIndex])
						}
						// ユーザーの足し上げ処理
						targetUserPileUp, userExists := lo.Find(targetPileUp.AssignedUser.Users, func(item openapi_models.PileUpByPerson) bool {
							return *item.User.Id == userId
						})
						if userExists {
							targetUserPileUp.Labels[validIndex] += workPerDay
							if pileUpLabelFormat(targetUserPileUp.Labels[validIndex]) > 1 {
								applyErrorStyle(&targetUserPileUp.Styles[validIndex])
							}
						}
					} else {
						// 未確定の場合 への積み上げ
						targetPileUp.NoOrdersReceivedPileUp.Labels[validIndex] += workPerDay
						if pileUpLabelFormat(targetPileUp.NoOrdersReceivedPileUp.Labels[validIndex]) > numberOfDepartmentUsers {
							applyErrorStyle(&targetPileUp.NoOrdersReceivedPileUp.Styles[validIndex])
						}

						// 未確定の場合の設備への積み上げ
						targetFacilityPileUp, facilityExists := lo.Find(targetPileUp.NoOrdersReceivedPileUp.Facilities, func(item openapi_models.PileUpByFacility) bool {
							return ticket.FacilityId == item.FacilityId
						})
						if facilityExists {
							targetFacilityPileUp.Labels[validIndex] += workPerDay
						}
					}
				}
			}
		} else {
			// 未アサインの場合は部署の設定がされていないと計上しない
			targetPileUp, exists := lo.Find(defaultPileUps, func(item openapi_models.DefaultPileUp) bool {
				return ticket.DepartmentId != nil && item.DepartmentId == *ticket.DepartmentId
			})
			if !exists {
				continue
			}

			// 未アサインの場合 [未アサインのその設備] [積み上げ、未アサイン積み上げ]に計上する
			for _, validIndex := range ticket.ValidIndexes {
				if len(targetPileUp.Labels) <= int(validIndex) {
					continue
				}
				// 部署人数合計
				numberOfDepartmentUsers := float64(len(
					lo.Filter(departmentUserMap[targetPileUp.DepartmentId], func(item db.User, index int) bool {
						return lo.Contains(validUserMap[validIndex], *item.Id)
					}),
				))

				// 部署への積み上げ
				targetPileUp.Labels[validIndex] += workPerDay
				if pileUpLabelFormat(targetPileUp.Labels[validIndex]) > numberOfDepartmentUsers {
					applyErrorStyle(&targetPileUp.Styles[validIndex])
				}

				if facilityMap[ticket.FacilityId].Type == constants.FacilityTypeOrdered {

					// 未アサインへの積み上げ
					targetPileUp.UnAssignedPileUp.Labels[validIndex] += workPerDay
					if pileUpLabelFormat(targetPileUp.UnAssignedPileUp.Labels[validIndex]) > numberOfDepartmentUsers {
						applyErrorStyle(&targetPileUp.UnAssignedPileUp.Styles[validIndex])
					}

					// 未アサインの設備の積み上げ
					targetFacilityPileUp, facilityExists := lo.Find(targetPileUp.UnAssignedPileUp.Facilities, func(item openapi_models.PileUpByFacility) bool {
						return ticket.FacilityId == item.FacilityId
					})
					if facilityExists {
						targetFacilityPileUp.Labels[validIndex] += workPerDay
					}
				} else {
					// 未確定の場合 への積み上げ
					targetPileUp.NoOrdersReceivedPileUp.Labels[validIndex] += workPerDay

					if pileUpLabelFormat(targetPileUp.NoOrdersReceivedPileUp.Labels[validIndex]) > numberOfDepartmentUsers {
						applyErrorStyle(&targetPileUp.NoOrdersReceivedPileUp.Styles[validIndex])
					}

					// 未確定の場合の設備への積み上げ
					targetFacilityPileUp, facilityExists := lo.Find(targetPileUp.NoOrdersReceivedPileUp.Facilities, func(item openapi_models.PileUpByFacility) bool {
						return ticket.FacilityId == item.FacilityId
					})
					if facilityExists {
						targetFacilityPileUp.Labels[validIndex] += workPerDay
					}
				}
			}
		}
	}

	return openapi_models.GetDefaultPileUpsResponse{
		DefaultPileUps:  defaultPileUps,
		GlobalStartDate: globalStartDate,
		DefaultValidUserIndexes: lo.Map(validUserIndexes, func(item db.ValidIndexUser, index int) openapi_models.DefaultValidIndexUsers {
			return openapi_models.DefaultValidIndexUsers{
				UserIds:    item.UserIds,
				ValidIndex: item.ValidIndex,
				IsHoliday: item.IsHoliday,
			}
		}),
	}, nil
}

// getDefaultPileUps 山積みの初期化を行う。期間は設備の最小開始日、最大終了日とする。
func getDefaultPileUps(departments []db.Department, defaultPileUps []openapi_models.DefaultPileUp, departmentUserMap map[int32][]db.User, facilities []db.Facility) ([]openapi_models.DefaultPileUp, time.Time) {
	minFacility := lo.MinBy(facilities, func(a db.Facility, b db.Facility) bool {
		return a.TermFrom.Before(b.TermFrom)
	})
	maxFacility := lo.MaxBy(facilities, func(a db.Facility, b db.Facility) bool {
		return a.TermTo.After(b.TermTo)
	})
	days := int(maxFacility.TermTo.Sub(minFacility.TermFrom).Hours() / 24)

	for _, department := range departments {
		// 配列は基本日数分用意する
		// PileUpByPersonの中身はその部署に所属する人分用意する
		// PileUpByFacilityの中身は設備分用意する
		defaultPileUps = append(defaultPileUps, openapi_models.DefaultPileUp{
			DepartmentId: *department.Id,
			Labels:       make([]float32, days),
			Styles:       createDefaultStyles(days),
			Display:      false,
			AssignedUser: openapi_models.AssignedPileUp{
				Users: lo.Map(departmentUserMap[*department.Id], func(item db.User, index int) openapi_models.PileUpByPerson {
					return openapi_models.PileUpByPerson{
						User: openapi_models.User{
							Id:               item.Id,
							DepartmentId:     item.DepartmentId,
							LimitOfOperation: item.LimitOfOperation,
							LastName:         item.LastName,
							FirstName:        item.FirstName,
							Password:         "",
							Email:            item.Email,
							CreatedAt:        item.CreatedAt,
							UpdatedAt:        item.UpdatedAt,
							Role:             item.Role,
						},
						Labels:   make([]float32, days),
						Styles:   createDefaultStyles(days),
						HasError: false,
					}
				}),
				Labels:  make([]float32, days),
				Styles:  createDefaultStyles(days),
				Display: false,
			},
			UnAssignedPileUp: openapi_models.UnAssingedPileUp{
				Facilities: lo.Map(facilities, func(item db.Facility, index int) openapi_models.PileUpByFacility {
					return openapi_models.PileUpByFacility{
						FacilityId: *item.Id,
						Labels:     make([]float32, days),
						Styles:     createDefaultStyles(days),
						HasError:   false,
					}
				}),
				Labels:  make([]float32, days),
				Styles:  createDefaultStyles(days),
				Display: false,
			},
			NoOrdersReceivedPileUp: openapi_models.NoOrdersReceivedPileUp{
				Facilities: lo.Map(facilities, func(item db.Facility, index int) openapi_models.PileUpByFacility {
					return openapi_models.PileUpByFacility{
						FacilityId: *item.Id,
						Labels:     make([]float32, days),
						Styles:     createDefaultStyles(days),
						HasError:   false,
					}
				}),
				Labels:  make([]float32, days),
				Styles:  createDefaultStyles(days),
				Display: false,
			},
		})
	}

	return defaultPileUps, minFacility.TermFrom
}

func applyNonEnrollmentStyle(defaultPileUps []openapi_models.DefaultPileUp, validUserMap map[int32][]int32) {
	for _, defaultPileUp := range defaultPileUps {
		for _, userPileUp := range defaultPileUp.AssignedUser.Users {
			for index, _ := range userPileUp.Styles {
				if !lo.Contains(validUserMap[int32(index)], *userPileUp.User.Id) {
					applyGrayOutStyle(&userPileUp.Styles[index])
				}
			}
		}
	}
}

func pileUpLabelFormat(v float32) float64 {
	return math.Round(float64(v)*10/8) / 10
}

func createDefaultStyles(days int) []map[string]interface{} {
	result := make([]map[string]interface{}, days)
	for i := 0; i < days; i++ {
		result[i] = make(map[string]interface{})
	}
	return result
}

func applyErrorStyle(v *map[string]interface{}) {
	if *v == nil {
		*v = make(map[string]interface{})
	}
	(*v)["color"] = errorStyle
}

func applyGrayOutStyle(v *map[string]interface{}) {
	if *v == nil {
		*v = make(map[string]interface{})
	}
	(*v)["background-color"] = grayOutStyle
}
