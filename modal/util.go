package modal

import (
	"strings"
)

//拼接k
func JointStr(s1 ...string) string {
	var build strings.Builder
	for k, v := range s1 {
		build.WriteString(v)
		if k != (len(s1) - 1) {
			build.WriteString(":")
		}
	}
	return build.String()

}
