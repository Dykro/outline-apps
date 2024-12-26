// Copyright 2024 The Outline Authors
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
	"encoding/base64"
	"net"
	"testing"

	"github.com/Jigsaw-Code/outline-sdk/transport"
	"github.com/stretchr/testify/require"
)

func TestParseShadowsocksConfig_URL(t *testing.T) {
	t.Run("Fully Encoded", func(t *testing.T) {
		encoded := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString([]byte("chacha20-ietf-poly1305:SECRET@example.com:1234?prefix=HTTP%2F1.1%20"))
		config, err := parseShadowsocksConfig("ss://" + string(encoded) + "#outline-123")
		require.NoError(t, err)
		require.Equal(t, "example.com:1234", config.Endpoint)
		require.Equal(t, "chacha20-ietf-poly1305", config.Cipher)
		require.Equal(t, "SECRET", config.Secret)
		require.Equal(t, "HTTP/1.1 ", config.Prefix)
	})

	t.Run("User Info Encoded", func(t *testing.T) {
		encoded := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString([]byte("chacha20-ietf-poly1305:SECRET"))
		config, err := parseShadowsocksConfig("ss://" + string(encoded) + "@example.com:1234?prefix=HTTP%2F1.1%20" + "#outline-123")
		require.NoError(t, err)
		require.Equal(t, "example.com:1234", config.Endpoint)
		require.Equal(t, "chacha20-ietf-poly1305", config.Cipher)
		require.Equal(t, "SECRET", config.Secret)
		require.Equal(t, "HTTP/1.1 ", config.Prefix)
	})

	t.Run("User Info Legacy Encoded", func(t *testing.T) {
		encoded := base64.StdEncoding.EncodeToString([]byte("chacha20-ietf-poly1305:SECRET"))
		config, err := parseShadowsocksConfig("ss://" + string(encoded) + "@example.com:1234?prefix=HTTP%2F1.1%20" + "#outline-123")
		require.NoError(t, err)
		require.Equal(t, "example.com:1234", config.Endpoint)
		require.Equal(t, "chacha20-ietf-poly1305", config.Cipher)
		require.Equal(t, "SECRET", config.Secret)
		require.Equal(t, "HTTP/1.1 ", config.Prefix)
	})

	t.Run("User Info No Encoding", func(t *testing.T) {
		configString := "ss://chacha20-ietf-poly1305:SECRET@example.com:1234"
		config, err := parseShadowsocksConfig(configString)
		require.NoError(t, err)
		require.Equal(t, "example.com:1234", config.Endpoint)
		require.Equal(t, "chacha20-ietf-poly1305", config.Cipher)
		require.Equal(t, "SECRET", config.Secret)
	})

	t.Run("Invalid Cipher Fails", func(t *testing.T) {
		configString := "ss://chacha20-ietf-poly13051234567@example.com:1234"
		_, err := parseShadowsocksParams(configString)
		require.Error(t, err)
	})

	t.Run("Unsupported Cipher Fails", func(t *testing.T) {
		configString := "ss://Y2hhY2hhMjAtaWV0Zi1wb2x5MTMwnTpLeTUyN2duU3FEVFB3R0JpQ1RxUnlT@example.com:1234"
		_, err := parseShadowsocksParams(configString)
		require.Error(t, err)
	})
}

func TestNewShadowsocksTransport(t *testing.T) {
	streamEndpoints := NewTypeParser(func(ctx context.Context, config ConfigNode) (*Endpoint[transport.StreamConn], error) {
		require.Equal(t, "example.com:1234", config)
		return &Endpoint[transport.StreamConn]{}, nil
	})
	packetEndpoints := NewTypeParser(func(ctx context.Context, config ConfigNode) (*Endpoint[net.Conn], error) {
		require.Equal(t, "example.com:1234", config)
		return &Endpoint[net.Conn]{}, nil
	})

	t.Run("Success", func(t *testing.T) {
		config := map[string]any{
			"endpoint": "example.com:1234",
			"cipher":   "chacha20-ietf-poly1305",
			"secret":   "SECRET",
			"prefix":   "outline-123",
		}
		transport, err := parseShadowsocksTransport(context.Background(), config, streamEndpoints.Parse, packetEndpoints.Parse)
		require.NoError(t, err)
		require.NotNil(t, transport)
	})

	t.Run("Fail on unsupported cipher", func(t *testing.T) {
		config := map[string]any{
			"endpoint": "example.com:1234",
			"cipher":   "NOT SUPPORTED",
			"secret":   "SECRET",
			"prefix":   "outline-123",
		}
		_, err := parseShadowsocksTransport(context.Background(), config, streamEndpoints.Parse, packetEndpoints.Parse)
		require.Error(t, err)
	})

	t.Run("Fail on extraneous field", func(t *testing.T) {
		config := map[string]any{
			"endpoint": "example.com:1234",
			"cipher":   "chacha20-ietf-poly1305",
			"secret":   "SECRET",
			"prefix":   "outline-123",
			"extra":    "NOT SUPPORTED",
		}
		_, err := parseShadowsocksTransport(context.Background(), config, streamEndpoints.Parse, packetEndpoints.Parse)
		require.Error(t, err)
	})
}
