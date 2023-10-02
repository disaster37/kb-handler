package kbhandler

import (
	"github.com/disaster37/go-kibana-rest/v8/kbapi"
	"github.com/disaster37/generic-objectmatcher/patch"
	jsonIterator "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

// RoleUpdate permit to update or create role
func (h *KibanaHandlerImpl) RoleUpdate(role *kbapi.KibanaRole) (err error) {
	h.log.Debugf("Update role %s", role.Name)

	_, err = h.client.KibanaRoleManagement.CreateOrUpdate(role)
	return err
}

// RoleDelete permit to delete role
func (h *KibanaHandlerImpl) RoleDelete(name string) (err error) {
	h.log.Debugf("Delete role %s", name)

	return h.client.KibanaRoleManagement.Delete(name)
}

// RoleGet permit to get a role
func (h *KibanaHandlerImpl) RoleGet(name string) (role *kbapi.KibanaRole, err error) {
	h.log.Debugf("Role name: %s", name)

	return h.client.KibanaRoleManagement.Get(name)
}

// RoleDiff permit to diff role
func (h *KibanaHandlerImpl) RoleDiff(actualObject, expectedObject, originalObject *kbapi.KibanaRole) (patchResult *patch.PatchResult, err error) {
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
