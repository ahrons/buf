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

// Package github provides TESTING functionality for GitHub.
//
// This should eventually be refactored to provide universal functionality.
package github

import (
	"context"
	"net/http"

	"github.com/bufbuild/buf/internal/pkg/app"
	"github.com/bufbuild/buf/internal/pkg/httpauth"
	"go.uber.org/zap"
)

// ArchiveReader reads GitHub archives.
type ArchiveReader interface {
	// GetArchive gets the GitHub archive and untars it to the output directory path.
	//
	// The root directory within the tarball is stripped.
	// If the directory already exists, this is a no-op.
	//
	// Only use for testing.
	GetArchive(
		ctx context.Context,
		container app.EnvStdinContainer,
		outputDirPath string,
		owner string,
		repository string,
		ref string,
	) error
}

// NewArchiveReader returns a new ArchiveReader.
func NewArchiveReader(
	logger *zap.Logger,
	httpClient *http.Client,
	httpAuthenticator httpauth.Authenticator,
) ArchiveReader {
	return newArchiveReader(
		logger,
		httpClient,
		httpAuthenticator,
	)
}
