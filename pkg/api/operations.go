package api

var AllOperations = OperationList{
	Operations: []Operation{
		OperationResultsRead,
		OperationStatusRead,
		OperationRead,
		OperationOpenShiftClusterRead,
		OperationOpenShiftClusterWrite,
		OperationOpenShiftClusterDelete,
		OperationOpenShiftClusterListCredentials,
		OperationOpenShiftClusterListAdminCredentials,
		OperationListInstallVersions,
		OperationSyncSetsRead,
		OperationSyncSetsWrite,
		OperationSyncSetsDelete,
		OperationMachinePoolsRead,
		OperationMachinePoolsWrite,
		OperationMachinePoolsDelete,
		OperationSyncIdentityProvidersRead,
		OperationSyncIdentityProvidersWrite,
		OperationSyncIdentityProvidersDelete,
		OperationOpenShiftClusterGetDetectors,
	},
}
