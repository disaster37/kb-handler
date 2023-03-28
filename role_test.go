package kbhandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/disaster37/go-kibana-rest/v8/kbapi"
	"github.com/google/go-cmp/cmp"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

var urlrole = fmt.Sprintf("%s/api/security/role/test", baseURL)

func (t *KibanaHandlerTestSuite) TestRoleGet() {

	rawRole := `
{
	"name": "test",
	"metadata" : {
		"version" : 1
	},
	"transient_metadata": {
		"enabled": true
	},
	"elasticsearch": {
		"cluster": [ ],
		"indices": [ ],
		"run_as": [ ]
	},
		"kibana": [
		{
		"base": [
			"read"
		],
		"feature": {},
		"spaces": [
			"marketing"
		]
		},
		{
		"base": [],
		"feature": {
			"discover": [
			"all"
			],
			"visualize": [
			"all"
			],
			"dashboard": [
			"all"
			],
			"dev_tools": [
			"read"
			],
			"advancedSettings": [
			"read"
			],
			"indexPatterns": [
			"read"
			],
			"graph": [
			"all"
			],
			"apm": [
			"read"
			],
			"maps": [
			"read"
			],
			"canvas": [
			"read"
			],
			"infrastructure": [
			"all"
			],
			"logs": [
			"all"
			],
			"uptime": [
			"all"
			]
		},
		"spaces": [
			"sales",
			"default"
		]
		}
	]
}
	`

	roleTest := &kbapi.KibanaRole{}
	if err := json.Unmarshal([]byte(rawRole), roleTest); err != nil {
		panic(err)
	}

	httpmock.RegisterResponder("GET", urlrole, func(req *http.Request) (*http.Response, error) {
		resp := httpmock.NewStringResponse(200, rawRole)
		return resp, nil
	})

	role, err := t.kbHandler.RoleGet("test")
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Empty(t.T(), cmp.Diff(roleTest, role))

	// When error
	httpmock.RegisterResponder("GET", urlrole, httpmock.NewErrorResponder(errors.New("fack error")))
	_, err = t.kbHandler.RoleGet("test")
	assert.Error(t.T(), err)
}

func (t *KibanaHandlerTestSuite) TestRoleDelete() {

	httpmock.RegisterResponder("DELETE", urlrole, func(req *http.Request) (*http.Response, error) {
		resp := httpmock.NewStringResponse(200, "")
		return resp, nil
	})

	err := t.kbHandler.RoleDelete("test")
	if err != nil {
		t.Fail(err.Error())
	}

	// When error
	httpmock.RegisterResponder("DELETE", urlrole, httpmock.NewErrorResponder(errors.New("fack error")))
	err = t.kbHandler.RoleDelete("test")
	assert.Error(t.T(), err)
}

func (t *KibanaHandlerTestSuite) TestRoleUpdate() {

	rawRole := `
{
	"name": "test",
	"metadata" : {
		"version" : 1
	},
	"elasticsearch": {
		"cluster" : [ ],
		"indices" : [ ]
	},
	"kibana": [
		{
		"base": [],
		"feature": {
			"discover": [
			"all"
			],
			"visualize": [
			"all"
			],
			"dashboard": [
			"all"
			],
			"dev_tools": [
			"read"
			],
			"advancedSettings": [
			"read"
			],
			"indexPatterns": [
			"read"
			],
			"graph": [
			"all"
			],
			"apm": [
			"read"
			],
			"maps": [
			"read"
			],
			"canvas": [
			"read"
			],
			"infrastructure": [
			"all"
			],
			"logs": [
			"all"
			],
			"uptime": [
			"all"
			]
		},
		"spaces": [
			"*"
		]
		}
	]
}
	`

	role := &kbapi.KibanaRole{}
	if err := json.Unmarshal([]byte(rawRole), role); err != nil {
		panic(err)
	}

	httpmock.RegisterResponder("PUT", urlrole, func(req *http.Request) (*http.Response, error) {
		resp := httpmock.NewStringResponse(200, rawRole)
		return resp, nil
	})
	httpmock.RegisterResponder("GET", urlrole, func(req *http.Request) (*http.Response, error) {
		resp := httpmock.NewStringResponse(200, rawRole)
		return resp, nil
	})

	err := t.kbHandler.RoleUpdate(role)
	if err != nil {
		t.Fail(err.Error())
	}

	// When error
	httpmock.RegisterResponder("PUT", urlrole, httpmock.NewErrorResponder(errors.New("fack error")))
	err = t.kbHandler.RoleUpdate(role)
	assert.Error(t.T(), err)
}

