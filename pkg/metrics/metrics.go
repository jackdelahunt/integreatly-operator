package metrics

import (
	"fmt"

	integreatlyv1alpha1 "github.com/integr8ly/integreatly-operator/apis/v1alpha1"
	"github.com/integr8ly/integreatly-operator/version"

	"github.com/prometheus/client_golang/prometheus"
)

// Custom metrics
var (
	OperatorVersion = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "integreatly_version_info",
			Help: "Integreatly operator information",
			ConstLabels: prometheus.Labels{
				"operator_version": version.GetVersion(),
				"version":          version.GetVersion(),
			},
		},
	)

	RHMIStatusAvailable = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "rhmi_status_available",
			Help: "RHMI status available",
		},
	)

	RHMIInfo = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rhmi_spec",
			Help: "RHMI info variables",
		},
		[]string{
			"use_cluster_storage",
			"master_url",
			"installation_type",
			"operator_name",
			"namespace",
			"namespace_prefix",
			"operators_in_product_namespace",
			"routing_subdomain",
			"self_signed_certs",
		},
	)

	RHMIVersion = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rhmi_version",
			Help: "RHMI versions",
		},
		[]string{
			"stage",
			"version",
			"to_version",
		},
	)

	RHMIStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rhmi_status",
			Help: "RHMI status of an installation",
		},
		[]string{
			"stage",
		},
	)

	RHMIPreflightStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rhmi_preflight_status",
			Help: "Preflight status of an RHMI installation",
		},
		[]string{
			"status",
		},
	)

	RHOAMVersion = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rhoam_version",
			Help: "RHOAM versions",
		},
		[]string{
			"stage",
			"version",
			"to_version",
		},
	)

	RHOAMStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rhoam_status",
			Help: "RHOAM status of an installation",
		},
		[]string{
			"stage",
		},
	)

	RHOAMPreflightStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rhoam_preflight_status",
			Help: "Preflight status of an RHOAM installation",
		},
		[]string{
			"status",
		},
	)

	ThreeScaleUserAction = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "threescale_user_action",
			Help: "Status of user CRUD action in 3scale",
		},
		[]string{
			"username",
			"action",
		},
	)

	Quota = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rhoam_quota",
			Help: "Status of the current quota config",
		},
		[]string{
			"stage",
			"quota",
			"toQuota",
		},
	)
)

// SetRHMIInfo exposes rhmi info metrics with labels from the installation CR
func SetRHMIInfo(installation *integreatlyv1alpha1.RHMI) {
	RHMIInfo.WithLabelValues(installation.Spec.UseClusterStorage,
		installation.Spec.MasterURL,
		installation.Spec.Type,
		installation.GetName(),
		installation.GetNamespace(),
		installation.Spec.NamespacePrefix,
		fmt.Sprintf("%t", installation.Spec.OperatorsInProductNamespace),
		installation.Spec.RoutingSubdomain,
		fmt.Sprintf("%t", installation.Spec.SelfSignedCerts),
	)
}

// SetRHMIStatus exposes rhmi_status metric for each stage
func SetRHMIStatus(installation *integreatlyv1alpha1.RHMI) {
	RHMIStatus.Reset()
	if string(installation.Status.Stage) != "" {
		RHMIStatus.With(prometheus.Labels{"stage": string(installation.Status.Stage)}).Set(float64(1))
	}

	RHOAMStatus.Reset()
	if string(installation.Status.Stage) != "" {
		RHOAMStatus.With(prometheus.Labels{"stage": string(installation.Status.Stage)}).Set(float64(1))
	}
}

func SetRhmiVersions(stage string, version string, toVersion string, firstInstallTimestamp int64) {
	RHMIVersion.Reset()
	RHMIVersion.WithLabelValues(stage, version, toVersion).Set(float64(firstInstallTimestamp))

	RHOAMVersion.Reset()
	RHOAMVersion.WithLabelValues(stage, version, toVersion).Set(float64(firstInstallTimestamp))
}

func SetPreflightStatus(status integreatlyv1alpha1.PreflightStatus) {
	RHOAMPreflightStatus.Reset()
	RHMIPreflightStatus.Reset()

	// -1 -> Fail
	//  0 -> In Progress
	//  1 -> Success
	preflightNumberStatus := 0
	if status == integreatlyv1alpha1.PreflightFail {
		preflightNumberStatus -= 1
	} else if status == integreatlyv1alpha1.PreflightSuccess {
		preflightNumberStatus += 1
	}

	RHOAMPreflightStatus.WithLabelValues(string(status)).Set(float64(preflightNumberStatus))
	RHMIPreflightStatus.WithLabelValues(string(status)).Set(float64(preflightNumberStatus))
}

func SetThreeScaleUserAction(httpStatus int, username, action string) {
	ThreeScaleUserAction.WithLabelValues(username, action).Set(float64(httpStatus))
}

func ResetThreeScaleUserAction() {
	ThreeScaleUserAction.Reset()
}

func SetQuota(stage string, quota string, toQuota string) {
	Quota.Reset()
	Quota.WithLabelValues(stage, quota, toQuota).Set(float64(1))
}
