/*
Copyright 2023 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package flux

import (
	"context"
	"fmt"

	"sigs.k8s.io/e2e-framework/pkg/env"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
)

var manager *Manager

const NoFluxInstallationFoundMsg = "flux needs to be installed within a cluster first"

// InstallFlux installs all flux components into the cluster. It is possible to specify a target namespace with flux.WithNamespace(). Default namespace is 'flux-system'
func InstallFlux(opts ...Option) env.Func {
	return func(ctx context.Context, c *envconf.Config) (context.Context, error) {
		manager = New(c.KubeconfigFile())
		err := manager.installFlux(opts...)
		if err != nil {
			return ctx, fmt.Errorf("installation of flux failed: %w", err)
		}
		return ctx, nil
	}
}

// CreateGitRepo creates a reference to a specific repository, it is a source for Kustomization or HelmRelease
func CreateGitRepo(gitRepoName, gitRepoURL string, opts ...Option) env.Func {
	return func(ctx context.Context, c *envconf.Config) (context.Context, error) {
		if manager == nil {
			return ctx, fmt.Errorf(NoFluxInstallationFoundMsg)
		}
		err := manager.createSource(Git, gitRepoName, gitRepoURL, opts...)
		if err != nil {
			return ctx, fmt.Errorf("git reporistory creation failed: %w", err)
		}
		return ctx, nil
	}
}

// CreateKustomization is used to point to a specific source and path for reconciliation
func CreateKustomization(kustomizationName, sourceRef string, opts ...Option) env.Func {
	return func(ctx context.Context, c *envconf.Config) (context.Context, error) {
		if manager == nil {
			return ctx, fmt.Errorf(NoFluxInstallationFoundMsg)
		}
		err := manager.createKustomization(kustomizationName, sourceRef, opts...)
		if err != nil {
			return ctx, fmt.Errorf("kustomization creation failed: %w", err)
		}
		return ctx, nil
	}
}

// UninstallFlux removes all flux components from a cluster
func UninstallFlux(opts ...Option) env.Func {
	return func(ctx context.Context, c *envconf.Config) (context.Context, error) {
		if manager == nil {
			return ctx, fmt.Errorf(NoFluxInstallationFoundMsg)
		}
		err := manager.uninstallFlux(opts...)
		if err != nil {
			return ctx, fmt.Errorf("uninstallation of flux failed: %w", err)
		}
		return ctx, nil
	}
}

// DeleteKustomization removes a specific Kustomization object from the cluster
func DeleteKustomization(kustomizationName string, opts ...Option) env.Func {
	return func(ctx context.Context, c *envconf.Config) (context.Context, error) {
		if manager == nil {
			return ctx, fmt.Errorf(NoFluxInstallationFoundMsg)
		}
		err := manager.deleteKustomization(kustomizationName, opts...)
		if err != nil {
			return ctx, fmt.Errorf("kustomization creation failed: %w", err)
		}
		return ctx, nil
	}
}

// DeleteGitRepo removes a specific GitRepository object from the cluster
func DeleteGitRepo(gitRepoName string, opts ...Option) env.Func {
	return func(ctx context.Context, c *envconf.Config) (context.Context, error) {
		if manager == nil {
			return ctx, fmt.Errorf(NoFluxInstallationFoundMsg)
		}
		err := manager.deleteSource(Git, gitRepoName, opts...)
		if err != nil {
			return ctx, fmt.Errorf("git reporistory deletion failed: %w", err)
		}
		return ctx, nil
	}
}
