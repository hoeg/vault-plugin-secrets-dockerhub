package dockerhub

import (
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
)

func getStringFrom(data *framework.FieldData, key string) string {
	v := data.Get(key)
	if s, ok := v.(string); ok && s != "" {
		return s
	}
	panic(fmt.Sprintf("no string from %s", key))
}
