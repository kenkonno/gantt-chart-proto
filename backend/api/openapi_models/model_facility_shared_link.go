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

type FacilitySharedLink struct {

	Id *int32 `json:"id,omitempty"`

	FacilityId int32 `json:"facility_id" binding:"min=1"`

	Uuid *string `json:"uuid,omitempty" binding:"min=1"`

	CreatedAt time.Time `json:"created_at,omitempty"`

	UpdatedAt int `json:"updated_at,omitempty"`
}