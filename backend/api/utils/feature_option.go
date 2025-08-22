package utils

import (
	"github.com/kenkonno/gantt-chart-proto/backend/api/constants"
	"github.com/kenkonno/gantt-chart-proto/backend/models/db"
	"github.com/kenkonno/gantt-chart-proto/backend/repository"
	"github.com/samber/lo"
)

func HasOption(optionName constants.Feature) bool {

	featureOptionRep := repository.NewFeatureOptionRepository()

	_, exists := lo.Find(featureOptionRep.FindAll(), func(item db.FeatureOption) bool {
		return optionName.Equals(item.Name)
	})

	return exists
}
