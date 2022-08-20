package k8s

import (
	"testing"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
)

func Test_ageToString(t *testing.T) {
	t.Parallel()

	t.Run("seconds", func(t *testing.T) {
		t.Parallel()

		res := ageToString(30)

		rq := require.New(t)
		rq.Equal("30s", res)

		res = ageToString(59)

		rq.Equal("59s", res)
	})
	t.Run("minutes", func(t *testing.T) {
		t.Parallel()

		res := ageToString(60)

		rq := require.New(t)
		rq.Equal("1m", res)

		res = ageToString(3599)

		rq.Equal("59m", res)
	})
	t.Run("hours", func(t *testing.T) {
		t.Parallel()

		res := ageToString(3600)

		rq := require.New(t)
		rq.Equal("1h", res)

		res = ageToString(86399)

		rq.Equal("23h", res)
	})
	t.Run("days", func(t *testing.T) {
		t.Parallel()

		res := ageToString(86400)

		rq := require.New(t)
		rq.Equal("1d", res)

		res = ageToString(2678399)

		rq.Equal("30d", res)
	})
	t.Run("months", func(t *testing.T) {
		t.Parallel()

		res := ageToString(2678400)

		rq := require.New(t)
		rq.Equal("1M", res)

		res = ageToString(31535999)

		rq.Equal("11M", res)
	})
	t.Run("months", func(t *testing.T) {
		t.Parallel()

		res := ageToString(31536000)

		rq := require.New(t)
		rq.Equal("1Y", res)

		res = ageToString(315360000)

		rq.Equal("10Y", res)
	})
}

func Test_getReadyOfListCont(t *testing.T) {
	t.Parallel()
	rq := require.New(t)

	t.Run("1/1", func(t *testing.T) {
		t.Parallel()

		ss := []corev1.ContainerStatus{
			{
				Ready: true,
			},
		}

		rq.Equal("1/1", getReadyOfListCont(ss))
	})
	t.Run("3/4", func(t *testing.T) {
		t.Parallel()

		ss := []corev1.ContainerStatus{
			{
				Ready: true,
			},
			{
				Ready: true,
			},
			{
				Ready: false,
			},
			{
				Ready: true,
			},
		}

		rq.Equal("3/4", getReadyOfListCont(ss))
	})
	t.Run("0/3", func(t *testing.T) {
		t.Parallel()

		ss := []corev1.ContainerStatus{
			{
				Ready: false,
			},
			{
				Ready: false,
			},
			{
				Ready: false,
			},
		}

		rq.Equal("0/3", getReadyOfListCont(ss))
	})
}

func Test_getRestartsCount(t *testing.T) {
	t.Parallel()
	rq := require.New(t)

	t.Run("0", func(t *testing.T) {
		t.Parallel()

		ss := []corev1.ContainerStatus{
			{
				Ready: false,
			},
			{
				Ready: false,
			},
			{
				Ready: false,
			},
		}
		rq.Equal(0, getRestartsCount(ss))
	})

	t.Run("2", func(t *testing.T) {
		t.Parallel()

		ss := []corev1.ContainerStatus{
			{
				Ready: true,
			},
			{
				Ready: false,
			},
			{
				Ready: true,
			},
		}
		rq.Equal(0, getRestartsCount(ss))
	})

	t.Run("4", func(t *testing.T) {
		t.Parallel()

		ss := []corev1.ContainerStatus{
			{
				Ready: true,
			},
			{
				Ready: true,
			},
			{
				Ready: true,
			},
			{
				Ready: true,
			},
		}
		rq.Equal(0, getRestartsCount(ss))
	})
}
