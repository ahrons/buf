// Copyright 2020 Buf Technologies Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bufconfig

import (
	"context"
	"io/ioutil"

	"github.com/bufbuild/buf/internal/buf/bufcheck/bufbreaking"
	"github.com/bufbuild/buf/internal/buf/bufcheck/buflint"
	"github.com/bufbuild/buf/internal/pkg/encoding"
	"github.com/bufbuild/buf/internal/pkg/instrument"
	"github.com/bufbuild/buf/internal/pkg/storage"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

type provider struct {
	logger                 *zap.Logger
	externalConfigModifier func(*ExternalConfig) error
}

func newProvider(logger *zap.Logger, options ...ProviderOption) *provider {
	provider := &provider{
		logger: logger,
	}
	for _, option := range options {
		option(provider)
	}
	return provider
}

func (p *provider) GetConfig(ctx context.Context, readBucket storage.ReadBucket) (_ *Config, retErr error) {
	defer instrument.Start(p.logger, "get_config").End()

	externalConfig := &ExternalConfig{}
	readObject, err := readBucket.Get(ctx, ConfigFilePath)
	if err != nil {
		if storage.IsNotExist(err) {
			return p.newConfig(externalConfig)
		}
		return nil, err
	}
	defer func() {
		retErr = multierr.Append(retErr, readObject.Close())
	}()
	data, err := ioutil.ReadAll(readObject)
	if err != nil {
		return nil, err
	}
	if err := encoding.UnmarshalYAMLStrict(data, externalConfig); err != nil {
		return nil, err
	}
	return p.newConfig(externalConfig)
}

func (p *provider) GetConfigForData(data []byte) (*Config, error) {
	defer instrument.Start(p.logger, "get_config_for_data").End()

	externalConfig := &ExternalConfig{}
	if err := encoding.UnmarshalJSONOrYAMLStrict(data, externalConfig); err != nil {
		return nil, err
	}
	return p.newConfig(externalConfig)
}

func (p *provider) newConfig(externalConfig *ExternalConfig) (*Config, error) {
	if p.externalConfigModifier != nil {
		if err := p.externalConfigModifier(externalConfig); err != nil {
			return nil, err
		}
	}
	breakingConfig, err := bufbreaking.NewConfig(externalConfig.Breaking)
	if err != nil {
		return nil, err
	}
	lintConfig, err := buflint.NewConfig(externalConfig.Lint)
	if err != nil {
		return nil, err
	}
	return &Config{
		Roots:    externalConfig.Build.Roots,
		Excludes: externalConfig.Build.Excludes,
		Breaking: breakingConfig,
		Lint:     lintConfig,
	}, nil
}
