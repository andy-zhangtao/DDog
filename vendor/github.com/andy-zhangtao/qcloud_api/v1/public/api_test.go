package public

import "testing"

func TestPublicParam(t *testing.T) {
	param := PublicParam("DescribeInstances","sh","123456")
	t.Log(param)
}
