package cpu

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCPU(t *testing.T) {

	_, err := Get()
	require.NoError(t, err)

	time.Sleep(time.Second)

	data, err := Get()
	require.NoError(t, err)
	require.NotNil(t, data)
	require.GreaterOrEqual(t, data.User, 0.0)
	require.LessOrEqual(t, data.User, 100.0)
	require.GreaterOrEqual(t, data.System, 0.0)
	require.LessOrEqual(t, data.System, 100.0)
	require.GreaterOrEqual(t, data.Idle, 0.0)
	require.LessOrEqual(t, data.Idle, 100.0)
}
