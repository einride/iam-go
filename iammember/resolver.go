package iammember

import (
	"context"
	"hash/crc64"
)

// crcTable used for calculating checksums.
var crcTable = crc64.MakeTable(crc64.ISO)

// Resolver resolves the IAM member identifiers for a caller context.
type Resolver interface {
	ResolveIAMMembers(context.Context) (ResolveResult, error)
}

// Metadata contains IAM members partitioned by which gRPC metadata key they were resolved from.
type Metadata map[string][]string

// Add a member to the specified metadata key.
func (m Metadata) Add(key, member string) {
	for _, existingMember := range m[key] {
		if existingMember == member {
			return
		}
	}
	m[key] = append(m[key], member)
}

// AddAll adds all members from another metadata instance.
func (m Metadata) AddAll(other Metadata) {
	for key, members := range other {
		for _, member := range members {
			m.Add(key, member)
		}
	}
}

// ResolveResult is the output from a Resolver.
type ResolveResult struct {
	// Checksum of the metadata that the members were resolved from.
	Checksum uint64
	// Members are the resolved IAM members.
	Members []string
	// Metadata are the resolved IAM members partitioned by which metadata key they were resolved from.
	Metadata Metadata
}

// AddChecksum adds the provided metadata key and value to the checksum.
func (r *ResolveResult) AddChecksum(key, value string) {
	r.Checksum = crc64.Update(r.Checksum, crcTable, []byte(key))
	r.Checksum = crc64.Update(r.Checksum, crcTable, []byte(value))
}

// Add a member resolved from the provided metadata key and value.
func (r *ResolveResult) Add(key string, member string) {
	var hasMember bool
	for _, existingMember := range r.Members {
		if member == existingMember {
			hasMember = true
			break
		}
	}
	if !hasMember {
		r.Members = append(r.Members, member)
	}
	if r.Metadata == nil {
		r.Metadata = make(Metadata)
	}
	r.Metadata.Add(key, member)
}

// AddAll adds all the resolved members from another ResolveResult.
func (r *ResolveResult) AddAll(other ResolveResult) {
	// Update checksum. Since we don't have the original metadata values, we simply add the other checksum.
	r.Checksum = crc64.Update(r.Checksum, crcTable, []byte{byte(other.Checksum)})
	// Add ordered members first to maintain order.
	for _, member := range other.Members {
		var hasMember bool
		for _, existingMember := range r.Members {
			if member == existingMember {
				hasMember = true
				break
			}
		}
		if !hasMember {
			r.Members = append(r.Members, member)
		}
	}
	if r.Metadata == nil {
		r.Metadata = make(Metadata)
	}
	r.Metadata.AddAll(other.Metadata)
}
