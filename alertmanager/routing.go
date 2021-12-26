package alertmanager

import (
	"github.com/K-Phoen/sdk"
)

type RoutingPolicyOption func(policy *RoutingPolicy)

type RoutingPolicy struct {
	builder *sdk.NotificationRoutingPolicy
}

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

func TagEq(tag string, value string) RoutingPolicyOption {
	return func(policy *RoutingPolicy) {
		policy.builder.ObjectMatchers = append(policy.builder.ObjectMatchers, sdk.AlertObjectMatcher{
			tag, "=", value,
		})
	}
}

func TagNeq(tag string, value string) RoutingPolicyOption {
	return func(policy *RoutingPolicy) {
		policy.builder.ObjectMatchers = append(policy.builder.ObjectMatchers, sdk.AlertObjectMatcher{
			tag, "!=", value,
		})
	}
}

func TagMatches(tag string, regex string) RoutingPolicyOption {
	return func(policy *RoutingPolicy) {
		policy.builder.ObjectMatchers = append(policy.builder.ObjectMatchers, sdk.AlertObjectMatcher{
			tag, "=~", regex,
		})
	}
}

func TagNotMatches(tag string, regex string) RoutingPolicyOption {
	return func(policy *RoutingPolicy) {
		policy.builder.ObjectMatchers = append(policy.builder.ObjectMatchers, sdk.AlertObjectMatcher{
			tag, "!~", regex,
		})
	}
}
