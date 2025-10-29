package helpers

import (
	"net/http"
	"strings"

	"github.com/goawwer/yamyard/utils"
)

type FilterAndSortingParameters struct {
	Key      string
	Value    string
	Sort     string
	Order    string
	Username string
}

func GetValidQueryParameters(r *http.Request, model any) FilterAndSortingParameters {
	var params FilterAndSortingParameters

	baseQuery := r.URL.Query()

	params.Order = baseQuery.Get("order")
	if params.Order != "asc" && params.Order != "desc" {
		params.Order = "asc"
	}
	baseQuery.Del("order")

	params.Sort = baseQuery.Get("sort")

	fieldMap := utils.GetDBFieldMap(model)
	if dbField, ok := fieldMap[strings.ToLower(params.Sort)]; ok {
		params.Sort = dbField
	}

	baseQuery.Del("sort")

	params.Username = baseQuery.Get("username")
	baseQuery.Del("username")

	for k, v := range baseQuery {
		params.Key = k
		params.Value = strings.Join(v, "")
		break
	}

	return params
}

func InitHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}
