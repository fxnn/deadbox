# Deadbox
[![0pdd.com](http://www.0pdd.com/svg?name=fxnn/deadbox)](http://www.0pdd.com/p?name=fxnn/deadbox)

Access data and devices in your private network without Dynamic DNS, port opening etc.

Implementation is done as far as there's some spare time.

I'm very happy to hear your feedback and ideas. Simply file an issue!

## Problem

When a user wants to access his private data from the internet, he currently has two possibilities:
* Either he uses a public cloud, like Google Drive, Microsoft OneDrive etc., which forces him to _always_ store his data on foreign servers.
* Or he runs a private cloud like ownCloud, Nextcloud or proprietary software found on NAS devices, and allows to access them from the internet.
  This induces a severe security risk, as using security holes would allow an attacker to enter the users private network and data.
  The risk of such a security hole is high, as Dynamic DNS and private clouds rely on a broad interface and many different technologies.
  
It would be better, if the user wouldn't need to neither store his data on foreign servers, nor open his private network to the public internet.

## Approach
_Deadbox_ is an application that combines the concepts of peer to peer and 
message bus to establish communication between public and private networks.
It consists of two parts: _workers_ and _drops_.
* The worker runs on devices in your private network and connects to the
  drop,
  which is located on a machine somewhere publicy available.
* The user accesses the drop through a web-based UI,
  sending usecase based requests like "search for a file by name" or
  "upload a file".
  The drop stores the requests in a queue.
* The worker requests all items in his queue from the drop,
  processes them and sends the answers back to the drop.
* Finally, the user loads the results from the drop.

## Advantages
The user's experience is similar to that of a typical cloud application.
He uploads and downloads files or accesses other functionality through a comfortable web-based UI.

Still, the system provides more security, due to the following reasons.
* The private network provides no interface vulnerable to probing techniques like port scanning.
* The private network exposes a very small interface, protected through well-known authentication techniques.
  It is additionally protected through the server, who won't process unauthenticated requests.
* The publicy available server does not store private information over a longer time.
  It will remove all transmitted data after request completion or within a configurable timeout.
  Therefore, attacking the servers data storage will not provide access to the majority of information.
* Moreover, using public key cryptography, end-to-end encyption can be used.
  Then, the server can't see the data at all.

By combining worker and drop in one single binary and adding some basic routing,
a peer to peer network emerges, allowing to combine workers from different 
private networks using tree-like or even rather random structures.
One deadbox instance could serve as worker and drop at the same time, loading
requests from foreign drops, providing them to other workers.
This offers a high flexibility in retrieving data and sending commands from/to
private devices.

## Disadvantages
Compared to a centralized cloud-based storage solution,
the _deadbox_ suffers under limited availability and increased configuration 
effort.
* The service is only available when the devices in the private network are on
  and connected.
  However, this could be circumvented by letting the drop not only store the
  encrypted requests, but also the encrypted responses for some limited time.
  Of course, this is not feasible in all use cases.
* The user needs to install, configure and maintain instances on at least two
  different devices.
  This could be circumvented by providing a central cloud service, maybe on a
  subscription base.

## Challenges
* While the worker needs to open the connection to the drop,
  the drop must be able to push requests to the worker, so that the user won't 
  experience a high latency.
  Similarly, the worker must be able to push his responses back to the drop and
  thus to the user.
* When thinking about authentication, it is clear that the worker should not 
  need to trust the drop.
  Otherwise, in terms of security, the drop would be a single point of failure.
  Instead, it would be better if both could rely on a common third party for 
  establishing authentication, as it is possible with e.g. OpenID.

## Use cases and examples
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

## Specification

### Security

* Drop secures it's REST interface using TLS.
  If the workers regards the TLS certificate as untrusted,
  [Public Key Pinning](https://developer.mozilla.org/en-US/docs/Web/HTTP/Public_Key_Pinning) might be an option.
* Worker identifies itself using client certificates while connecting with TLS.
* Drop compares worker's certificate against whitelist as means of authorization.
* User encrypts requests using the worker's public key.
  User retrieves that public key from the Drop, which received it during registration.
* Worker encrypts responses using the user's public key, which is included in the request.

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
  received messages and to identify the worker's instance against a drop.
* Enabled receivers, with the appropriate configuration for each.
 * A file system receiver must be configured with the file system location, include/exclude patterns and allowed/disallowed operations.
 * An executing receiver must be configured with the command to execute when invoked.
 * Per receiver, it must be possible to enable/disable it for individual authenticated users.

### `drop.yml`

The file `drop.yml` contains the configuration for a drop.
It must configure at least the following aspects.

* Hostname and port to bind on.
* A whitelist of allowed workers, identified by a fingerprint of their key.

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

* `GET /worker/` allows users to retrieve all workers available to him.
* `POST /worker/` lets a worker register with the drop or update his 
  information.
  This endpoint is idempotent.
* `GET /worker/{workerId}/request` lets a worker retrieve all pending requests 
  targeted to him.
* `POST /worker/{workerId}/request` allows users to file a new 
  request targeted to an worker.
  This endpoint is idempotent.
* `POST /worker/{workerId}/response` is called by a worker to
  transport the response to a request.
  This endpoint is idempotent.
* `GET /worker/{workerId}/response/{requestId}` allows users to retrieve the
  response to a previously filed request, if already available.

As for a worker entry, at least the following information need to be contained.

* The worker's unique identification.
* A timestamp, after which the worker may be regarded as outdated and removed.
* The request types supported by the worker, given as URN.
* The worker's public key, used to encrypt the requests.

As for a queue entry, which is a request, at least the following information
need to be contained.

* The request's unique identification.
* A timestamp, after which the request may be regarded as outdated and removed.
* The request type, given as URN.
* The request payload.
* If we want to apply routing over multiple hops, we might also need a unique
  identifier of the originating drop, and possibly an ordered list of all
  drops that transported the request.

As for a response, at least the following information need to be contained.

* The original request's unique identification.
* A timestamp, after which the response may be regarded as outdated and removed.
* The response payload.
* If we want to apply routing over multiple hops, we might also need a unique
  identifier of the originating drop, and possibly an ordered list of all drops
  that are still to address.
