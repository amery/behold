Behold is a generic querying and indexing layer for Go
======================================================

Behold is based on the *amazing* work by Tim Shannon on
[BoltHold](https://github.com/timshannon/bolthold) and
[BadgerHold](https://github.com/timshannon/badgerhold), and
adapts to _any database_ by using interfaces.

The goal is to create a simple higher level interface that
simplifies dealing with Go Types and finding data.

Queries
-------
Queries are chain-able construct that filter out any data that
doesn't match the criteria.

See also
--------
* [BoltHold](https://github.com/timshannon/bolthold)
* [BadgerHold](https://github.com/timshannon/badgerhold)
