/*
 * GanttChartApi
 *
 * No description provided (generated by Openapi Generator https://github.com/openapi_modelstools/openapi-generator)
 *
 * API version: 1.0
 * Generated by: OpenAPI Generator (https://openapi_models-generator.tech)
 */

package openapi_models

type GetPileUpsResponseListInner struct {

	FacilityId int32 `json:"facilityId" form:"facilityId"`

	Holidays []Holiday `json:"holidays" form:"holidays"`

	GanttGroups []GanttGroup `json:"ganttGroups" form:"ganttGroups"`
}
