package background

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/hexops/autogold"
	"github.com/stretchr/testify/require"

	"github.com/sourcegraph/sourcegraph/internal/search/result"
)

func TestWebhook(t *testing.T) {
	eu, err := url.Parse("https://sourcegraph.com")
	require.NoError(t, err)

	action := actionArgs{
		MonitorDescription: "My test monitor",
		ExternalURL:        eu,
		MonitorID:          42,
		Query:              "repo:camdentest -file:id_rsa.pub BEGIN",
		Results:            []*result.CommitMatch{&diffResultMock, &commitResultMock},
		IncludeResults:     false,
		HostList:           "loopback",
	}

	t.Run("no error", func(t *testing.T) {
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			b, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			autogold.Equal(t, autogold.Raw(b))
			w.WriteHeader(200)
		}))
		defer s.Close()

		err := postWebhook(context.Background(), s.URL, action)
		require.NoError(t, err)
	})

	// If these tests fail, be sure to check that the changes are correct here:
	// https://app.slack.com/block-kit-builder/T02FSM7DL#%7B%22blocks%22:%5B%5D%7D
	t.Run("golden with results", func(t *testing.T) {
		actionCopy := action
		actionCopy.IncludeResults = true

		j, err := json.Marshal(generateWebhookPayload(actionCopy))
		require.NoError(t, err)

		autogold.Equal(t, autogold.Raw(j))
	})

	t.Run("error is returned", func(t *testing.T) {
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			b, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			autogold.Equal(t, autogold.Raw(b))
			w.WriteHeader(500)
		}))
		defer s.Close()

		err := postWebhook(context.Background(), s.URL, action)
		require.Error(t, err)
	})

	t.Run("loopback requests are blocked", func(t *testing.T) {
		// only allow external requests, excluding loopback
		action.HostList = "external"
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			b, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			autogold.Equal(t, autogold.Raw(b))
			w.WriteHeader(500)
		}))
		defer s.Close()

		err := postWebhook(context.Background(), s.URL, action)
		require.Error(t, err)
	})
}

func TestTriggerTestWebhookAction(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		autogold.Equal(t, autogold.Raw(b))
		w.WriteHeader(200)
	}))
	defer s.Close()

	err := SendTestWebhook(context.Background(), "My test monitor", s.URL, "loopback")
	require.NoError(t, err)
}
