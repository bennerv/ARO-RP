package main

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os/exec"

	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/creack/pty"
	"github.com/sirupsen/logrus"

	"github.com/Azure/ARO-RP/pkg/util/stringutils"
)

func expect(r io.Reader, expected []byte) error {
	b := make([]byte, len(expected))

	n, err := io.ReadAtLeast(r, b, len(b))
	if err != nil {
		return err
	}

	if !bytes.Equal(expected, b[:n]) {
		b := make([]byte, 256)
		_, err := io.ReadAtLeast(r, b, len(b))
		if err != nil {
			return err
		}
		return fmt.Errorf("expected %q, got %q", string(expected), string(b))
	}

	return nil
}

func (s *sshTool) az(ctx context.Context, log *logrus.Entry) error {
	spp := s.oc.Properties.ServicePrincipalProfile

	cmd := exec.Command("az", "login", "--service-principal", "-u", string(spp.ClientID), "-t", string(s.sub.Properties.TenantID))

	pty, err := pty.Start(cmd)
	if err != nil {
		log.Error(err)
		return err
	}

	err = expect(pty, []byte("Password: "))
	if err != nil {
		return err
	}

	_, err = pty.Write(append([]byte(spp.ClientSecret), '\n'))
	if err != nil {
		return err
	}

	err = cmd.Wait()
	if err != nil {
		return err
	}

	id, err := azure.ParseResourceID(s.oc.ID)
	if err != nil {
		return err
	}

	err = exec.Command("az", "account", "set", "-s", id.SubscriptionID).Run()
	if err != nil {
		return err
	}

	resourceGroup := stringutils.LastTokenByte(s.oc.Properties.ClusterProfile.ResourceGroupID, '/')

	return exec.Command("az", "configure", "--defaults", "group="+resourceGroup).Run()
}