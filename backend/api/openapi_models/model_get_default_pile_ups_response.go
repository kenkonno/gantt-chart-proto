/*
 * GanttChartApi
 *
 * No description provided (generated by Openapi Generator https://github.com/openapi_modelstools/openapi-generator)
 *
 * API version: 1.0
 * Generated by: OpenAPI Generator (https://openapi_models-generator.tech)
 */

package openapi_models

import (
	"time"
)

type GetDefaultPileUpsResponse struct {
	DefaultPileUps []DefaultPileUp `json:"defaultPileUps"`

	GlobalStartDate time.Time `json:"globalStartDate"`

	DefaultValidUserIndexes []DefautValidIndexUsers `json:"defaultValidUserIndexes"`
}
