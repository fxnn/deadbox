# Idea

## Problem

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

## Approach
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

## Advantages
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

## Disadvantages
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

## Challenges
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

## Related Work
The following projects are similar to _deadbox_.
* [Syncthing](https://syncthing.net) is an open source, multi platform synchronization solution.
  If you want to synchronize the same set of directories between the same set of devices for some time, use Syncthing! 
  However, if you'd rather have some publicly reachable endpoint for your private devices and use it to perform actions (of any kind) on your device, use _deadbox_.
* [ngrok](https://ngrok.com/) is a proprietary tunnel to your device.
  If you want to access any kind of webapp or TCP service (like Nextcloud), running on your device, itself running behind a firewall, from the public internet, use ngrok.
  However, if you'd rather have a fully open source solution, without the need to run multiple server components, or if you want to support non-interactive scenarios, where you want to enqueue actions (like sending files) while the target device is offline, use _deadbox_.
