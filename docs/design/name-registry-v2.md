---
title: Name Registry V2 (or: what is NameReg for?)
author: Silas Davis <silas.davis@monax.io>
category: Accounts and State
created: 2021-03-17
---

## Summary

NameReg has evolved to be used for references to addreses

## Motivation
<!--The motivation is critical It should clearly explain why the existing protocol is inadequate to address the problem that the ADR solves. ADR submissions without sufficient motivation may be rejected outright.-->

### Current implementation
The existing Name Registry (herein: namereg) is a simple key-value store that has been a long-standing part of Burrow (and ErisDB before it). It allows users to store name record with an owner and some arbitrary string data. Each entry has an expiry in block height. Before the expiry only the owner may update the entry residing a particular name, after expiry anyone may overwrite the entry. The name is globally unique and only one entry may be associated with a name at any point. The schema for an entry is:

```proto
message Entry {
    // registered name for the entry
    string Name = 1;
    // address that created the entry
    bytes Owner = 2;
    // data to store under this name
    string Data = 3;
    // block at which this entry expires
    uint64 Expires = 4;
}
```

```proto
message EntryV2 {
    // registered name for the entry
    string Namespace = 1;
    
    string Name = 1;
    // address that created the entry
    bytes Owner = 2;
    // data to store under this name
    Reference Reference = 3;
    // block at which this entry expires
    uint64 Expires = 4;
}
```

## Specification
<!--The technical specification should describe the syntax and semantics of any new feature.-->
The technical specification should describe the syntax and semantics of any new feature.

## Rationale
<!--The rationale fleshes out the specification by describing what motivated the design and why particular design decisions were made. It should describe alternate designs that were considered and related work. The rationale may also provide evidence of consensus within the community, and should discuss important objections or concerns raised during discussion.-->
The rationale fleshes out the specification by describing what motivated the design and why particular design decisions were made. It should describe alternate designs that were considered and related work. The rationale may also provide evidence of consensus within the community, and should discuss important objections or concerns raised during discussion.-->

## Backwards Compatibility
<!--All ADRs that introduce backwards incompatibilities must include a section describing these incompatibilities and their severity. The ADR must explain how the author proposes to deal with these incompatibilities. ADR submissions without a sufficient backwards compatibility treatise may be rejected outright.-->
All ADRs that introduce backwards incompatibilities must include a section describing these incompatibilities and their severity. The ADR must explain how the author proposes to deal with these incompatibilities. ADR submissions without a sufficient backwards compatibility treatise may be rejected outright.

## Test Cases
<!--Test cases for an implementation are mandatory for ADRs that are affecting consensus changes. Other ADRs can choose to include links to test cases if applicable.-->
Test cases for an implementation are mandatory for ADRs that are affecting consensus changes. Other ADRs can choose to include links to test cases if applicable.
