package query

import (
	"bytes"
	"google.golang.org/appengine"
	"google.golang.org/appengine/aetest"
	"google.golang.org/appengine/datastore"
	"net/url"
	"testing"
	"time"
)

type widget struct {
	Cost    float64
	Count   int64
	ForSale bool
}

func TestDatastoreWithValues(t *testing.T) {

	ctx, done, _ := aetest.NewContext()
	defer done()

	// set up the datastore with some Widgets!
	results := []widget{}

	// http://stackoverflow.com/questions/25070974/google-app-engine-golang-datastore-query-getall-not-working-locally
	var w widget
	k1, _ := datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "Widget", nil), &widget{1.0, 8, true})
	k2, _ := datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "Widget", nil), &widget{9.0, 10, true})
	k3, _ := datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "Widget", nil), &widget{10.0, 803, true})
	k4, _ := datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "Widget", nil), &widget{9.0, 50, true})
	k5, _ := datastore.Put(ctx, datastore.NewIncompleteKey(ctx, "Widget", nil), &widget{150.0, 10, false})
	datastore.Get(ctx, k1, &w)
	datastore.Get(ctx, k2, &w)
	datastore.Get(ctx, k3, &w)
	datastore.Get(ctx, k4, &w)
	datastore.Get(ctx, k5, &w)

	u, _ := url.Parse("http://example.com/?Cost=9.0&Count_lt=100&ForSale=true&_order=-Count")
	q, err := DatastoreWithValues("Widget", u.Query())
	if err != nil {
		t.Fatalf("Unexpected error %s from DatastoreWithValues", err)
	}

	w0 := widget{9.0, 50, true}
	w1 := widget{9.0, 10, true}

	_, err = q.GetAll(ctx, &results)
	if err != nil {
		t.Errorf("Unexpected error %s from q.GetAll", err)
	} else if len(results) != 2 {
		t.Errorf("Unexpected number of results %d %+v", len(results), results)
	} else if results[0] != w0 {
		t.Errorf("Unexpected results[0]: %+v", results[0])
	} else if results[1] != w1 {
		t.Errorf("Unexpected results[1]: %+v", results[1])
	}

	u, _ = url.Parse("http://example.com/?Cost=9.0&Count_lt=100&ForSale=true&_limit=9999")
	if _, err = DatastoreWithValues("Widget", u.Query()); err != ErrDatastoreLimitTooLarge {
		t.Fatalf("Expected ErrDatastoreLimitTooLarge, got %s", err)
	}

	// TODO(goldibex): test _start and _end
}

func Test_getFilterStr(t *testing.T) {

	var buf bytes.Buffer
	if r := getFilterStr("foo_gt", &buf); r != "foo >" {
		t.Errorf("Expected foo >, got %s", r)
	}
	if r := getFilterStr("bar_lt", &buf); r != "bar <" {
		t.Errorf("Expected bar <, got %s", r)
	}
	if r := getFilterStr("foo_ge", &buf); r != "foo >=" {
		t.Errorf("Expected foo >=, got %s", r)
	}
	if r := getFilterStr("bar_le", &buf); r != "bar <=" {
		t.Errorf("Expected bar <=, got %s", r)
	}
	if r := getFilterStr("foobar", &buf); r != "foobar =" {
		t.Errorf("Expected foobar =, got %s", r)
	}

}

func Test_getFilterValue(t *testing.T) {

	ctx, done, _ := aetest.NewContext()
	defer done()

	if r := getFilterValue("\"true\""); r.(string) != "true" {
		t.Errorf("Expected \"true\", got %s", r)
	}
	if r := getFilterValue("true"); r.(bool) != true {
		t.Errorf("Expected true, got %s", r)
	}
	if r := getFilterValue("false"); r.(bool) != false {
		t.Errorf("Expected false, got %s", r)
	}
	if r := getFilterValue("42"); r.(int64) != 42 {
		t.Errorf("Expected 42, got %s", r)
	}
	if r := getFilterValue("1.99998"); r.(float64) != 1.99998 {
		t.Errorf("Expectd 1.99998, got %s", r)
	}
	k := datastore.NewKey(ctx, "Foo", "", 1, nil).Encode()
	if r := getFilterValue(k).(*datastore.Key); r.Encode() != k {
		t.Errorf("Expected %s, got %s", k, r)
	}
	tm := time.Now()
	if r := getFilterValue(tm.Format(time.RFC3339)).(time.Time); r.Format(time.RFC3339) != tm.Format(time.RFC3339) {
		t.Errorf("Expected %s, got %s", tm, r)
	}
	if r := getFilterValue("lat_1.2_lng_3.4").(appengine.GeoPoint); r.Lat != 1.2 || r.Lng != 3.4 {
		t.Errorf("Got unexpected GeoPoint %+v", r)
	}
	if r := getFilterValue("foo").(string); r != "foo" {
		t.Errorf("Expected foo, got %s", r)
	}
}
