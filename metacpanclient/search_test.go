package metacpanclient

import (
	goJson "encoding/json"
	"reflect"
	"testing"
)

func TestSearchAdjust(t *testing.T) {
	t.Parallel()
	args := tUnmarshal(testTestNewSearchInput0)
	s := tNewSearch(args)
	got := tRemarshal(s.Query)
	expected := tUnmarshal(testTestNewSearchOutput0)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("got %v, expected %v", got, expected)
	}
	args = tUnmarshal(testTestNewSearchInput1)
	s = tNewSearch(args)
	got = tRemarshal(s.Query)
	expected = tUnmarshal(testTestNewSearchOutput1)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("got %v, expected %v", got, expected)
	}
	args = tUnmarshal(testTestNewSearchInput2)
	s = tNewSearch(args)
	got = tRemarshal(s.Query)
	expected = tUnmarshal(testTestNewSearchOutput2)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("got %#v, expected %#v", got, expected)
	}
}

type tFakeResult = *wrapper[struct{}]

func tNewSearch(m map[string]interface{}) *search[tFakeResult] {
	mc := &Client{}

	s, err := newSearch[tFakeResult](&searchConfig{
		mc:    mc,
		query: m,
	})
	if err != nil {
		panic(err)
	}
	return s
}

func tRemarshal(m map[string]interface{}) map[string]interface{} {
	var ret map[string]interface{}
	b, err := goJson.Marshal(m)
	if err != nil {
		panic(err)
	}
	err = goJson.Unmarshal(b, &ret)
	if err != nil {
		panic(err)
	}
	return ret
}

func tUnmarshal(s string) map[string]interface{} {
	var m map[string]interface{}
	b := []byte(s)
	err := goJson.Unmarshal(b, &m)
	if err != nil {
		panic(err)
	}
	return m
}

const (
	testEnableDebug         = false
	testTestNewSearchInput0 = `{
	    "either": [
	        {
	            "name": "Dave *"
	        },
	        {
	            "name": "David *"
	        }
	    ]
	}`
	testTestNewSearchOutput0 = `{
	    "bool": {
	        "minimum_should_match": 1,
	        "should": [
	            {
	                "wildcard": {
	                    "name": "Dave *"
	                }
	            },
	            {
	                "wildcard": {
	                    "name": "David *"
	                }
	            }
	        ]
	    }
	}`
	testTestNewSearchInput1 = `{
	    "either": [
	        {
	            "name": "Dave *"
	        },
	        {
	            "name": "David *"
	        }
	    ],
	    "not": [
	        {
	            "name": "Dave S*"
	        },
	        {
	            "name": "David S*"
	        }
	    ]
	}`
	testTestNewSearchOutput1 = `{
	    "bool": {
	        "minimum_should_match": 1,
	        "must_not": [
	            {
	                "wildcard": {
	                    "name": "Dave S*"
	                }
	            },
	            {
	                "wildcard": {
	                    "name": "David S*"
	                }
	            }
	        ],
	        "should": [
	            {
	                "wildcard": {
	                    "name": "Dave *"
	                }
	            },
	            {
	                "wildcard": {
	                    "name": "David *"
	                }
	            }
	        ]
	    }
	}`
	testTestNewSearchInput2 = `{
	    "either": [
	        {
	            "all": [
	                {
	                    "name": "Dave *"
	                },
	                {
	                    "email": "*gmail.com"
	                }
	            ]
	        },
	        {
	            "all": [
	                {
	                    "name": "David *"
	                },
	                {
	                    "email": "*gmail.com"
	                }
	            ]
	        }
	    ]
	}`
	testTestNewSearchOutput2 = `{
	    "bool": {
	        "minimum_should_match": 1,
	        "should": [
	            {
	                "bool": {
	                    "must": [
	                        {
	                            "wildcard": {
	                                "name": "Dave *"
	                            }
	                        },
	                        {
	                            "wildcard": {
	                                "email": "*gmail.com"
	                            }
	                        }
	                    ]
	                }
	            },
	            {
	                "bool": {
	                    "must": [
	                        {
	                            "wildcard": {
	                                "name": "David *"
	                            }
	                        },
	                        {
	                            "wildcard": {
	                                "email": "*gmail.com"
	                            }
	                        }
	                    ]
	                }
	            }
	        ]
	    }
	}`
)
