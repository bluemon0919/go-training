/*
 * webapi RowData
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package webapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// A RowdataApiController binds http requests to an api service and writes the service results to the http response
type RowdataApiController struct {
	service RowdataApiServicer
}

// NewRowdataApiController creates a default api controller
func NewRowdataApiController(s RowdataApiServicer) Router {
	return &RowdataApiController{service: s}
}

// Routes returns all of the api route for the RowdataApiController
func (c *RowdataApiController) Routes() Routes {
	return Routes{
		{
			"GetRowData",
			strings.ToUpper("Get"),
			"/api/rowdata",
			c.GetRowData,
		},
	}
}

// GetRowData - get rowdata information
func (c *RowdataApiController) GetRowData(w http.ResponseWriter, r *http.Request) {
	/*
		query := r.URL.Query()
		title := query.Get("title")
		result, err := c.service.GetRowData(title)
		if err != nil {
			w.WriteHeader(500)
			return
		}
	*/
	var rowdatas []RowData
	for i := 0; i < 15; i++ {
		rowdata := RowData{
			Title:  "みんなのラジオ番組",
			Detail: fmt.Sprintf("月曜日 %d:00-%d:00", i, i+1),
		}
		rowdatas = append(rowdatas, rowdata)
	}
	json.NewEncoder(w).Encode(rowdatas)
	//EncodeJSONResponse(result, nil, w)
}
