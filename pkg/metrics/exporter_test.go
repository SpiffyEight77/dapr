/*
Copyright 2021 The Dapr Authors
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

package metrics

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dapr/dapr/pkg/healthz"
	"github.com/dapr/kit/logger"
)

func TestMetricsExporter(t *testing.T) {
	logger := logger.NewLogger("test.logger")

	t.Run("returns default options", func(t *testing.T) {
		e := New(Options{
			Enabled: DefaultFlagOptions().enabled,
			Port:    DefaultFlagOptions().port,
			Log:     logger,
			Healthz: healthz.New(),
		})
		assert.Equal(t, "9090", e.(*exporter).port)
		assert.True(t, e.(*exporter).enabled)
	})

	t.Run("skip starting metric server but wait for context cancellation", func(t *testing.T) {
		e := New(Options{
			Enabled: false,
			Port:    "9090",
			Log:     logger,
			Healthz: healthz.New(),
		})

		ctx, cancel := context.WithCancel(t.Context())
		errCh := make(chan error)
		go func() {
			errCh <- e.Start(ctx)
		}()

		cancel()

		select {
		case err := <-errCh:
			require.NoError(t, err)
		case <-time.After(time.Second):
			t.Error("expected metrics Run() to return in time when context is cancelled")
		}
	})
}
