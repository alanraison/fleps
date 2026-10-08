// Package model contains Go types for the Google Workspace card schema
// (google.apps.card.v1), as used by Google Chat apps and Workspace add-ons.
//
// The types serialise to and from the canonical JSON representation of the
// protobuf messages. Protobuf "oneof" groups are represented as a set of
// pointer fields of which at most one may be non-nil; call Validate on a Card
// (or any individual component) to check this constraint.
//
// See https://developers.google.com/workspace/add-ons/reference/rpc/google.apps.card.v1
package model
