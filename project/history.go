package project

import (
	"github.com/monax/relic"
)

// Can be used to set the commit hash version of the binary at build time with:
// `go build -ldflags "-X github.com/hyperledger/burrow/project.commit=$(git rev-parse --short HEAD)" ./cmd/burrow`

var commit = ""

func Commit() string {
	return commit
}

func FullVersion() string {
	version := History.CurrentVersion().String()
	if commit != "" {
		return version + "+commit." + commit
	}
	return version
}

// The releases described by version string and changes, newest release first.
// The current release is taken to be the first release in the slice, and its
// version determines the single authoritative version for the next release.
//
// To cut a new release add a release to the front of this slice then run the
// release tagging script: ./scripts/tag_release.sh
var History relic.ImmutableHistory = relic.NewHistory("Hyperledger Burrow").MustDeclareReleases(
	"0.18.1",
	`Fixes for deployment`,
	"0.18.0",
	`Prerelease`,
	"0.17.1",
	`Minor tweaks to docker build file`,

	"0.17.0",
	`This is a service release with some significant ethereum/solidity compatibility improvements and new logging features. It includes:

- [Upgrade to use Tendermint v0.9.2](https://github.com/hyperledger/burrow/pull/595)
- [Implemented dynamic memory](https://github.com/hyperledger/burrow/pull/607) assumed by the EVM bytecode produce by solidity, fixing various issues.
- Logging sinks and configuration - providing a flexible mechanism for configuring log flows and outputs see [logging section in readme](https://github.com/hyperledger/burrow#logging). Various other logging enhancements.
- Fix event unsubscription
- Remove module-specific versioning
- Rename suicide to selfdestruct
- SNative tweaks

Known issues:

- SELFDESTRUCT opcode causes a panic when an account is removed. A [fix](https://github.com/hyperledger/burrow/pull/605) was produced but was [reverted](https://github.com/hyperledger/burrow/pull/636) pending investigation of a possible regression.`,

	"0.16.3",
	`This release adds an stop-gap fix to the Transact method so that it never
transfers value with the CallTx is generates.

We hard-code amount = fee so that no value is transferred
regardless of fee sent. This fixes an invalid jump destination error arising
from transferring value to non-payable functions with newer versions of solidity.
By doing this we can resolve some issues with users of the v0 RPC without making
a breaking API change.`,

	"0.16.2",
	`This release finalises our accession to the Hyperledger project and updates our root package namespace to github.com/hyperledger/burrow.

It also includes a bug fix for rpc/V0 so that BroadcastTx can accept any transaction type and various pieces of internal clean-up.`,

	"0.16.1",
	`This release was an internal rename to 'Burrow' with some minor other attendant clean up.`,

	"0.16.0",
	`This is a consolidation release that fixes various bugs and improves elements
of the architecture across the Monax Platform to support a quicker release
cadence.

#### Features and improvements (among others)
- [pull-510](https://github.com/hyperledger/burrow/pull/510) upgrade consensus engine to Tendermint v0.8.0
- [pull-507](https://github.com/hyperledger/burrow/pull/507) use sha3 for snative addresses for future-proofing
- [pull-506](https://github.com/hyperledger/burrow/pull/506) alignment and consolidation for genesis and config between tooling and chains
- [pull-504](https://github.com/hyperledger/burrow/pull/504) relicense eris-db to Apache 2.0
- [pull-500](https://github.com/hyperledger/burrow/pull/500) introduce more strongly types secure native contracts
- [pull-499](https://github.com/hyperledger/burrow/pull/499) introduce word256 and remove dependency on tendermint/go-common
- [pull-493](https://github.com/hyperledger/burrow/pull/493) re-introduce GenesisTime in GenesisDoc

- Logging system overhauled based on the central logging interface of go-kit log. Configuration lacking in this release but should be in 0.16.1. Allows powerful routing, filtering, and output options for better operations and increasing the observability of an eris blockchain. More to follow.
- Genesis making is improved and moved into eris-db.
- Config templating is moved into eris-db for better synchronisation of server config between the consumer of it (eris-db) and the producers of it (eris cli and other tools).
- Some documentation updates in code and in specs.
- [pull-462](https://github.com/hyperledger/burrow/pull/499) Makefile added to capture conventions around building and testing and replicate them across different environments such as continuous integration systems.

#### Bugfixes (among others)
- [pull-516](https://github.com/hyperledger/burrow/pull/516) Organize and add unit tests for rpc/v0
- [pull-453](https://github.com/hyperledger/burrow/pull/453) Fix deserialisation for BroadcastTx on rpc/v0
- [pull-476](https://github.com/hyperledger/burrow/pull/476) patch EXTCODESIZE for native contracts as solc ^v0.4 performs a safety check for non-zero contract code
- [pull-468](https://github.com/hyperledger/burrow/pull/468) correct specifications for params on unsubscribe on rpc/tendermint
- [pull-465](https://github.com/hyperledger/burrow/pull/465) fix divergence from JSON-RPC spec for Response object
- [pull-366](https://github.com/hyperledger/burrow/pull/366) correction to circle ci script
- [pull-379](https://github.com/hyperledger/burrow/pull/379) more descriptive error message for eris-client
`,

	"0.15.0",
	"This release was elided to synchronise release versions with tooling",

	"0.14.0",
	"This release was elided to synchronise release versions with tooling",

	"0.13.0",
	"This release was elided to synchronise release versions with tooling",

	"0.12.0",
	`This release marks the start of Eris-DB as the full permissioned blockchain node
 of the Eris platform with the Tendermint permissioned consensus engine.
 This involved significant refactoring of almost all parts of the code,
 but provides a solid foundation to build the next generation of advanced
 permissioned smart contract blockchains.

 Many changes are under the hood but here are the main externally
 visible changes:

- Features and improvements
  - Upgrade to Tendermint 0.6.0 in-process consensus
  - Support DELEGATECALL opcode in Ethereum Virtual Machine (important for solidity library calls)
  - ARM support
  - Docker image size reduced
  - Introduction of eris-client companion library for interacting with
  eris:db
  - Improved single configuration file for all components written by eris-cm
  - Allow multiple event subscriptions from same host under rpc/tendermint


- Tool changes  
  - Use glide instead of godeps for dependencies


- Testing
  - integration tests over simulated RPC calls
  - significantly improved unit tests
  - the ethereum virtual machine and the consensus engine are now top-level
  components and are exposed to continuous integration tests


- Bugfixes (incomplete list)
  - [EVM] Fix calculation of child CALL gaslimit (allowing solidity library calls to work properly)
  - [RPC/v0] Fix blocking event subscription in transactAndHold (preventing return in Javascript libraries)
  - [Blockchain] Fix getBlocks to respect block height cap.
`,
)
