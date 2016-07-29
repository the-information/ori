package dsimport

import (
	"encoding/json"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"reflect"
	"testing"
	"time"
)

func Test_importEntity(t *testing.T) {

	e := entity{}
	results := []datastore.Property{}

	fullTestJSON := []byte(`{
			"CreatedAt": {
				"Type": "time",
				"Value": "2011-06-12T12:30:00Z"
			},
			"Name": "Jane Q. Public",
			"LotteryNumbers": [0,7,19,36],
			"BMI": 21.2
		}`)

	expectedResults := []datastore.Property{
		{
			Name:  "CreatedAt",
			Value: time.Date(2011, 06, 12, 12, 30, 0, 0, time.UTC),
		},
		{
			Name:  "Name",
			Value: "Jane Q. Public",
		},
		{
			Name:     "LotteryNumbers",
			Multiple: true,
			Value:    int64(0),
		},
		{
			Name:     "LotteryNumbers",
			Multiple: true,
			Value:    int64(7),
		},
		{
			Name:     "LotteryNumbers",
			Multiple: true,
			Value:    int64(19),
		},
		{
			Name:     "LotteryNumbers",
			Multiple: true,
			Value:    int64(36),
		},
		{
			Name:  "BMI",
			Value: float64(21.2),
		},
	}

	if err := json.Unmarshal(fullTestJSON, &e); err != nil {
		t.Errorf("Unexpected error %s during json.Unmarshal", err)
	} else if err := e.FetchProperties(context.Background(), &results); err != nil {
		t.Errorf("Unexpected error %s during FetchProperties", err)
	}

	if !reflect.DeepEqual(results, expectedResults) {
		t.Errorf("Got unexpected result set %+v", results)
	}

}
