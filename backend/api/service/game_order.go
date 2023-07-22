package service

import (
	"github.com/golang-module/carbon/v2"
	"github.com/kenkonno/gantt-chart-proto/backend/api/utils"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
)

// TODO: UserIDに紐づけて見たそうなやつを出す
func GetGameOrder(baseDate carbon.Carbon) []db.Game {
	gameRep := repository.NewGameRepository()

	games := gameRep.GetDefaultOrder(baseDate.Format(utils.DateTimeFormatUTC))

	return games

}
