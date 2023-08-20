# Thema

Thema is a system for writing schemas. Much like JSON Schema or OpenAPI, it is general-purpose and its most obvious application is as an [IDL](https://en.wikipedia.org/wiki/Interface_description_language). However, those systems treat _changing_ schemas as out of scope: a single version of a schema for some object is the atomic unit, and versioning is left to opaque strings in external systems like git or HTTP. Thema, by contrast, makes schema change a first-class system property: the atomic unit is the _set_ of schema for some object, iteratively appended to over time as requirements evolve.

Thema's approach is novel, so an analogy to the familiar may help. ["Branching by abstraction"](https://martinfowler.com/bliki/BranchByAbstraction.html) suggests that you refactor large applications not with long-running VCS branches and big-bang merges, but by letting old and new code live side-by-side on `main`, and choosing between them with logical gates, like [feature flags](https://featureflags.io/feature-flags/). Thema is "schema versioning by abstraction": all versions of a schema live side-by-side on `main`, within logical structures Thema defines.

This holistic view allows Thema to act like a typechecker, but for change-safety _between_ schema versions: either schema versions must be backwards compatible, or there must exist logic to translate a valid instance of schema from one schema version to the next. [CUE](https://cuelang.org), the language in which Thema schemas are written, allows Thema to [mechanically verify these properties](#Maturity).

These capabilities make Thema a general framework for decoupling the evolution of communicating systems. This can be outward-facing: Thema's guardrails allow anyone to create APIs with Stripe's renowned [backwards compatibility](https://stripe.com/docs/upgrades) guarantees. Or it can be inward-facing: or to change the messages passed in a mesh of microservices without intricately orchestrating deployment.

Learn more in our [docs](https://github.com/grafana/thema/tree/main/docs), or in this [overview video](https://www.youtube.com/watch?v=PpoS_ThntEM)! (Some things have been renamed since that video, but the logic is unchanged.)

## Usage

Thema defines the way schemas are written, organizing each object's history into a "lineage." Once authored, Thema also provides tools for working with lineages via a few [basic operations](https://github.com/grafana/thema/blob/main/docs/overview.md#about-thema-operations). There are a few different usage patterns, all largely equivalent in capability:

* **CLI:** a CLI command that provides access to Thema's basic operations, one lineage per invocation. Use it for fast exploration and testing of schemas, or as a tool in CI.
* **Server:** An HTTP server that provides access to Thema's basic operations for a configurable set of lineages. Run it as a stateless sidecar in your infrastructure or microservice mesh.
* **Library:** a library, importable in your application code, that provides a convenient interface to Thema's basic operations, as well as helpers for common usage patterns. Naturally the most flexible, and the recommended approach for creating new helpers, such as code generators, API generators, or a whole Kubernetes operator framework. (Currently only for Go[^evaluator])

The CLI and server modes are bundled together in the `thema` command. To install:

```bash
go install github.com/grafana/thema/cmd/thema@latest
```

## Maturity

Thema is a young project. The goals are large, but bounded: we will know when the core system is complete. And it mostly is, now - though some breaking changes to how schemas are written are planned before reaching stability.

It is not yet recommended to replace established, stable systems with Thema, but experimenting with doing so is reasonable (and appreciated!). For newer projects, Thema may be a good choice today; the decision is likely to come down to whether the long-term benefit of a simpler architecture for authoring, composing and evolving schema will offset the short-term cost of some incomplete functionality and breaking changes.

## Prior/Related Art

A number of systems partially overlap with Thema - for some data, rolling together a set of schema with the relations between those schema.

* [Project Cambria](https://www.inkandswitch.com/cambria/) - Thema's closest analogue. Limited in verifiability by (intentionally) being without a notion of linear schema ordering and versioning, and because schema and translations are written in a Turing complete language (Typescript).
* [Kubernetes resources and webhook conversions](https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definition-versioning/#specify-multiple-versions) - Similar goals: multiple versions of resources (schema) and convertibility between them. Limited in verifiability by relying on convention for grouping schemas, and by expressing translation in a Turing complete language (Go).
* [Stripe's HTTP API](https://stripe.com/docs/upgrades) - exhibits the backwards compatibility properties an API can have that arise from a schema system with translatability.


[^evaluator]:
    Using Thema as a library in a language depends on a CUE evaluator for that language. Currently, the only CUE evaluator is written in Go.
