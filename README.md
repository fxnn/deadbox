# Deadbox

A dead letter box for your data.

Access data and devices in your private network without Dynamic DNS, port opening etc.
A program in your private network connects to an internet-accessible server, ready to receive encrypted and authorized requests.

[![Build Status](https://travis-ci.org/fxnn/deadbox.svg?branch=master)](https://travis-ci.org/fxnn/deadbox)
[![codecov](https://codecov.io/gh/fxnn/deadbox/branch/master/graph/badge.svg)](https://codecov.io/gh/fxnn/deadbox)
[![0pdd.com](http://www.0pdd.com/svg?name=fxnn/deadbox)](http://www.0pdd.com/p?name=fxnn/deadbox)
[![GoDoc](https://godoc.org/github.com/fxnn/deadbox?status.svg)](https://godoc.org/github.com/fxnn/deadbox)

*I'm very happy to hear your feedback and ideas. Simply file an issue!*

## Idea

### Problem

When a user wants to access his private data from the internet, he currently has two possibilities:
* Either he uses a public cloud, like Google Drive, Microsoft OneDrive etc.
  This forces the user to _always_ store his data on foreign servers.
  He is **not able to control** who may access these data, as no popular (free of charge) cloud provider has support for encryption.
* Or he runs a private cloud (ownCloud, Nextcloud) on NAS devices.
  This forces the user to open a port to his private network behind a public DNS name, inducing a severe security risk.
  Security holes in the router, operating system or application software allow an attacker to enter the users private network and data.
  The **risk of such a security hole is high**, as Dynamic DNS and private clouds rely on a broad interface and many different technologies.

Having a way of remotely accessing private data,
without neither storing them on foreign servers nor opening up one's private network using large interfaces,
would significantly improve security of all parts of the private network.

### Approach
_Deadbox_ is an application that combines the concepts of peer to peer and 
message bus to establish communication between public and private networks.
It consists of two parts: _Workers_ and _Drops_.
* The Worker runs on devices in your private network and connects to the
  Drop,
  which is located on a machine somewhere publicly available.
* The User accesses the Drop through a **web-based UI**,
  sending usecase based requests like "search for a file by name" or
  "upload a file".
  The Drop stores the requests in a queue.
  The Worker requests all items in his queue from the Drop,
  processes them and sends the answers back to the Drop.
  From there, they can be requested by the User's UI.
* All requests and respones are **end-to-end encrypted**.
  The Drop cannot read them.
  Though the Drop stores the Worker's certified public keys, 
  the User doesn't need to trust the Drop, as he needs to know the public
  key's fingerprint (it's the Worker's connection id).
  Thus, a secure connection setup is possible.
* The User's UI runs in the web browser and is therefore easily available.
  While the browser needs to be trusted, the server delivering the webapp
  remains untrusted, as we allow for the webapp's easy verification. 

### Advantages
The User's experience is similar to that of a typical cloud application.
He uploads and downloads files or accesses other functionality through a comfortable web-based UI.

Still, the system provides more security, due to the following reasons.
* The private network provides no interface vulnerable to probing techniques like port scanning.
* The private network exposes a very small interface, protected through well-known authentication techniques.
  It is additionally protected through the Drop, who won't process unauthenticated requests.
* The publicly available Drop does not store unencrypted private data at all.
  The data stored by the Drop are deleted after request completion or within a configurable timeout.
  Therefore, attacking the Drop's data storage will not allow to access private information.
* Communication with the private network may only happen via the Drop,
  which therefore acts as a shield against DoS or bruteforce attacks.
  The availability of the Drop might be in risk, but the private network is secured.
* The webapp can be trusted, as it's carefully engineered and easily verified.

By combining Worker and Drop in one single binary and adding some basic routing,
a peer to peer network emerges, allowing to combine Workers from different 
private networks using tree-like or even rather random structures.
One _deadbox_ instance could serve as worker and drop at the same time, loading
requests from foreign Drops, providing them to other Workers.

This offers a high flexibility in retrieving data and sending commands from/to
private devices.

### Disadvantages
Compared to a centralized cloud-based storage solution,
the _deadbox_ suffers under limited availability and increased configuration 
effort.
* The service is only available when the devices in the private network are on
  and connected.
  However, this could be circumvented by letting the Drop not only store the
  encrypted requests, but also the encrypted responses for some limited time.
  Of course, this is not feasible in all use cases and, in any case, leads
  to possibly high latency.
* The user needs to install, configure and maintain instances on at least two
  different devices, from which one must be in the public internet.
  This could be circumvented by providing a central cloud service, maybe on a
  subscription base.

### Challenges
* The JavaScript webapp needs to be carefully engineered, so that no malware
  can be injected into the code.
  Also, the means of verification must be very simple and effective.
* While the Worker needs to open the connection to the Drop,
  the Drop must be able to push requests to the worker, so that the user won't 
  experience a high latency.
  Similarly, the Worker must be able to push his responses back to the Drop and
  thus to the user.
* When thinking about authentication, it is clear that the Worker should not 
  need to trust the Drop.
  Otherwise, in terms of security, the drop would be a single point of failure.
  Instead, it would be better if both could rely on a common third party for 
  establishing authentication, as it is possible with e.g. OpenID.

### Use cases and examples
The _deadbox_ could be used in following scenarios.
* A user wants to access his home NAS from the internet, while retaining a high level of security.
* A user wants to orchestrate smart home devices.
* A company wants to access and control multiple devices in different networks from a single interface. 

Possible use cases are as follows.
* The user finds a file on one of his devices by the files name, location or contents.
* The user downloads a file from one of his devices.
* The user downloads a folder from one of his devices.
* The user uploads a file to one of his devices.
* The user establishes a regular synchronization of folders between a portable device and a device inside his private network.
* The user turns on a device using wake on lan.
* The user runs preconfigured scripts and programs on one of his devices.
* The user controls a smart home device, like setting the room temperature.
* The user retrieves data from a smart home device, like a webcam or a thermostat.

### Related Work
The following projects are similar to _deadbox_.
* [Syncthing](https://syncthing.net) is an open source, multi platform synchronization solution.
  If you want to synchronize the same set of directories between the same set of devices for some time, use Syncthing! 
  However, if you'd rather have some publicly reachable endpoint for your private devices and use it to perform actions (of any kind) on your device, use _deadbox_.
* [ngrok](https://ngrok.com/) is a proprietary tunnel to your device.
  If you want to access any kind of webapp or TCP service (like Nextcloud), running on your device, itself running behind a firewall, from the public internet, use ngrok.
  However, if you'd rather have a fully open source solution, without the need to run multiple server components, or if you want to support non-interactive scenarios, where you want to enqueue actions (like sending files) while the target device is offline, use _deadbox_.

## Specification

This spec serves as a guideline for the current implementation.
However, it is subject to change, and some points might be of quite low priority.

### Security

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

### Command Line Interface

The main program is `deadbox`, providing at least the following arguments.

* `-worker /path/to/worker.yml` makes the pocess act as an worker, connecting 
  to a drop and waiting for commands.
* `-drop /path/to/drop.yml` makes the process act as a drop, opening 
  a port providing a REST interface for both users and workers.
* `-auth /path/to/auth.yml` configures authentication for both workers and
  drops.

### `worker.yml`

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

### `drop.yml`

The file `drop.yml` contains the configuration for a drop.
It must configure at least the following aspects.

* Hostname and port to bind on.
* A whitelist of allowed workers, identified by their unique identifiers.

### `auth.yml`

The file `auth.yml` configures the authentication method to be used by both 
workers and drops.
It must contain at least the following information.

* Authentication method to be used, e.g. htpasswd file or OAuth 2.0 server.
* Parameters for the authentication method, e.g. location of the passwd file or URL and access data to the OAuth 2.0 server.
 * Note, that drop and worker might need different parameters.
   While the drop must be able to establish user authentication, the worker 
   must simply be able to consume a request.

### REST interface

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

### Public Key Fingerprinting

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

