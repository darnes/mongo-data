package plugin

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	// "time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/data"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Make sure Datasource implements required interfaces. This is important to do
// since otherwise we will only get a not implemented error response from plugin in
// runtime. In this example datasource instance implements backend.QueryDataHandler,
// backend.CheckHealthHandler interfaces. Plugin should not implement all these
// interfaces- only those which are required for a particular task.
var (
	_ backend.QueryDataHandler      = (*Datasource)(nil)
	_ backend.CheckHealthHandler    = (*Datasource)(nil)
	_ instancemgmt.InstanceDisposer = (*Datasource)(nil)
)

type configModel struct {
	ConnectionString string `json:"mongoConnectionString"`
	// Collection       string `json:"collection"`
}

// NewDatasource creates a new datasource instance.
func NewDatasource(s backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {

	backend.Logger.Info("new-data-source", "set", s)
	cm := configModel{}
	err := json.Unmarshal(s.JSONData, &cm)
	if err != nil {
		backend.Logger.Error("failed to unmarshal config from JSONData", err)
		return nil, errors.New("wrong configuration")
	}
	backend.Logger.Info("checking config ", "c", cm, "cs", cm.ConnectionString)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(cm.ConnectionString))
	// here are mostly URL parsing errors
	if err != nil {
		backend.Logger.Info("failed connect to mongo", "err", err)
		backend.Logger.Error("failed connect to DB", err)
		return nil, err
	}
	backend.Logger.Info("client after connection", "c", client)

	// actual connection validity check
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		backend.Logger.Info("failed ping mongo", "err", err)
		backend.Logger.Error("failed ping DB", err)
		return nil, err
	}

	return &Datasource{Client: client}, nil
}

// Datasource is an example datasource which can respond to data queries, reports
// its health and has streaming skills.
type Datasource struct {
	Client *mongo.Client
}

// Dispose here tells plugin SDK that plugin wants to clean up resources when a new instance
// created. As soon as datasource settings change detected by SDK old datasource instance will
// be disposed and a new one will be created using NewSampleDatasource factory function.
func (d *Datasource) Dispose() {
	// Clean up datasource instance resources.
	err := d.Client.Disconnect(context.TODO())
	if err != nil {
		backend.Logger.Error("failed-to-disconnect-db", "err", err)
	}
}

// QueryData handles multiple queries and returns multiple responses.
// req contains the queries []DataQuery (where each query contains RefID as a unique identifier).
// The QueryDataResponse contains a map of RefID to the response for each query, and each response
// contains Frames ([]*Frame).
func (d *Datasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	// create response struct
	response := backend.NewQueryDataResponse()
	// db.movies.aggregate({$project: {"_id":0, type:1, r:1}})

	// loop over queries and execute them individually.
	for _, q := range req.Queries {
		res := d.query(ctx, req.PluginContext, q)
		// save the response in a hashmap
		// based on with RefID as identifier
		response.Responses[q.RefID] = res
	}

	return response, nil
}

type queryModel struct {
	QueryText  string `json:"queryText"`
	Collection string `json:"collection"`
	DbName     string `json:"db"`
}

func (d *Datasource) query(_ context.Context, pCtx backend.PluginContext, query backend.DataQuery) backend.DataResponse {
	var response backend.DataResponse

	// reading query
	var qm queryModel

	// Unmarshal the JSON into our queryModel.
	err := json.Unmarshal(query.JSON, &qm)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("json unmarshal: %v", err.Error()))
	}
	backend.Logger.Info("parsed query", "qm", qm)
	var extJsonQuery interface{}

	err = bson.UnmarshalExtJSON([]byte(qm.QueryText), true, &extJsonQuery)
	if err != nil {
		backend.Logger.Error("failed to parse query as bson", err)
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("failed to parse query as bson %v", err))
	}

	backend.Logger.Info("parsed json-query: ", "q", extJsonQuery)

	coll := d.Client.Database(qm.DbName).Collection(qm.Collection)

	cursor, err := coll.Aggregate(context.TODO(), extJsonQuery)
	if err != nil {
		backend.Logger.Error("query failed", err)
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("query failed %v", err))
	}

	// getting the result
	var results []bson.M
	cursor.All(context.TODO(), &results)
	backend.Logger.Info("query result", "result", results)

	// create data frame response.
	// For an overview on data frames and how grafana handles them:
	// https://grafana.com/docs/grafana/latest/developers/plugins/data-frames/

	frame := data.NewFrame("response")
	// making response out of results

	frame.Fields = append(frame.Fields,
		bsonToFrames(results)...,
	)

	// add the frames to the response.
	response.Frames = append(response.Frames, frame)

	return response
}

func bsonToFrames(dbResp []bson.M) []*data.Field {
	var fields []*data.Field
	if len(dbResp) > 0 {
		for key, v := range dbResp[0] {
			if reflect.TypeOf(v).Name() == "string" {
				fields = append(fields, data.NewField(key, nil, getValues[string](key, dbResp)))

			} else if reflect.TypeOf(v).Name() == "int" {
				fields = append(fields, data.NewField(key, nil, getValues[int](key, dbResp)))

			} else if reflect.TypeOf(v).Name() == "int32" {
				fields = append(fields, data.NewField(key, nil, getValues[int32](key, dbResp)))
			} else if reflect.TypeOf(v).Name() == "float64" {
				fields = append(fields, data.NewField(key, nil, getValues[float64](key, dbResp)))
			} else {
				// $group: {_id: "$type", total: {$avg: "$r"}, }}
				//ObjectID is not implemented yet
				// it has nice .toString  method can be used.
				// no idea how handy it may be.
				backend.Logger.Info("unknown type", "t", reflect.TypeOf(v).Name())
			}
		}
	}

	// 	backend.Logger.Info("added value", "key", key, "val", values)
	// }
	return fields
}

func getValues[T int32 | int | string | float64](key string, source []bson.M) []T {
	var res []T

	for _, se := range source {
		el := se[key]
		tEl, _ := el.(T)
		res = append(res, tEl)
	}
	return res
}

// CheckHealth handles health checks sent from Grafana to the plugin.
// The main use case for these health checks is the test button on the
// datasource configuration page which allows users to verify that
// a datasource is working as expected.
func (d *Datasource) CheckHealth(_ context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	var status = backend.HealthStatusError
	var message = "Failed to verify DB connection "

	if d.Client != nil {
		err := d.Client.Ping(context.TODO(), nil)
		if err != nil {
			backend.Logger.Error("failed to ping client on health-check", "e", err)
			return nil, errors.New("can not establish connection to DB from health-check")
		}
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusOk,
			Message: "connection to DB established",
		}, nil
	}
	backend.Logger.Info("health-check: d.client is nill, defaulting to failed check")
	return &backend.CheckHealthResult{
		Status:  status,
		Message: message,
	}, nil
}
