package iamcaller

import (
	"sort"

	iamv1 "go.einride.tech/iam/proto/gen/einride/iam/v1"
)

// Add adds the metadata resolved from the provided key to the provider caller info.
func Add(caller *iamv1.Caller, key string, metadata *iamv1.Caller_Metadata) {
MemberLoop:
	for _, member := range metadata.Members {
		for _, existingMember := range caller.Members {
			if member == existingMember {
				continue MemberLoop
			}
		}
		caller.Members = append(caller.Members, member)
	}
	sort.Slice(caller.Members, func(i, j int) bool {
		return caller.Members[i] < caller.Members[j]
	})
	if caller.Metadata == nil {
		caller.Metadata = map[string]*iamv1.Caller_Metadata{}
	}
	caller.Metadata[key] = metadata
}
