package main

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"context"
	"flag"
	"strings"

	mgmtfeatures "github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-07-01/features"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"
)

func nic(ctx context.Context, log *logrus.Entry) error {
	nicName := flag.Arg(1)
	nic, err := azure.ParseResourceID(nicName)
	if err != nil {
		return err
	}

	authorizer, err := auth.NewAuthorizerFromCLIWithResource("https://management.azure.com")
	if err != nil {
		return err
	}

	client := mgmtfeatures.NewResourcesClient(nic.SubscriptionID)
	client.Authorizer = authorizer

	log.Info("listing all resources in resource group")
	var resources []mgmtfeatures.GenericResourceExpanded
	page, err := client.ListByResourceGroup(ctx, nic.ResourceGroup, "", "provisioningState", nil)
	if err != nil {
		return err
	}

	for page.NotDone() {
		resources = append(resources, page.Values()...)
		err = page.NextWithContext(ctx)
		if err != nil {
			return err
		}
	}

	for _, resource := range resources {
		switch strings.ToLower(*resource.Type) {
		case "microsoft.network/networkinterfaces":
			spew.Dump(resource)
		}
	}

	return nil
}
