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

var urlUserSpace = fmt.Sprintf("%s/api/spaces/space/test", baseURL)

func (t *KibanaHandlerTestSuite) TestUserSpaceGet() {

	rawUserSpace := `
{
	"id": "test",
	"name": "test",
	"description" : "This is the Marketing Space",
	"color": "#aabbcc",
	"initials": "MK",
	"disabledFeatures": [],
	"imageUrl": ""
}
	`

	userSpaceTest := &kbapi.KibanaSpace{}
	if err := json.Unmarshal([]byte(rawUserSpace), userSpaceTest); err != nil {
		panic(err)
	}

	httpmock.RegisterResponder("GET", urlUserSpace, func(req *http.Request) (*http.Response, error) {
		resp := httpmock.NewStringResponse(200, rawUserSpace)
		return resp, nil
	})

	userSpace, err := t.kbHandler.UserSpaceGet("test")
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Empty(t.T(), cmp.Diff(userSpaceTest, userSpace))

	// When error
	httpmock.RegisterResponder("GET", urlUserSpace, httpmock.NewErrorResponder(errors.New("fack error")))
	_, err = t.kbHandler.UserSpaceGet("test")
	assert.Error(t.T(), err)
}

func (t *KibanaHandlerTestSuite) TestUserSpaceDelete() {

	httpmock.RegisterResponder("DELETE", urlUserSpace, func(req *http.Request) (*http.Response, error) {
		resp := httpmock.NewStringResponse(200, "")
		return resp, nil
	})

	err := t.kbHandler.UserSpaceDelete("test")
	if err != nil {
		t.Fail(err.Error())
	}

	// When error
	httpmock.RegisterResponder("DELETE", urlUserSpace, httpmock.NewErrorResponder(errors.New("fack error")))
	err = t.kbHandler.UserSpaceDelete("test")
	assert.Error(t.T(), err)
}

func (t *KibanaHandlerTestSuite) TestUserSpaceCreate() {

	url := fmt.Sprintf("%s/api/spaces/space", baseURL)

	rawUserSpace := `
{
	"id": "test",
	"name": "test",
	"description" : "This is the Marketing Space",
	"color": "#aabbcc",
	"initials": "MK",
	"disabledFeatures": [],
	"imageUrl": ""
}
	`

	userSpace := &kbapi.KibanaSpace{}
	if err := json.Unmarshal([]byte(rawUserSpace), userSpace); err != nil {
		panic(err)
	}

	httpmock.RegisterResponder("POST", url, func(req *http.Request) (*http.Response, error) {
		resp := httpmock.NewStringResponse(200, rawUserSpace)
		return resp, nil
	})

	err := t.kbHandler.UserSpaceCreate(userSpace)
	if err != nil {
		t.Fail(err.Error())
	}

	// When error
	httpmock.RegisterResponder("POST", url, httpmock.NewErrorResponder(errors.New("fack error")))
	err = t.kbHandler.UserSpaceCreate(userSpace)
	assert.Error(t.T(), err)
}

func (t *KibanaHandlerTestSuite) TestUserSpaceUpdate() {

	rawUserSpace := `
{
	"id": "test",
	"name": "test",
	"description" : "This is the Marketing Space",
	"color": "#aabbcc",
	"initials": "MK",
	"disabledFeatures": [],
	"imageUrl": ""
}
	`

	userSpace := &kbapi.KibanaSpace{}
	if err := json.Unmarshal([]byte(rawUserSpace), userSpace); err != nil {
		panic(err)
	}

	httpmock.RegisterResponder("PUT", urlUserSpace, func(req *http.Request) (*http.Response, error) {
		resp := httpmock.NewStringResponse(200, rawUserSpace)
		return resp, nil
	})

	err := t.kbHandler.UserSpaceUpdate(userSpace)
	if err != nil {
		t.Fail(err.Error())
	}

	// When error
	httpmock.RegisterResponder("PUT", urlUserSpace, httpmock.NewErrorResponder(errors.New("fack error")))
	err = t.kbHandler.UserSpaceUpdate(userSpace)
	assert.Error(t.T(), err)
}

