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

package bufcli

import (
	"context"

	"github.com/bufbuild/buf/internal/buf/bufanalysis"
	"github.com/bufbuild/buf/internal/buf/bufbuild"
	"github.com/bufbuild/buf/internal/buf/bufconfig"
	"github.com/bufbuild/buf/internal/buf/buffetch"
	"github.com/bufbuild/buf/internal/buf/bufimage"
	"github.com/bufbuild/buf/internal/pkg/app"
	"go.uber.org/zap"
)

// Env is an environment.
type Env interface {
	Image() bufimage.Image
	Config() *bufconfig.Config
}

// EnvReader is an environment reader.
type EnvReader interface {
	// GetEnv gets an environment for the fetch value.
	//
	// If externalFilePaths is empty, this builds all files under Buf control.
	GetEnv(
		ctx context.Context,
		container app.EnvStdinContainer,
		value string,
		configOverride string,
		externalFilePaths []string,
		externalFileFilePathsAllowNotExist bool,
		excludeSourceInfo bool,
	) (Env, []bufanalysis.FileAnnotation, error)
	// GetImageEnv is the same as GetEnv but only allows image values and never builds.
	GetImageEnv(
		ctx context.Context,
		container app.EnvStdinContainer,
		value string,
		configOverride string,
		externalFilePaths []string,
		externalFileFilePathsAllowNotExist bool,
	) (Env, error)
	// GetSourceEnv is the same as GetEnv but only allows source values and always builds.
	GetSourceEnv(
		ctx context.Context,
		container app.EnvStdinContainer,
		value string,
		configOverride string,
		externalFilePaths []string,
		externalFileFilePathsAllowNotExist bool,
		excludeSourceInfo bool,
	) (Env, []bufanalysis.FileAnnotation, error)
	// ListFiles lists the files.
	ListFiles(
		ctx context.Context,
		container app.EnvStdinContainer,
		value string,
		configOverride string,
	) ([]bufimage.FileRef, error)
	// GetConfig gets the config.
	GetConfig(
		ctx context.Context,
		configOverride string,
	) (*bufconfig.Config, error)
}

// NewEnvReader returns a new EnvReader.
func NewEnvReader(
	logger *zap.Logger,
	fetchRefParser buffetch.RefParser,
	fetchReader buffetch.Reader,
	configProvider bufconfig.Provider,
	buildFileRefProvider bufbuild.FileRefProvider,
	buildBuilder bufbuild.Builder,
	valueFlagName string,
	configOverrideFlagName string,
) EnvReader {
	return newEnvReader(
		logger,
		fetchRefParser,
		fetchReader,
		configProvider,
		buildFileRefProvider,
		buildBuilder,
		valueFlagName,
		configOverrideFlagName,
	)
}

// ImageWriter is an image writer.
type ImageWriter interface {
	// PutImage writes the image to the value.
	//
	// The file must be an image format.
	// This is a no-np if value is the equivalent of /dev/null.
	PutImage(
		ctx context.Context,
		container app.EnvStdoutContainer,
		value string,
		image bufimage.Image,
		asFileDescriptorSet bool,
		excludeImports bool,
	) error
}

// NewImageWriter returns a new ImageWriter.
func NewImageWriter(
	logger *zap.Logger,
	fetchRefParser buffetch.RefParser,
	fetchWriter buffetch.Writer,
) ImageWriter {
	return newImageWriter(
		logger,
		fetchRefParser,
		fetchWriter,
	)
}
