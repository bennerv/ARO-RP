package azureclient

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/go-autorest/autorest/azure"

	"github.com/Azure/ARO-RP/pkg/util/azureclient/azuresdk/common"
	utilgraph "github.com/Azure/ARO-RP/pkg/util/graph"
)

// AROEnvironment contains additional, cloud-specific information needed by ARO.
type AROEnvironment struct {
	azure.Environment
	ActualCloudName          string
	GenevaMonitoringEndpoint string
	AppSuffix                string
	AppLensEndpoint          string
	AppLensScope             string
	AppLensTenantID          string
	PkiIssuerUrlTemplate     string
	PkiCaName                string
	AuthzRemotePDPEndPoint   string
	AzureRbacPDPEnvironment
	Cloud cloud.Configuration
	// Microsoft identity platform scopes used by ARO
	// See https://learn.microsoft.com/EN-US/azure/active-directory/develop/scopes-oidc#the-default-scope
	ResourceManagerScope   string
	KeyVaultScope          string
	MicrosoftGraphScope    string
	CosmosDBDNSSuffixScope string
}

// AzureRbacPDPEnvironment contains cloud specific instance of Authz RBAC PDP Remote Server
type AzureRbacPDPEnvironment struct {
	Endpoint   string
	OAuthScope string
}

var (
	// PublicCloud contains additional ARO information for the public Azure cloud environment.
	PublicCloud = AROEnvironment{
		Environment:              azure.PublicCloud,
		ActualCloudName:          "AzureCloud",
		GenevaMonitoringEndpoint: "https://gcs.prod.monitoring.core.windows.net/",
		AppSuffix:                "aro.azure.com",
		AppLensEndpoint:          "https://diag-runtimehost-prod.trafficmanager.net/api/invoke",
		AppLensScope:             "0d7b6142-46a3-426a-ad6d-eed97c2a48ee",
		AppLensTenantID:          "33e01921-4d64-4f8c-a055-5bdaffd5e33d",
		PkiIssuerUrlTemplate:     "https://issuer.pki.azure.com/dsms/issuercertificates?getissuersv3&caName=%s",
		PkiCaName:                "ame",
		Cloud:                    cloud.AzurePublic,
		AzureRbacPDPEnvironment: AzureRbacPDPEnvironment{
			Endpoint:   "https://%s.authorization.azure.net/providers/Microsoft.Authorization/checkAccess?api-version=2021-06-01-preview",
			OAuthScope: "https://authorization.azure.net/.default",
		},
		ResourceManagerScope:   azure.PublicCloud.ResourceManagerEndpoint + "/.default",
		KeyVaultScope:          azure.PublicCloud.ResourceIdentifiers.KeyVault + "/.default",
		MicrosoftGraphScope:    azure.PublicCloud.MicrosoftGraphEndpoint + "/.default",
		CosmosDBDNSSuffixScope: azure.PublicCloud.CosmosDBDNSSuffix + "/.default",
	}

	// USGovernmentCloud contains additional ARO information for the US Gov cloud environment.
	USGovernmentCloud = AROEnvironment{
		Environment:              azure.USGovernmentCloud,
		ActualCloudName:          "AzureUSGovernment",
		GenevaMonitoringEndpoint: "https://gcs.monitoring.core.usgovcloudapi.net/",
		AppSuffix:                "aro.azure.us",
		AppLensEndpoint:          "https://diag-runtimehost-prod-bn1-001.azurewebsites.us/api/invoke",
		AppLensScope:             "https://microsoft.onmicrosoft.com/runtimehost",
		AppLensTenantID:          "cab8a31a-1906-4287-a0d8-4eef66b95f6e",
		Cloud:                    cloud.AzureGovernment,
		// the .us tls cert is issued by DigiCerts, and no additional certs are needed from pki
		PkiIssuerUrlTemplate: "",
		PkiCaName:            "",
		AzureRbacPDPEnvironment: AzureRbacPDPEnvironment{
			Endpoint:   "https://%s.authorization.azure.us/providers/Microsoft.Authorization/checkAccess?api-version=2021-06-01-preview",
			OAuthScope: "https://authorization.azure.us/.default",
		},
		ResourceManagerScope:   azure.USGovernmentCloud.ResourceManagerEndpoint + "/.default",
		KeyVaultScope:          azure.USGovernmentCloud.ResourceIdentifiers.KeyVault + "/.default",
		MicrosoftGraphScope:    azure.USGovernmentCloud.MicrosoftGraphEndpoint + "/.default",
		CosmosDBDNSSuffixScope: azure.USGovernmentCloud.CosmosDBDNSSuffix + "/.default",
	}
)