func (t *KibanaHandlerTestSuite) TestUserSpaceDiff() {
	var actual, expected, original *kbapi.KibanaSpace

	rawUserSpace := `
{
	"id": "test",
	"name": "test",
	"description" : "This is the Marketing Space",
	"color": "#aabbcc",
	"initials": "MK",
	"imageUrl": ""
}
	`

	expected = &kbapi.KibanaSpace{}
	if err := json.Unmarshal([]byte(rawUserSpace), expected); err != nil {
		panic(err)
	}

	// When User space not exist yet
	actual = nil
	diff, err := t.kbHandler.UserSpaceDiff(actual, expected, nil)
	if err != nil {
		t.Fail(err.Error())
	}
	assert.False(t.T(), diff.IsEmpty())
	assert.Equal(t.T(), expected, diff.Patched)

	// When user space is the same
	actual = &kbapi.KibanaSpace{}
	if err := json.Unmarshal([]byte(rawUserSpace), &actual); err != nil {
		panic(err)
	}
	diff, err = t.kbHandler.UserSpaceDiff(actual, expected, actual)
	if err != nil {
		t.Fail(err.Error())
	}
	assert.True(t.T(), diff.IsEmpty())
	assert.Equal(t.T(), expected, diff.Patched)

	// When user space is not the same
	rawUserSpace = `
{
	"id": "test",
	"name": "test",
	"description" : "This is the Marketing Space 2",
	"color": "#aabbcc",
	"initials": "MK",
	"imageUrl": ""
}
	`
	expected = &kbapi.KibanaSpace{}
	if err := json.Unmarshal([]byte(rawUserSpace), expected); err != nil {
		panic(err)
	}
	diff, err = t.kbHandler.UserSpaceDiff(actual, expected, actual)
	if err != nil {
		t.Fail(err.Error())
	}
	assert.False(t.T(), diff.IsEmpty())
	assert.Equal(t.T(), expected, diff.Patched)

	// When kibana add default values
	rawUserSpace = `
{
	"id": "test",
	"name": "test",
	"description" : "This is the Marketing Space",
	"color": "#aabbcc",
	"initials": "MK",
	"imageUrl": "fake"
}
	`
	actual = &kbapi.KibanaSpace{}
	if err := json.Unmarshal([]byte(rawUserSpace), actual); err != nil {
		panic(err)
	}
	rawUserSpace = `
{
	"id": "test",
	"name": "test",
	"description" : "This is the Marketing Space",
	"color": "#aabbcc",
	"initials": "MK",
	"imageUrl": ""
}
	`
	expected = &kbapi.KibanaSpace{}
	if err := json.Unmarshal([]byte(rawUserSpace), expected); err != nil {
		panic(err)
	}

	original = &kbapi.KibanaSpace{}
	if err := json.Unmarshal([]byte(rawUserSpace), original); err != nil {
		panic(err)
	}

	diff, err = t.kbHandler.UserSpaceDiff(actual, expected, original)
	if err != nil {
		t.Fail(err.Error())
	}
	assert.True(t.T(), diff.IsEmpty())
	assert.Equal(t.T(), actual, diff.Patched)

}

func (t *KibanaHandlerTestSuite) TestUserSpaceCopyObject() {

	url := "/api/spaces/_copy_saved_objects"

	httpmock.RegisterResponder("POST", url, func(req *http.Request) (*http.Response, error) {
		resp := httpmock.NewStringResponse(200, `
{
	"test": {
		"success": true,
		"successCount": 1,
		"successResults": [
			{
				"id": "fake",
				"type": "index-pattern",
				"destinationId": "bc3c9c70-bf6f-4bec-b4ce-f4189aa9e26b",
				"meta": {
					"icon": "indexPatternApp",
					"title": "my-pattern-*"
				}
			}
		]
	}
}
		`)
		return resp, nil
	})

	copySpec := &kbapi.KibanaSpaceCopySavedObjectParameter{
		Spaces:            []string{"test"},
		IncludeReferences: true,
		Overwrite:         true,
		Objects: []kbapi.KibanaSpaceObjectParameter{
			{
				Type: "index-pattern",
				ID:   "fake",
			},
		},
	}

	err := t.kbHandler.UserSpaceCopyObject("default", copySpec)
	if err != nil {
		t.Fail(err.Error())
	}

	// When error
	httpmock.RegisterResponder("POST", url, httpmock.NewErrorResponder(errors.New("fack error")))
	err = t.kbHandler.UserSpaceCopyObject("default", copySpec)
	assert.Error(t.T(), err)
}
