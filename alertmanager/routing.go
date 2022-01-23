package alertmanager

import (
	"github.com/K-Phoen/sdk"
)

// RoutingPolicyOption represents an option that can be used to configure a
// routing policy.
type RoutingPolicyOption func(policy *RoutingPolicy)

// RoutingPolicy represents a routing policy.
type RoutingPolicy struct {
	builder *sdk.NotificationRoutingPolicy
}

// Policy defines a routing policy that applies to the given contact point.
// All the options given on this policy will be combined using a logical "AND".
func Policy(contactPoint string, opts ...RoutingPolicyOption) RoutingPolicy {
	policy := &RoutingPolicy{
		builder: &sdk.NotificationRoutingPolicy{
			Receiver:       contactPoint,
			ObjectMatchers: nil,
		},
	}

	for _, opt := range opts {
		opt(policy)
	}

	return *policy
}

// TagEq defines an equality ("=") constraint between the given tag and value.
func TagEq(tag string, value string) RoutingPolicyOption {
	return func(policy *RoutingPolicy) {
		policy.builder.ObjectMatchers = append(policy.builder.ObjectMatchers, sdk.AlertObjectMatcher{
			tag, "=", value,
		})
	}
}

// TagNeq defines a non-equality ("!=") constraint between the given tag and value.
func TagNeq(tag string, value string) RoutingPolicyOption {
	return func(policy *RoutingPolicy) {
		policy.builder.ObjectMatchers = append(policy.builder.ObjectMatchers, sdk.AlertObjectMatcher{
			tag, "!=", value,
		})
	}
}

// TagMatches defines a similarity ("=~") constraint between the given tag and regex.
func TagMatches(tag string, regex string) RoutingPolicyOption {
	return func(policy *RoutingPolicy) {
		policy.builder.ObjectMatchers = append(policy.builder.ObjectMatchers, sdk.AlertObjectMatcher{
			tag, "=~", regex,
		})
	}
}

// TagNotMatches defines a non-similarity ("!~") constraint between the given tag and regex.
func TagNotMatches(tag string, regex string) RoutingPolicyOption {
	return func(policy *RoutingPolicy) {
		policy.builder.ObjectMatchers = append(policy.builder.ObjectMatchers, sdk.AlertObjectMatcher{
			tag, "!~", regex,
		})
	}
}