// EnvironmentFromName returns the AROEnvironment corresponding to the common name specified.
func EnvironmentFromName(name string) (AROEnvironment, error) {
	switch strings.ToUpper(name) {
	case "AZUREPUBLICCLOUD":
		return PublicCloud, nil
	case "AZUREUSGOVERNMENTCLOUD":
		return USGovernmentCloud, nil
	}
	return AROEnvironment{}, fmt.Errorf("cloud environment %q is unsupported by ARO", name)
}

// ArmClientOptions returns an arm.ClientOptions to be passed in when instantiating
// Azure SDK for Go clients.
func (e *AROEnvironment) ArmClientOptions(middlewares ...policy.Policy) *arm.ClientOptions {
	return &arm.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Cloud:           e.Cloud,
			Retry:           common.RetryOptions,
			PerCallPolicies: append([]policy.Policy{NewLoggingPolicy()}, middlewares...),
		},
	}
}

func (e *AROEnvironment) ClientCertificateCredentialOptions(additionalTenants []string) *azidentity.ClientCertificateCredentialOptions {
	return &azidentity.ClientCertificateCredentialOptions{
		AdditionallyAllowedTenants: additionalTenants,
		ClientOptions:              e.AzureClientOptions(),
		// Required for Subject Name/Issuer (SNI) authentication
		SendCertificateChain: true,
	}
}

func (e *AROEnvironment) ClientSecretCredentialOptions() *azidentity.ClientSecretCredentialOptions {
	return &azidentity.ClientSecretCredentialOptions{
		ClientOptions: e.AzureClientOptions(),
	}
}

func (e *AROEnvironment) DefaultAzureCredentialOptions() *azidentity.DefaultAzureCredentialOptions {
	return &azidentity.DefaultAzureCredentialOptions{
		ClientOptions: e.AzureClientOptions(),
	}
}

func (e *AROEnvironment) EnvironmentCredentialOptions() *azidentity.EnvironmentCredentialOptions {
	return &azidentity.EnvironmentCredentialOptions{
		ClientOptions: e.AzureClientOptions(),
	}
}

func (e *AROEnvironment) ManagedIdentityCredentialOptions() *azidentity.ManagedIdentityCredentialOptions {
	return &azidentity.ManagedIdentityCredentialOptions{
		ClientOptions: e.AzureClientOptions(),
	}
}

func (e *AROEnvironment) WorkloadIdentityCredentialOptions() *azidentity.WorkloadIdentityCredentialOptions {
	return &azidentity.WorkloadIdentityCredentialOptions{
		ClientOptions: e.AzureClientOptions(),
	}
}

func (e *AROEnvironment) AzureClientOptions() azcore.ClientOptions {
	return azcore.ClientOptions{
		Cloud: e.Cloud,
	}
}

func (e *AROEnvironment) NewGraphServiceClient(tokenCredential azcore.TokenCredential) (*utilgraph.GraphServiceClient, error) {
	scopes := []string{e.MicrosoftGraphScope}
	client, err := utilgraph.NewGraphServiceClientWithCredentials(tokenCredential, scopes)
	if err != nil {
		return nil, err
	}

	// Watch this issue for a better way to configure the graph endpoint.
	// https://github.com/microsoftgraph/msgraph-sdk-go/issues/235
	client.GetAdapter().SetBaseUrl(e.MicrosoftGraphEndpoint + "v1.0")

	return client, nil
}
