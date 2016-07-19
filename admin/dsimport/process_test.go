package dsimport

import (
	"github.com/qedus/nds"
	"google.golang.org/appengine/aetest"
	"google.golang.org/appengine/datastore"
	"reflect"
	"strings"
	"testing"
	"time"
)

type fakeUser struct {
	Name           string
	LotteryNumbers []int64
	CreatedAt      time.Time
}

func TestProcess(t *testing.T) {

	ctx, done, _ := aetest.NewContext()
	defer done()

	smallTest := strings.NewReader(`{
		"User/jsmith": {
			"Name": "John Smith",
			"LotteryNumbers": [1,2,3,4,5],
			"CreatedAt": {
				"Type": "time",
				"Value": "1993-05-01T12:31:00.000Z"
			}
		},
		"User/jdoe": {
			"Name": "Jane Doe",
			"LotteryNumbers": [2,4,6,8,10],
			"CreatedAt": {
				"Type": "time",
				"Value": "1992-01-30T08:01:00.000Z"
			}
		}
	}`)

	if err := Process(ctx, smallTest); err != nil {
		t.Errorf("Unexpected error %s", err)
	}

	key1 := datastore.NewKey(ctx, "User", "jdoe", 0, nil)
	key2 := datastore.NewKey(ctx, "User", "jsmith", 0, nil)

	value1 := fakeUser{}
	value2 := fakeUser{}

	if err := nds.Get(ctx, key1, &value1); err != nil {
		t.Errorf("Unexpected error %s retrieving value1", err)
	} else if err := nds.Get(ctx, key2, &value2); err != nil {
		t.Errorf("Unexpected error %s retrieving value2", err)
	}

	value1.CreatedAt = value1.CreatedAt.UTC()
	value2.CreatedAt = value2.CreatedAt.UTC()

	if !reflect.DeepEqual(value1, fakeUser{
		Name:           "Jane Doe",
		LotteryNumbers: []int64{2, 4, 6, 8, 10},
		CreatedAt:      time.Date(1992, 1, 30, 8, 1, 0, 0, time.UTC),
	}) {
		t.Errorf("Unexpected value in value1: %+v", value1)
	}

	if !reflect.DeepEqual(value2, fakeUser{
		Name:           "John Smith",
		LotteryNumbers: []int64{1, 2, 3, 4, 5},
		CreatedAt:      time.Date(1993, 5, 1, 12, 31, 0, 0, time.UTC),
	}) {
		t.Errorf("Unexpected value in value2: %+v", value1)
	}

}
