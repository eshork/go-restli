package suite

import (
	"testing"

	conflictresolution "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/conflictResolution"
	"github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite"
	"github.com/stretchr/testify/require"

	. "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite/typerefs/collectionTyperef"
)

func (s *TestServer) CollectionTyperefGet(t *testing.T, c Client) {
	id := testsuite.Url("http://rest.li")
	message, err := c.Get(id)
	require.NoError(t, err)
	require.Equal(t, &conflictresolution.Message{Message: "test message"}, message)
}
