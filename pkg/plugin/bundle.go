/*
Copyright 2022 The Kubernetes Authors.

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

package plugin

import (
	"fmt"

	"sigs.k8s.io/kubebuilder/v4/pkg/config"
)

type bundle struct {
	name    string
	version Version
	plugins []Plugin

	supportedProjectVersions []config.Version
	deprecateWarning         string
}

// BundleOption define the options to create the bundle
type BundleOption func(*bundle)

// WithName allow set the name of the Bundle Plugin
func WithName(name string) BundleOption {
	return func(opts *bundle) {
		opts.name = name
	}
}

// WithVersion allow set the version of the Bundle Plugin
func WithVersion(version Version) BundleOption {
	return func(opts *bundle) {
		opts.version = version
	}
}

// WithPlugins allow set the plugins which will be used in the composition for the Bundle Plugin
func WithPlugins(plugins ...Plugin) BundleOption {
	return func(opts *bundle) {
		opts.plugins = plugins
	}
}

// WithDeprecationMessage allow set a deprecate message when needed
func WithDeprecationMessage(msg string) BundleOption {
	return func(opts *bundle) {
		opts.deprecateWarning = msg
	}
}

// NewBundleWithOptions creates a new Bundle with the provided BundleOptions.
// The list of supported project versions is computed from the provided plugins in options.
func NewBundleWithOptions(opts ...BundleOption) (Bundle, error) {
	bundleOpts := bundle{}

	for _, opts := range opts {
		opts(&bundleOpts)
	}

	supportedProjectVersions := CommonSupportedProjectVersions(bundleOpts.plugins...)
	if len(supportedProjectVersions) == 0 {
		return nil, fmt.Errorf("in order to bundle plugins, they must all support at least one common project version")
	}

	// Plugins may be bundles themselves, so unbundle here
	// NOTE(Adirio): unbundling here ensures that Bundle.Plugin always returns a flat list of Plugins instead of also
	//               including Bundles, and therefore we don't have to use a recursive algorithm when resolving.
	allPlugins := make([]Plugin, 0, len(bundleOpts.plugins))
	for _, plugin := range bundleOpts.plugins {
		if pluginBundle, isBundle := plugin.(Bundle); isBundle {
			allPlugins = append(allPlugins, pluginBundle.Plugins()...)
		} else {
			allPlugins = append(allPlugins, plugin)
		}
	}

	return bundle{
		name:                     bundleOpts.name,
		version:                  bundleOpts.version,
		plugins:                  allPlugins,
		supportedProjectVersions: supportedProjectVersions,
		deprecateWarning:         bundleOpts.deprecateWarning,
	}, nil
}

// Name implements Plugin
func (b bundle) Name() string {
	return b.name
}

// Version implements Plugin
func (b bundle) Version() Version {
	return b.version
}

// SupportedProjectVersions implements Plugin
func (b bundle) SupportedProjectVersions() []config.Version {
	return b.supportedProjectVersions
}

// Plugins implements Bundle
func (b bundle) Plugins() []Plugin {
	return b.plugins
}

// DeprecationWarning return the warning message
func (b bundle) DeprecationWarning() string {
	return b.deprecateWarning
}
