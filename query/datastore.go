package query

import (
	"bytes"
	"errors"
	"fmt"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"net/url"
	"strconv"
	"time"
)

var ErrDatastoreLimitTooLarge = errors.New("The limit specified by this query was greater than 1000")

/*
DatastoreWithValues produces a *datastore.Query given a set of URL query parameters.
Given a query string that looks like the following,
	Size=large&Price_gt=100&Price_lt=500&ForSale=true&SaleUntil_le=2020-01-01T00%3A00%3A00Z&_order=-Price
This code:
	q := query.DatastoreWithValues("Widget", r.URL.Query())
is equivalent to this code:
	q := datastore.NewQuery("Widget").
		Filter("Size =", "large").
		Filter("Price >", 100).
		Filter("Price <", 500).
		Filter("ForSale =", true).
		Filter("SaleUntil", time.Parse("2020-01-01T00:00:00Z")
		Order("-Price")

The following four query keys are treated specially:

	"_order" will be used to set the ordering of the query results.
	"_start" and "_end", if supplied, are interpreted as encoded datastore.Cursor objects.
	If they are not valid encoded cursors, DatastoreWithValues will fail.
	"_limit" is interpreted as an integer to be used with q.Limit(). If its value is greater than 1000 or it cannot
	be converted to an integer, DatastoreWithValues will fail.

All other query parameters are interpreted as filters according to the following algorithm:

The query key's last three characters are checked. If they are one of the following four values,
the last three characters are stripped from the key and the given operator is used. Otherwise,
strict equality is assumed.

	- "_gt": ">"
	- "_lt": "<"
	- "_ge": ">="
	- "_le": "<="

The query value will be unmarshaled into an interface{} as follows. If the query value
is surrounded by quotes (i.e., it looked like `foo=%2cbar%2d` in the query string,
so it's actually `foo="bar"` unescaped), it will be treated as a string and unwrapped.
Non-quoted query values will be unmarshaled into an interface{} value as follows:

	- Literal "true": boolean true.
	- Literal "false": boolean false.
	- An int64 as interpreted by strconv.ParseInt: int64.
	- A float64 as interpreted by strconv.ParseFloat: float64.
	- A base64 string datastore.DecodeKey can interpret: *datastore.Key.
	- An ISO8601/RFC3339 time string as interpreted by time.Parse: time.Time.
	- A string with the format lat_(float64)_lng_(float64): appengine.GeoPoint.
	- All other values: string.

*/
func DatastoreWithValues(kind string, params url.Values) (*datastore.Query, error) {

	var buf bytes.Buffer
	q := datastore.NewQuery(kind)

	for k, _ := range params {
		v := params.Get(k)
		switch k {
		case "_order":
			q = q.Order(v)
		case "_start":
			if dsCursor, err := datastore.DecodeCursor(v); err != nil {
				return nil, err
			} else {
				q = q.Start(dsCursor)
			}
		case "_end":
			if dsCursor, err := datastore.DecodeCursor(v); err != nil {
				return nil, err
			} else {
				q = q.End(dsCursor)
			}
		case "_limit":
			if count, err := strconv.Atoi(v); err != nil {
				return nil, err
			} else if count > 1000 {
				return nil, ErrDatastoreLimitTooLarge
			} else {
				q = q.Limit(count)
			}
		default:
			q = q.Filter(getFilterStr(k, &buf), getFilterValue(v))
		}
	}

	return q, nil

}

func getFilterStr(k string, buf *bytes.Buffer) string {

	fieldName := k[:len(k)-3]
	operator := k[len(k)-3:]

	buf.Reset()

	switch operator {
	default:
		buf.WriteString(k)
		buf.WriteString(" =")
		return buf.String()
	case "_gt":
		buf.WriteString(fieldName)
		buf.WriteString(" >")
		return buf.String()
	case "_ge":
		buf.WriteString(fieldName)
		buf.WriteString(" >=")
		return buf.String()
	case "_lt":
		buf.WriteString(fieldName)
		buf.WriteString(" <")
		return buf.String()
	case "_le":
		buf.WriteString(fieldName)
		buf.WriteString(" <=")
		return buf.String()
	}

}

func getFilterValue(v string) interface{} {

	var lat, lng float64
	if len(v) > 0 && v[0] == '"' && v[len(v)-1] == '"' {
		return v[1 : len(v)-1]
	} else if v == "true" {
		return true
	} else if v == "false" {
		return false
	} else if intVal, err := strconv.ParseInt(v, 10, 64); err == nil {
		return intVal
	} else if floatVal, err := strconv.ParseFloat(v, 64); err == nil {
		return floatVal
	} else if dsKey, err := datastore.DecodeKey(v); err == nil {
		return dsKey
	} else if date, err := time.Parse(time.RFC3339, v); err == nil {
		return date
	} else if count, err := fmt.Sscanf(v, "lat_%f_lng_%f", &lat, &lng); count == 2 && err == nil {
		return appengine.GeoPoint{lat, lng}
	} else {
		return v
	}

}
