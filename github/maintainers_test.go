package github

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseMaintainer(t *testing.T) {
	require := require.New(t)

	input := "My Name <me@example.com> (@my_user)"
	expected := &User{Name: "My Name", Email: "me@example.com", Handle: "my_user"}
	actual, err := ParseMaintainer(input)
	require.NoError(err)
	require.Equal(expected, actual)

	input = "My Name <me@example.com> (my_user)"
	expected = &User{Name: "My Name", Email: "me@example.com", Handle: "my_user"}
	actual, err = ParseMaintainer(input)
	require.NoError(err)
	require.Equal(expected, actual)

	input = "* My Name <me@example.com> (@my_user)"
	expected = &User{Name: "My Name", Email: "me@example.com", Handle: "my_user"}
	actual, err = ParseMaintainer(input)
	require.NoError(err)
	require.Equal(expected, actual)
}
