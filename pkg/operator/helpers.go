package operator

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"context"

	configv1 "github.com/openshift/api/config/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	arov1alpha1 "github.com/Azure/ARO-RP/pkg/operator/apis/aro.openshift.io/v1alpha1"
	"github.com/Azure/ARO-RP/pkg/util/version"
)

func GatewayEnabled(cluster *arov1alpha1.Cluster) bool {
	return len(cluster.Spec.GatewayDomains) > 0
}

// PodSecurityAdmissionControl is an admissions controller
// for pods which replaces pod security policies, enabled on
// OpenShift 4.11 and up
func PodSecurityAdmissionControl(ctx context.Context, client client.Client) (bool, error) {
	cv := &configv1.ClusterVersion{}

	err := client.Get(ctx, types.NamespacedName{Name: "version"}, cv)
	if err != nil {
		return false, err
	}

	vers, err := version.ParseVersion(cv.Spec.DesiredUpdate.Version)
	if err != nil {
		return false, err
	}

	if vers.Lt(version.NewVersion(4, 11)) {
		return false, nil
	}

	return true, nil
}
