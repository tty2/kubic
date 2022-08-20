package shared

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_TabItem_String(t *testing.T) {
	t.Parallel()
	rq := require.New(t)

	t.Run("namespace", func(t *testing.T) {
		t.Parallel()

		rq.Equal(namespacesTabTitle, NamespacesTab.String())
	})

	t.Run("deployments", func(t *testing.T) {
		t.Parallel()

		rq.Equal(deploymentsTabTitle, DeploymentsTab.String())
	})

	t.Run("pods", func(t *testing.T) {
		t.Parallel()

		rq.Equal(podsTabTitle, PodsTab.String())
	})

	t.Run("any", func(t *testing.T) {
		t.Parallel()

		rq.Equal("", AnyTab.String())
	})
}

func Test_GetTabItems(t *testing.T) {
	t.Parallel()
	rq := require.New(t)

	t.Run("ok", func(t *testing.T) {
		t.Parallel()

		tt := GetTabItems()

		rq.Len(tt, 3)
		rq.Equal(NamespacesTab, tt[0])
		rq.Equal(DeploymentsTab, tt[1])
		rq.Equal(PodsTab, tt[2])
	})
}
