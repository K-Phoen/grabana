// Package vmux provides utilities that make it easy to implement "version
// multiplexing" with your Thema lineages.
//
// Version multiplexing is taking an input of []byte expected to represent an
// object schematized by some Thema lineage, and automatically decoding,
// validating, and translating it to a single schema version in that lineage.
// The effect of this is a Go program that "accepts all, sees one" - all
// historical schema versions are accepted, but the program can be written as
// though only a single version exists.
//
// The generic utilities in this package reduce version muxing to a single
// function call. Still, they are pure convenience: this package relies
// solely on the public interface of [github.com/grafana/thema], and can
// be reimplemented with different tradeoffs elsewhere if needed.
package vmux
