package iamcaller

import (
	"fmt"

	"go.einride.tech/iam/iammember"
	iamv1 "go.einride.tech/iam/proto/gen/einride/iam/v1"
)

// Validate checks that resolved caller info is valid.
func Validate(caller *iamv1.Caller) error {
	if err := validateMembers(caller); err != nil {
		return fmt.Errorf("validate caller: %w", err)
	}
	if err := validateMetadata(caller); err != nil {
		return fmt.Errorf("validate caller: %w", err)
	}
	return nil
}

func validateMembers(caller *iamv1.Caller) error {
ValidateMembersLoop:
	for _, member := range caller.GetMembers() {
		// Validate the member.
		if err := iammember.Validate(member); err != nil {
			return err
		}
		// Validate that the member is present in the metadata.
		for _, metadata := range caller.GetMetadata() {
			for _, metadataMember := range metadata.GetMembers() {
				if member == metadataMember {
					continue ValidateMembersLoop
				}
			}
		}
		return fmt.Errorf("member '%s' not in metadata", member)
	}
	return nil
}

func validateMetadata(caller *iamv1.Caller) error {
	for key, metadata := range caller.GetMetadata() {
	ValidateMetadataLoop:
		for _, metadataMember := range metadata.GetMembers() {
			// Validate the metadata member.
			if err := iammember.Validate(metadataMember); err != nil {
				return err
			}
			// Validate that the metadata member is present in the top-level members.
			for _, callerMember := range caller.GetMembers() {
				if metadataMember == callerMember {
					continue ValidateMetadataLoop
				}
			}
			return fmt.Errorf("member '%s' from metadata '%s' not in caller members", metadataMember, key)
		}
	}
	return nil
}
