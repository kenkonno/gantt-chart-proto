/*
 * GanttChartApi
 *
 * No description provided (generated by Openapi Generator https://github.com/openapi_modelstools/openapi-generator)
 *
 * API version: 1.0
 * Generated by: OpenAPI Generator (https://openapi_models-generator.tech)
 */

package openapi_models

type NoOrdersReceivedPileUp struct {

	Facilities []PileUpByFacility `json:"facilities" form:"facilities"`

	Labels []float32 `json:"labels" form:"labels"`

	Styles []map[string]interface{} `json:"styles" form:"styles"`

	Display bool `json:"display" form:"display"`
}
