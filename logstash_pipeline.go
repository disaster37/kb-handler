package kbhandler

import (
	"github.com/disaster37/go-kibana-rest/v8/kbapi"
	"github.com/disaster37/kb-handler/v8/patch"
	jsonIterator "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

// LogstashPipelineUpdate permit to create or update Logstash pipeline
func (h *KibanaHandlerImpl) LogstashPipelineUpdate(pipeline *kbapi.LogstashPipeline) (err error) {
	h.log.Debugf("Update Logstash pipeline %s", pipeline.ID)

	_, err = h.client.KibanaLogstashPipeline.CreateOrUpdate(pipeline)
	return err
}

// LogstashPipelineDelete permit to delete Logstash pipeline
func (h *KibanaHandlerImpl) LogstashPipelineDelete(name string) (err error) {
	h.log.Debugf("Delete Logstash pipeline %s", name)

	return h.client.KibanaLogstashPipeline.Delete(name)
}

// LogstashPipelineGet permit to get Logstash pipeline
func (h *KibanaHandlerImpl) LogstashPipelineGet(name string) (pipeline *kbapi.LogstashPipeline, err error) {
	h.log.Debugf("Get Logstash pipeline %s", name)

	return h.client.KibanaLogstashPipeline.Get(name)
}

// LogstashPipelineDiff permit to diff Logstash pipeline
func (h *KibanaHandlerImpl) LogstashPipelineDiff(actualObject, expectedObject, originalObject *kbapi.LogstashPipeline) (patchResult *patch.PatchResult, err error) {
	// If not yet exist
	if actualObject == nil {
		expected, err := jsonIterator.ConfigCompatibleWithStandardLibrary.Marshal(expectedObject)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to convert expected object to byte sequence")
		}

		return &patch.PatchResult{
			Patch:    expected,
			Current:  expected,
			Modified: expected,
			Original: nil,
			Patched:  expectedObject,
		}, nil
	}

	return patch.DefaultPatchMaker.Calculate(actualObject, expectedObject, originalObject)
}
