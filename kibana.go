package kbhandler

import (
	"github.com/disaster37/go-kibana-rest/v8"
	"github.com/disaster37/go-kibana-rest/v8/kbapi"
	"github.com/disaster37/generic-objectmatcher/patch"
	"github.com/sirupsen/logrus"
)

type KibanaHandler interface {
	Client() (client *kibana.Client)
	SetLogger(log *logrus.Entry)

	// User space scope
	UserSpaceCreate(kibanaSpace *kbapi.KibanaSpace) (err error)
	UserSpaceUpdate(kibanaSpace *kbapi.KibanaSpace) (err error)
	UserSpaceDelete(name string) (err error)
	UserSpaceGet(name string) (userspace *kbapi.KibanaSpace, err error)
	UserSpaceDiff(actualObject, expectedObject, originalObject *kbapi.KibanaSpace) (patchResult *patch.PatchResult, err error)
	UserSpaceCopyObject(userSpaceOrigin string, copySpec *kbapi.KibanaSpaceCopySavedObjectParameter) (err error)

	// Role scope
	RoleUpdate(role *kbapi.KibanaRole) (err error)
	RoleDelete(name string) (err error)
	RoleGet(name string) (role *kbapi.KibanaRole, err error)
	RoleDiff(actualObject, expectedObject, originalObject *kbapi.KibanaRole) (patchResult *patch.PatchResult, err error)

	// Logstash pipeline scope
	LogstashPipelineUpdate(pipeline *kbapi.LogstashPipeline) (err error)
	LogstashPipelineDelete(name string) (err error)
	LogstashPipelineGet(name string) (pipeline *kbapi.LogstashPipeline, err error)
	LogstashPipelineDiff(actualObject, expectedObject, originalObject *kbapi.LogstashPipeline) (patchResult *patch.PatchResult, err error)
}

type KibanaHandlerImpl struct {
	client *kibana.Client
	log    *logrus.Entry
}

func NewKibanaHandler(cfg kibana.Config, log *logrus.Entry) (KibanaHandler, error) {

	client, err := kibana.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return &KibanaHandlerImpl{
		client: client,
		log:    log,
	}, nil
}

func (h *KibanaHandlerImpl) SetLogger(log *logrus.Entry) {
	h.log = log
}

func (h *KibanaHandlerImpl) Client() *kibana.Client {
	return h.client
}
