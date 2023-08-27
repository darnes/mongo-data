package plugin

import (
	// "context"
	"fmt"
	"reflect"

	"testing"

	"go.mongodb.org/mongo-driver/bson"
	// "github.com/grafana/grafana-plugin-sdk-go/backend"
)

// func typedArray(el interface{}) []string {
// 	if reflect.TypeOf(el).Name() == "string" {
// 		return []string{el}
// 	}
// }

func bsonToLists(dbResp []bson.M) {
	if len(dbResp) > 0 {
		for key, v := range dbResp[0] {
			if reflect.TypeOf(v).Name() == "string" {
				l := getValuesTest[string](key, dbResp)
				fmt.Printf("string list %v of type `%s` \n", l, reflect.TypeOf(l))
			} else if reflect.TypeOf(v).Name() == "int" {
				l := getValuesTest[int](key, dbResp)
				fmt.Printf("int list %v of type `%s` \n", l, reflect.TypeOf(l))

			} else {
				fmt.Printf("unknown type `%s` \n", reflect.TypeOf(v).Name())
			}

		}
	}
}

func getValuesTest[T int | string](key string, source []bson.M) []T {
	var res []T

	for _, se := range source {
		el, _ := se[key]
		tEl, _ := el.(T)
		res = append(res, tEl)
	}
	return res
}

func TestData(t *testing.T) {
	inp := []bson.M{}
	inp = append(inp, bson.M{"a": "asd", "b": 3})
	inp = append(inp, bson.M{"a": "asdddd", "b": 4})
	// fmt.Printf("inp %v", inp)
	bsonToLists(inp)
	fmt.Println("next one")
	// fmt.Printf("frames %v", bsonToFrames(inp))
	// ds := Datasource{}

	// resp, err := ds.QueryData(
	// 	context.Background(),
	// 	&backend.QueryDataRequest{
	// 		Queries: []backend.DataQuery{
	// 			{RefID: "A"},
	// 		},
	// 	},
	// )
	// if err != nil {
	// 	t.Error(err)
	// }

	// if len(resp.Responses) != 1 {
	// 	t.Fatal("QueryData must return a response")
	// }
	// t.("something")

	// t.Fatal("Failedssss")
}
