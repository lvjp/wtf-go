package e2e

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/prometheus/common/expfmt"
	"github.com/prometheus/common/model"
	"github.com/stretchr/testify/require"
)

func TestMetrics(t *testing.T) {
	var endpoint string
	if endpoint = os.Getenv("TEST_WTF_GO_BACKEND_ENDPOINT"); endpoint == "" {
		t.Skip("TEST_WTF_GO_BACKEND_ENDPOINT environment variable is not set, skipping test")
	}

	req, err := http.NewRequestWithContext(
		t.Context(),
		http.MethodGet,
		endpoint+"/metrics",
		nil,
	)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	parser := expfmt.NewTextParser(model.UTF8Validation)
	metrics, err := parser.TextToMetricFamilies(resp.Body)
	require.NoError(t, err)

	require.Contains(t, metrics, "go_info")
	require.Contains(t, metrics, "wtf_go_info")
	require.Contains(t, metrics, "wtf_go_start_date_timestamp")

	// Check that wtf_go_start_date_timestamp value is in the past
	startDateMetric := metrics["wtf_go_start_date_timestamp"]
	require.NotNil(t, startDateMetric)
	require.Len(t, startDateMetric.Metric, 1) // Should have exactly one metric

	startTimestamp := float64(startDateMetric.Metric[0].GetGauge().GetValue())
	currentTimestamp := float64(time.Now().Unix())

	require.Less(t, startTimestamp, currentTimestamp, "wtf_go_start_date_timestamp should be in the past")
}
