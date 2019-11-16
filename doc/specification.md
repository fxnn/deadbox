# Specification

This spec serves as a guideline for the current implementation.
However, it is subject to change, and some points might be of quite low priority.

## Security

* Drop secures its REST interface using TLS.
  If the Worker cannot establish a trust chain for the drops
  certificate, [Public Key Pinning](https://developer.mozilla.org/en-US/docs/Web/HTTP/Public_Key_Pinning)
  might be an option.
* The Worker provides a public key to the Drop.
  The key pair could be generated on first startup.
* When the User wants to send a request to a Worker, he identifies the
  Worker using a unique identifier.
  * The identifier must be pre-known to the user, the Drop should not
    provide a list of known Workers, which makes it harder to find
    vulnerable endpoints.
  * The identifier _equals the public key fingerprint_.
    This makes identiy spoofing harder and allows the User to trust
    that the Worker he connects to is the one he actually expects,
    without needing to trust the Drop.
* User retrieves the Worker's public key from the Drop, verifies it against
  the Worker's identifier and uses it to establish a symmetric session key
  with the Worker.
* Both requests and responses between User and Worker are now fully
  encrypted.
  Neither the Drop, nor a man in the middle attacker are able to read
  the data.

## Command Line Interface

The main program is `deadbox`, providing at least the following arguments.

* `-worker /path/to/worker.yml` makes the pocess act as an worker, connecting 
  to a drop and waiting for commands.
* `-drop /path/to/drop.yml` makes the process act as a drop, opening 
  a port providing a REST interface for both users and workers.
* `-auth /path/to/auth.yml` configures authentication for both workers and
  drops.

## `worker.yml`

The file `worker.yml` contains the configuration for an worker.
It must configure at least the following aspects.

* A textual identifier of the worker instance.
  Should be unique amongst all workers on the same drop.
* URL of one (or even multiple) drops.
* A reference to a private (and possibly public) key file, used to decrypt 
  received messages and to identify and authentificate the Worker against
  users and the drop.
  The keys will be generated automatically if not existant.
* Enabled receivers, with the appropriate configuration for each.
 * A file system receiver must be configured with the file system location, include/exclude patterns and allowed/disallowed operations.
 * An executing receiver must be configured with the command to execute when invoked.
 * Per receiver, it must be possible to enable/disable it for individual authenticated users.

## `drop.yml`

The file `drop.yml` contains the configuration for a drop.
It must configure at least the following aspects.

* Hostname and port to bind on.
* A whitelist of allowed workers, identified by their unique identifiers.

## `auth.yml`

The file `auth.yml` configures the authentication method to be used by both 
workers and drops.
It must contain at least the following information.

* Authentication method to be used, e.g. htpasswd file or OAuth 2.0 server.
* Parameters for the authentication method, e.g. location of the passwd file or URL and access data to the OAuth 2.0 server.
 * Note, that drop and worker might need different parameters.
   While the drop must be able to establish user authentication, the worker 
   must simply be able to consume a request.

## REST interface

The drop offers a REST interface to be consumed by users and workers.
It will provide at least the following endpoints.

* `POST /worker/` lets a Worker register with the Drop or update his
  information.
  This endpoint is idempotent.
* `GET /worker/{workerId}` allows users to retrieve data about a Worker
  by means of a pre-known `workerId`.
* `GET /worker/{workerId}/request` lets a Worker retrieve all pending requests
  targeted to him.
* `POST /worker/{workerId}/request` allows a User to file a new
  request targeted to a Worker.
  This endpoint is idempotent.
* `POST /worker/{workerId}/response/{requestId}` is called by a Worker to
  transport the response to a request.
  This endpoint is idempotent.
* `GET /worker/{workerId}/response/{requestId}` allows users to retrieve the
  response to a previously filed request, if already available.

As for a worker entry, at least the following information need to be contained.

* The Workers unique identification.
* A timestamp, after which the Worker may be regarded as outdated and
  removed.
  The timestamp must be in the future.
  The Drop might apply a smaller timeout than requested.
* The request types supported by the Worker, given as URN.
* The Workers public key, used to encrypt the requests, given in ASN1
  format.

As for a queue entry, which is a request, at least the following information
need to be contained.

* The requests unique identification.
* A timestamp, after which the request may be regarded as outdated and
  removed.
  The Drop might apply a smaller timeout than requested.
* The request type, given as URN.
  The request type should be contained in the list of supported request
  types as provided by the Worker.
* The request payload.
* If we want to apply routing over multiple hops, we might also need a unique
  identifier of the originating drop, and possibly an ordered list of all
  drops that transported the request.

As for a response, at least the following information need to be contained.

* The original requests unique identification.
* A timestamp, after which the response may be regarded as outdated and
  removed.
  The Drop might apply a smaller timeout than requested.
* The response payload.
* If we want to apply routing over multiple hops, we might also need a unique
  identifier of the originating drop, and possibly an ordered list of all drops
  that are still to address.

## Public Key Fingerprinting

The task of creating a public key fingerprint, and therefore a Worker
Id, is defined by the following algorithm.
The algorithm can be configured using the parameters

* `HashFunction`. The cryptographic hash function to be used, i.e.
  SHA-256.
* `Encoding`. The function for mapping the hash sum to its string
  representation, i.e. Base32.
  Base32 is chosen as default, because it performs good when it has to
  be read by human users.
* `ChallengeLevel`. At a value of `0`, it is very cheap to calculate
  the fingerprint, but it's also quite cheap to calculate a hash
  collision.
* `FingerprintLength`. When equal to the number of bytes the
  `HashFunction` produces, the fingerprint consists of the full hashsum,
  providing maximum protection against hash collision search, but also
  making the fingerprint difficult to use.

The algorithm consists of the following steps.

* Bring the public key into a binary representation,
* Add other information which must be protected against tampering:
  * the `ChallengeLevel` itself,
* Add a `challengeSolution`, which is initially `0`,
* Calculate the hashsum of these data using `HashFunction`,
* Verify that the leftmost `ChallengeLevel` bits of the hashsum are
  zero,
  * If not so, increase the `challengeSolution` by one and recalculate
    the hash,
* Encode the hashsum without the first `(ChallengeLevel+7)/8` bytes
  using the `Encoding`,
* Group the encoded hashsum in pairs of two characters,
  shorten it to the first `FingerprintLength` groups and
  separate the pairs by colons `:`.

The params `ChallengeLevel` and `FingerprintLength` should be chosen in
a way that low `FingerprintLength` is compensated with a high
`CalculationCost` and vice versa.

