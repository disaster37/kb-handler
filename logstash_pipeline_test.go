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

var urlLogstashPipeline = fmt.Sprintf("%s/api/logstash/pipeline/test", baseURL)

func (t *KibanaHandlerTestSuite) TestLogstashPipelineGet() {

	rawPipeline := `
{
	"id": "test",
	"description": "Just a simple pipeline",
	"username": "elastic",
	"pipeline": "input { stdin {} } output { stdout {} }",
	"settings": {
		"queue.type": "persistent"
	}
}
	`

	pipelineTest := &kbapi.LogstashPipeline{}
	if err := json.Unmarshal([]byte(rawPipeline), pipelineTest); err != nil {
		panic(err)
	}

	httpmock.RegisterResponder("GET", urlLogstashPipeline, func(req *http.Request) (*http.Response, error) {
		resp := httpmock.NewStringResponse(200, rawPipeline)
		return resp, nil
	})

	pipeline, err := t.kbHandler.LogstashPipelineGet("test")
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Empty(t.T(), cmp.Diff(pipelineTest, pipeline))

	// When error
	httpmock.RegisterResponder("GET", urlLogstashPipeline, httpmock.NewErrorResponder(errors.New("fack error")))
	_, err = t.kbHandler.LogstashPipelineGet("test")
	assert.Error(t.T(), err)
}

func (t *KibanaHandlerTestSuite) TestLogstashPipelineDelete() {

	httpmock.RegisterResponder("DELETE", urlLogstashPipeline, func(req *http.Request) (*http.Response, error) {
		resp := httpmock.NewStringResponse(200, "")
		return resp, nil
	})

	err := t.kbHandler.LogstashPipelineDelete("test")
	if err != nil {
		t.Fail(err.Error())
	}

	// When error
	httpmock.RegisterResponder("DELETE", urlLogstashPipeline, httpmock.NewErrorResponder(errors.New("fack error")))
	err = t.kbHandler.LogstashPipelineDelete("test")
	assert.Error(t.T(), err)
}

func (t *KibanaHandlerTestSuite) TestLogstashPipelineUpdate() {

	rawPipeline := `
{
	"id": "test",
	"pipeline": "input { stdin {} } output { stdout {} }",
	"settings": {
		"queue.type": "persisted"
	}
}
	`

	pipeline := &kbapi.LogstashPipeline{}
	if err := json.Unmarshal([]byte(rawPipeline), pipeline); err != nil {
		panic(err)
	}

	httpmock.RegisterResponder("PUT", urlLogstashPipeline, func(req *http.Request) (*http.Response, error) {
		resp := httpmock.NewStringResponse(200, rawPipeline)
		return resp, nil
	})
	httpmock.RegisterResponder("GET", urlLogstashPipeline, func(req *http.Request) (*http.Response, error) {
		resp := httpmock.NewStringResponse(200, rawPipeline)
		return resp, nil
	})

	err := t.kbHandler.LogstashPipelineUpdate(pipeline)
	if err != nil {
		t.Fail(err.Error())
	}

	// When error
	httpmock.RegisterResponder("PUT", urlLogstashPipeline, httpmock.NewErrorResponder(errors.New("fack error")))
	err = t.kbHandler.LogstashPipelineUpdate(pipeline)
	assert.Error(t.T(), err)
}

func (t *KibanaHandlerTestSuite) TestLogstashPipelineDiff() {
	var actual, expected, original *kbapi.LogstashPipeline

	rawPipeline := `
{
	"id": "test",
	"description": "Just a simple pipeline",
	"username": "elastic",
	"pipeline": "input { stdin {} } output { stdout {} }",
	"settings": {
		"queue.type": "persistent"
	}
}
	`

	expected = &kbapi.LogstashPipeline{}
	if err := json.Unmarshal([]byte(rawPipeline), expected); err != nil {
		panic(err)
	}

	// When pipeline not exist yet
	actual = nil
	diff, err := t.kbHandler.LogstashPipelineDiff(actual, expected, nil)
	if err != nil {
		t.Fail(err.Error())
	}
	assert.False(t.T(), diff.IsEmpty())
	assert.Equal(t.T(), expected, diff.Patched)

	// When pipeline is the same
	actual = &kbapi.LogstashPipeline{}
	if err := json.Unmarshal([]byte(rawPipeline), &actual); err != nil {
		panic(err)
	}
	diff, err = t.kbHandler.LogstashPipelineDiff(actual, expected, actual)
	if err != nil {
		t.Fail(err.Error())
	}
	assert.True(t.T(), diff.IsEmpty())
	assert.Equal(t.T(), expected, diff.Patched)

	// When pileine is not the same
	rawPipeline = `
{
	"id": "test",
	"description": "Just a simple pipeline",
	"username": "elastic",
	"pipeline": "input { stdin {} } output { stdout {} }",
	"settings": {
		"queue.type": "memory"
	}
}
	`
	expected = &kbapi.LogstashPipeline{}
	if err := json.Unmarshal([]byte(rawPipeline), expected); err != nil {
		panic(err)
	}
	diff, err = t.kbHandler.LogstashPipelineDiff(actual, expected, actual)
	if err != nil {
		t.Fail(err.Error())
	}
	assert.False(t.T(), diff.IsEmpty())
	assert.Equal(t.T(), expected, diff.Patched)

	// When kibana add default values
	rawPipeline = `
{
	"id": "test",
	"description": "Just a simple pipeline",
	"username": "elastic",
	"pipeline": "input { stdin {} } output { stdout {} }",
	"settings": {
		"queue.type": "persistent",
		"default": "fake"
	}
}
	`
	actual = &kbapi.LogstashPipeline{}
	if err := json.Unmarshal([]byte(rawPipeline), actual); err != nil {
		panic(err)
	}
	rawPipeline = `
{
	"id": "test",
	"description": "Just a simple pipeline",
	"username": "elastic",
	"pipeline": "input { stdin {} } output { stdout {} }",
	"settings": {
		"queue.type": "persistent"
	}
}
	`
	expected = &kbapi.LogstashPipeline{}
	if err := json.Unmarshal([]byte(rawPipeline), expected); err != nil {
		panic(err)
	}

	original = &kbapi.LogstashPipeline{}
	if err := json.Unmarshal([]byte(rawPipeline), original); err != nil {
		panic(err)
	}

	diff, err = t.kbHandler.LogstashPipelineDiff(actual, expected, original)
	if err != nil {
		t.Fail(err.Error())
	}
	assert.True(t.T(), diff.IsEmpty())
	assert.Equal(t.T(), actual, diff.Patched)

}
