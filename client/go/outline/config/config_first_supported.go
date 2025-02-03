// Copyright 2025 The Outline Authors
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

package config

import (
	"context"
	"errors"
	"fmt"
)

type FirstSupportedConfig struct {
	Options []any
}

func parseFirstSupported[Output any](ctx context.Context, configMap map[string]any, parseE ParseFunc[Output]) (Output, error) {
	var zero Output
	var config FirstSupportedConfig
	if err := mapToAny(configMap, &config); err != nil {
		return zero, fmt.Errorf("invalid config format: %w", err)
	}

	if len(config.Options) == 0 {
		return zero, errors.New("empty list of options")
	}

	for _, ec := range config.Options {
		endpoint, err := parseE(ctx, ec)
		if errors.Is(err, errors.ErrUnsupported) {
			continue
		}
		return endpoint, err
	}
	return zero, fmt.Errorf("no suported option found: %w", errors.ErrUnsupported)
}
