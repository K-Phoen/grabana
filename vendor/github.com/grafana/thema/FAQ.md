# FAQ

## Why, why, _whyyy_ does the universe need yet another schema system?

We asked the question, "what if a schema system managed backwards compatibility by construction, rather than post hoc, best effort, via separate tooling?" Thema is the answer we came up with - the minimal set of logical constructs for addressing the problem:

1. Schemas themselves
2. Fully general backwards compatibility checks
3. Logic for translating instances across schema versions

No existing system was reasonably capable of encapsulating all of these into a unified, portable, verifiable structure.

## Where does the name "Thema" come from?

Thema is a [portmanteau](https://en.wikipedia.org/wiki/Portmanteau) of "[Ship of Theseus](https://en.wikipedia.org/wiki/Ship_of_Theseus) with "Schema".

The Ship of Theseus is a thought experiment that asks about the identity of objects as their constituent parts are replaced over time. Thema is a system that allows us to specify a "thing" in terms of its constituent parts (as all schema do), as well as the evolution of those parts incrementally over time, while still retaining its identity as that "thing".

## You can't fool me. Breaking changes are breaking - how can they possibly be made non-breaking?

That's true. A breaking change to a contract like a schema is still breaking.

What Thema does is change the nature of the contract between communicating systems. Instead of agreeing on a single schema, systems agree on the whole Thema lineage as the contract, with all the guarantees about translation between schema versions that that entails.

For example, say system `A` accepts messages which is comprised of a single field named `foo`, which has value of type `int64`. System `B` accepts the contract, and starts sending messages to `A` according to this schema. What Thema changes is not the schema in use by either `B` or `A` at any one time, but the agreement `A` and `B` implicitly make when `B` decides to start communicating with `A`:

* **Traditional Schema:** `A` promises that messages with field `foo` containing an `int64` value will be valid in perpetuity.
* **Thema:** `A` promises that messages with a field `foo` containing an `int64` will either be valid itself, or will be translatable into a valid message, in perpetuity.

Thema shifts the contract up a level of abstraction - from rigid adherence to the contents of an individual schema, to the meta-property of relations between the schemas within a lineage.

## Is Thema as expressive as other schema systems?

Thema is just a thin layer of naming patterns and logical constraints atop of CUE itself, which makes this largely a question about CUE's expressiveness.

For the most part, yes, CUE is comparably expressive to other common schema systems, like JSON Schema and OpenAPI. There are some areas where CUE is less expressive, and some where it's more. (TODO - links to more relevant info)

## What definition of "backwards compatibility" does Thema use in its checks?

[CUE's definition of subsumption](https://cuelang.org/docs/concepts/logic): does `A` subsume `B`? If so, then `A` is backwards compatible with `B`.

This definition is precise. Persnickety, even. But a design premise of Thema is that, because Thema makes breaking changes a manageable problem, precision is preferable to permissiveness: having to write a lens for eyeroll-inducing reasons is better than a pseudoguarantee people can't really depend on.

## Aren't breaking changes evil? Isn't Thema encouraging bad behavior?

If you are committed to believing this, we cannot offer definitive contradictory proof.

Our foundational belief is that, while breaking changes can cause considerable pain, that pain has not been, and is unlikely to ever be, sufficient basis for system authors to stop making breaking changes.

Given this premise, the best course of action is to create patterns that allow breaking changes made by schema authors to be effectively managed by schema consumers. Thema is the simplest such pattern we can imagine: it turns "breaking" changes from hard, brittle failures into softer questions of risk management.

## Why did Thema make up a new version numbering system instead of just using [semver](https://semver.org)?

Thema versions, unlike most version numbering systems, are not an arbitrary declaration by the schema author. Rather, version numbers are derived from the position of the schema within the lineage's list of sequences. Sequence position, in turn, is governed by Thema's checked invariants on backwards compatibility and lens existence.

The association of version numbers with verifiable properties grants Thema's versions transparency and precision - within the bounds of what is expressible as schema. Systems like Semantic Versioning are, by contrast, opaque; it wouldn't be unreasonable to highlight this difference by calling Thema's version numbering system "Syntactic Versioning".

## How do I express prerelease-type concepts: "alpha", "beta", etc.?

You don't. 

Semantic versioning [explicitly](https://semver.org/#spec-item-9) grants prereleases an exception to its compatibility semantics. This makes each [contiguous series of] prerelease a subsequence where "anything goes."

Thema takes the stance that it is preferable to _never_ suspend version number-implied guarantees, and instead lean hard into the system of lenses, translations, and lacunas. In other words, it's fine to experiment and make breaking changes within your Thema, so long as you write lenses and lacunas that can lead your users' objects to solid ground.

Support for indicating a maturity level on individual schema may be added in the future. But it would have no bearing on core Thema guarantees. Instead, maturity would be an opaque string, used purely for signalling between humans: "we're really not sure about this yet; future lenses for translating from this schema may be sloppy!"