func (t *KibanaHandlerTestSuite) TestRoleDiff() {
	var actual, expected, original *kbapi.KibanaRole

	rawRole := `
{
	"metadata" : {
		"version" : 1
	},
	"elasticsearch": {},
	"kibana": [
		{
		"feature": {
			"discover": [
			"all"
			],
			"visualize": [
			"all"
			],
			"dashboard": [
			"all"
			],
			"dev_tools": [
			"read"
			],
			"advancedSettings": [
			"read"
			],
			"indexPatterns": [
			"read"
			],
			"graph": [
			"all"
			],
			"apm": [
			"read"
			],
			"maps": [
			"read"
			],
			"canvas": [
			"read"
			],
			"infrastructure": [
			"all"
			],
			"logs": [
			"all"
			],
			"uptime": [
			"all"
			]
		},
		"spaces": [
			"*"
		]
		}
	]
}
	`

	expected = &kbapi.KibanaRole{}
	if err := json.Unmarshal([]byte(rawRole), expected); err != nil {
		panic(err)
	}

	// When role not exist yet
	actual = nil
	diff, err := t.kbHandler.RoleDiff(actual, expected, nil)
	if err != nil {
		t.Fail(err.Error())
	}
	assert.False(t.T(), diff.IsEmpty())
	assert.Equal(t.T(), expected, diff.Patched)

	// When role is the same
	actual = &kbapi.KibanaRole{}
	if err := json.Unmarshal([]byte(rawRole), &actual); err != nil {
		panic(err)
	}
	diff, err = t.kbHandler.RoleDiff(actual, expected, actual)
	if err != nil {
		t.Fail(err.Error())
	}
	assert.True(t.T(), diff.IsEmpty())
	assert.Equal(t.T(), expected, diff.Patched)

	// When role is not the same
	rawRole = `
{
	"metadata" : {
		"version" : 1
	},
	"elasticsearch": {
	},
	"kibana": [
		{
		"feature": {
			"discover": [
			"all"
			],
			"visualize": [
			"all"
			],
			"dashboard": [
			"read"
			],
			"dev_tools": [
			"read"
			],
			"advancedSettings": [
			"read"
			],
			"indexPatterns": [
			"read"
			],
			"graph": [
			"all"
			],
			"apm": [
			"read"
			],
			"maps": [
			"read"
			],
			"canvas": [
			"read"
			],
			"infrastructure": [
			"all"
			],
			"logs": [
			"all"
			],
			"uptime": [
			"all"
			]
		},
		"spaces": [
			"*"
		]
		}
	]
}
	`
	expected = &kbapi.KibanaRole{}
	if err := json.Unmarshal([]byte(rawRole), expected); err != nil {
		panic(err)
	}
	diff, err = t.kbHandler.RoleDiff(actual, expected, actual)
	if err != nil {
		t.Fail(err.Error())
	}
	assert.False(t.T(), diff.IsEmpty())
	assert.Equal(t.T(), expected, diff.Patched)

	// When kibana add default values
	rawRole = `
{
	"metadata" : {
		"version" : 1,
		"default" : "fake"
	},
	"elasticsearch": {
	},
	"kibana": [
		{
		"feature": {
			"discover": [
			"all"
			],
			"visualize": [
			"all"
			],
			"dashboard": [
			"all"
			],
			"dev_tools": [
			"read"
			],
			"advancedSettings": [
			"read"
			],
			"indexPatterns": [
			"read"
			],
			"graph": [
			"all"
			],
			"apm": [
			"read"
			],
			"maps": [
			"read"
			],
			"canvas": [
			"read"
			],
			"infrastructure": [
			"all"
			],
			"logs": [
			"all"
			],
			"uptime": [
			"all"
			]
		},
		"spaces": [
			"*"
		]
		}
	]
}
	`
	actual = &kbapi.KibanaRole{}
	if err := json.Unmarshal([]byte(rawRole), actual); err != nil {
		panic(err)
	}
	rawRole = `
{
	"metadata" : {
		"version" : 1
	},
	"elasticsearch": {
	},
	"kibana": [
		{
		"feature": {
			"discover": [
			"all"
			],
			"visualize": [
			"all"
			],
			"dashboard": [
			"all"
			],
			"dev_tools": [
			"read"
			],
			"advancedSettings": [
			"read"
			],
			"indexPatterns": [
			"read"
			],
			"graph": [
			"all"
			],
			"apm": [
			"read"
			],
			"maps": [
			"read"
			],
			"canvas": [
			"read"
			],
			"infrastructure": [
			"all"
			],
			"logs": [
			"all"
			],
			"uptime": [
			"all"
			]
		},
		"spaces": [
			"*"
		]
		}
	]
}
	`
	expected = &kbapi.KibanaRole{}
	if err := json.Unmarshal([]byte(rawRole), expected); err != nil {
		panic(err)
	}

	original = &kbapi.KibanaRole{}
	if err := json.Unmarshal([]byte(rawRole), original); err != nil {
		panic(err)
	}

	diff, err = t.kbHandler.RoleDiff(actual, expected, original)
	if err != nil {
		t.Fail(err.Error())
	}
	assert.True(t.T(), diff.IsEmpty())
	assert.Equal(t.T(), actual, diff.Patched)

}
