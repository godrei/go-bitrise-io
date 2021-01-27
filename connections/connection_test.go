package connection

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/bitrise-io/xcode-project/pretty"
	"github.com/stretchr/testify/require"
)

func TestConnection(t *testing.T) {
	var provider AppleDeveloperConnectionProvider
	provider = NewBitriseClient(http.DefaultClient)
	conn, err := provider.GetAppleDeveloperConnection(os.Getenv("BITRISE_BUILD_URL"), os.Getenv("BITRISE_BUILD_API_TOKEN"))
	fmt.Printf("err: %s\n", err)
	fmt.Printf("conn: %s\n", pretty.Object(conn))
	require.Equal(t, 1, 2)
}
