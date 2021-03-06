package endpoint

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"

	"github.com/vterdunov/janna-api/internal/service"
)

// Endpoints collects all of the endpoints that compose the Service.
type Endpoints struct {
	InfoEndpoint endpoint.Endpoint

	HealthzEndpoint endpoint.Endpoint
	ReadyzEndpoint  endpoint.Endpoint

	VMListEndpoint   endpoint.Endpoint
	VMInfoEndpoint   endpoint.Endpoint
	VMDeleteEndpoint endpoint.Endpoint
	VMFindEndpoint   endpoint.Endpoint

	VMDeployEndpoint endpoint.Endpoint

	VMSnapshotsListEndpoint       endpoint.Endpoint
	VMSnapshotCreateEndpoint      endpoint.Endpoint
	VMSnapshotDeleteEndpoint      endpoint.Endpoint
	VMRestoreFromSnapshotEndpoint endpoint.Endpoint

	VMPowerEndpoint endpoint.Endpoint

	VMRolesListEndpoint endpoint.Endpoint
	VMAddRoleEndpoint   endpoint.Endpoint

	VMScreenshotEndpoint endpoint.Endpoint

	VMRenameEndpoint endpoint.Endpoint

	RoleListEndpoint endpoint.Endpoint

	TaskInfoEndpoint endpoint.Endpoint

	OpenAPIEndpoint endpoint.Endpoint
}

// New returns an Endpoints struct where each endpoint invokes
// the corresponding method on the provided service.
func New(s service.Service, logger log.Logger) Endpoints {
	infoEndpoint := MakeInfoEndpoint(s)
	infoEndpoint = LoggingMiddleware(log.With(logger, "endpoint", "Info"))(infoEndpoint)

	healthzEndpoint := MakeHealthzEndpoint(s)

	readyzEndpoint := MakeReadyzEndpoint(s)

	vmListEndpoint := MakeVMListEndpoint(s)
	vmListEndpoint = LoggingMiddleware(log.With(logger, "endpoint", "VMList"))(vmListEndpoint)

	vmInfoEndpoint := MakeVMInfoEndpoint(s)
	vmInfoEndpoint = LoggingMiddleware(log.With(logger, "endpoint", "VMInfo"))(vmInfoEndpoint)

	vmDeleteEndpoint := MakeVMDeleteEndpoint(s)
	vmDeleteEndpoint = LoggingMiddleware(log.With(logger, "endpoint", "VMDelete"))(vmDeleteEndpoint)

	vmFindEndpoint := MakeVMFindEndpoint(s)
	vmFindEndpoint = LoggingMiddleware(log.With(logger, "endpoint", "VMFind"))(vmFindEndpoint)

	vmDeployEndpoint := MakeVMDeployEndpoint(s, logger)
	vmDeployEndpoint = LoggingMiddleware(log.With(logger, "endpoint", "VMDeploy"))(vmDeployEndpoint)

	vmSnapshotsListEndpoint := MakeVMSnapshotsListEndpoint(s)
	vmSnapshotsListEndpoint = LoggingMiddleware(log.With(logger, "endpoint", "VMSnapshotsList"))(vmSnapshotsListEndpoint)

	vmSnapshotCreateEndpoint := MakeVMSnapshotCreateEndpoint(s)
	vmSnapshotCreateEndpoint = LoggingMiddleware(log.With(logger, "endpoint", "VMSnapshotCreate"))(vmSnapshotCreateEndpoint)

	vmRestoreFromSnapshotEndpoint := MakeVMRestoreFromSnapshotEndpoint(s)
	vmRestoreFromSnapshotEndpoint = LoggingMiddleware(log.With(logger, "endpoint", "VMRestoreFromSnapshot"))(vmRestoreFromSnapshotEndpoint)

	vmSnapshotDeleteEndpoint := MakeVMSnapshotDeleteEndpoint(s)
	vmSnapshotDeleteEndpoint = LoggingMiddleware(log.With(logger, "endpoint", "VMSnapshotDelete"))(vmSnapshotDeleteEndpoint)

	vmPowerEndpoint := MakeVMPowerEndpoint(s)
	vmPowerEndpoint = LoggingMiddleware(log.With(logger, "endpoint", "VMPower"))(vmPowerEndpoint)

	vmRolesListEndpoint := MakeVMRolesListEndpoint(s)
	vmRolesListEndpoint = LoggingMiddleware(log.With(logger, "endpoint", "VMRolesListEndpoint"))(vmRolesListEndpoint)

	vmAddROleEndpoint := MakeVMAddRoleEndpoint(s)
	vmAddROleEndpoint = LoggingMiddleware(log.With(logger, "endpoint", "VMAddRoleEndpoint"))(vmAddROleEndpoint)

	vmScreenshotEndpoint := MakeVMScreenshotEndpoint(s)
	vmScreenshotEndpoint = LoggingMiddleware(log.With(logger, "endpoint", "VMScreenshotEndpoint"))(vmScreenshotEndpoint)

	vmRenameEndpoint := MakeVMRenameEndpoint(s)
	vmRenameEndpoint = LoggingMiddleware(log.With(logger, "endpoint", "VMRenameEndpoint"))(vmRenameEndpoint)

	roleListEndpoint := MakeRolesListEndpoint(s)
	roleListEndpoint = LoggingMiddleware(log.With(logger, "endpoint", "RoleListEndpoint"))(roleListEndpoint)

	taskInfoEndpoint := MakeTaskInfoEndpoint(s)
	taskInfoEndpoint = LoggingMiddleware(log.With(logger, "endpoint", "TaskInfoEndpoint"))(taskInfoEndpoint)

	openAPIEndpoint := MakeOpenAPIEndpoint(s)
	openAPIEndpoint = LoggingMiddleware(log.With(logger, "endpoint", "OpenAPIEndpoint"))(openAPIEndpoint)

	return Endpoints{
		InfoEndpoint: infoEndpoint,

		HealthzEndpoint: healthzEndpoint,
		ReadyzEndpoint:  readyzEndpoint,

		VMListEndpoint:   vmListEndpoint,
		VMInfoEndpoint:   vmInfoEndpoint,
		VMDeleteEndpoint: vmDeleteEndpoint,
		VMFindEndpoint:   vmFindEndpoint,

		VMDeployEndpoint: vmDeployEndpoint,

		VMSnapshotsListEndpoint:       vmSnapshotsListEndpoint,
		VMSnapshotCreateEndpoint:      vmSnapshotCreateEndpoint,
		VMSnapshotDeleteEndpoint:      vmSnapshotDeleteEndpoint,
		VMRestoreFromSnapshotEndpoint: vmRestoreFromSnapshotEndpoint,

		VMPowerEndpoint: vmPowerEndpoint,

		VMRolesListEndpoint: vmRolesListEndpoint,
		VMAddRoleEndpoint:   vmAddROleEndpoint,

		VMScreenshotEndpoint: vmScreenshotEndpoint,

		VMRenameEndpoint: vmRenameEndpoint,

		RoleListEndpoint: roleListEndpoint,

		TaskInfoEndpoint: taskInfoEndpoint,

		OpenAPIEndpoint: openAPIEndpoint,
	}
}

// Failer is an interface that should be implemented by response types.
// Response encoders can check if responses are Failer, and if so they've
// failed, and if so encode them using a separate write path based on the error.
type Failer interface {
	Failed() error
}
