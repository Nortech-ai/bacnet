package bacnet

import (
	"testing"

	"github.com/Nortech-ai/bacnet/test_utils"
)

func TestParseWhois(t *testing.T) {
	test_utils.TestParseWhois(t, Parse)
}
func TestParseReadProperty(t *testing.T) {
	test_utils.TestParseReadProperty(t, Parse)
}
func TestParseUnicastIam(t *testing.T) {
	test_utils.TestParseUnicastIam(t, Parse)
}
func TestParseReadPropertyMultiple(t *testing.T) {
	test_utils.TestParseReadPropertyMultiple(t, Parse)
}
