package bl

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUpdateSomething(t *testing.T) {
	repository := NewMockSomethingRepository(t)
	repository.EXPECT().UpdateSomething(t.Context()).Return(nil).Once()

	svc := SomethingService{Repository: repository}

	err := svc.UpdateSomething(t.Context())
	require.NoError(t, err)
}
