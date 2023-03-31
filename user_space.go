package kbhandler

import (
	"github.com/disaster37/go-kibana-rest/v8/kbapi"
	"github.com/disaster37/kb-handler/v8/patch"
	jsonIterator "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

// UserSpaceCreate permit to create new user space
func (h *KibanaHandlerImpl) UserSpaceCreate(kibanaSpace *kbapi.KibanaSpace) (err error) {
	h.log.Debugf("Create user space %s", kibanaSpace.Name)

	_, err = h.client.KibanaSpaces.Create(kibanaSpace)
	return err
}

// UserSpaceUpdate permit to update user space
func (h *KibanaHandlerImpl) UserSpaceUpdate(kibanaSpace *kbapi.KibanaSpace) (err error) {
	h.log.Debugf("Update user space %s", kibanaSpace.Name)

	_, err = h.client.KibanaSpaces.Update(kibanaSpace)
	return err
}

// UserSpaceDelete permit to delete user space
func (h *KibanaHandlerImpl) UserSpaceDelete(name string) (err error) {
	h.log.Debugf("Name: %s", name)

	return h.client.KibanaSpaces.Delete(name)
}

// UserSpaceGet permit to get user space
func (h *KibanaHandlerImpl) UserSpaceGet(name string) (userspace *kbapi.KibanaSpace, err error) {
	h.log.Debugf("Name: %s", name)

	return h.client.KibanaSpaces.Get(name)
}

func (h *KibanaHandlerImpl) UserSpaceDiff(actualObject, expectedObject, originalObject *kbapi.KibanaSpace) (patchResult *patch.PatchResult, err error) {
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

func (h *KibanaHandlerImpl) UserSpaceCopyObject(userSpaceOrigin string, copySpec *kbapi.KibanaSpaceCopySavedObjectParameter) (err error) {
	h.log.Debugf("From User space: %s", userSpaceOrigin)

	return h.client.KibanaSpaces.CopySavedObjects(copySpec, userSpaceOrigin)
}
