package grafana

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClient_CreateDashboards(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload := &CreateDashboardReq{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "failed to unmarshal req", http.StatusBadRequest)
			return
		}
		b, err := ioutil.ReadFile("testdata/good_resp.json")
		if err != nil {
			http.Error(w, "failed to respond", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(b)
	}))
	t.Parallel()

	t.Run("good 1 dashboard", func(t *testing.T) {
		cl := &Client{cl: srv.Client(), token: "1234", url: srv.URL}

		b, err := ioutil.ReadFile("testdata/good_dashboard.json")
		require.NoError(t, err)

		dashboards := map[string][]byte{
			"good_dashboard.json": b,
		}
		require.NoError(t, cl.CreateDashboards(dashboards))
	})
}
