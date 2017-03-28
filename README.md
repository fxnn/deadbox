# Private Agent
Access data and devices in your private network without Dynamic DNS, port opening etc.

**Note:** This is a _concept_, still without implementation!
I'm very happy to hear your feedback and ideas. Simply file an issue!

## Problem

When a user wants to access his private data from the internet, he currently has two possibilities:
* Either he uses a public cloud, like Google Drive, Microsoft OneDrive etc., which forces him to _always_ store his data on foreign servers.
* Or he runs a private cloud like ownCloud, Nextcloud or proprietary software found on NAS devices, and allows to access them from the internet.
  This induces a severe security risk, as using security holes would allow an attacker to enter the users private network and data.
  The risk of such a security hole is high, as Dynamic DNS and private clouds rely on a broad interface and many different technologies.
  
It would be better, if the user wouldn't need to neither store his data on foreign servers, nor open his private network to the public internet.

## Approach
_Private Agent_ is an application that combines the concepts of peer to peer and message bus to establish communication between public and private networks.
It consists of two parts: _agents_ and _commanders_.
* The agent runs on devices in your private network and connects to the commander,
  which is located on a machine somewhere publicy available.
* The user accesses the commander through a web-based UI,
  sending usecase based requests like "search for a file by name" or "upload a file".
  The commander stores the requests in a queue.
* The client processes all requests from the queue and sends the answer back to the commander,
  which provides the results to the user.

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

By combining agent and commander in one single binary and adding some basic routing,
a peer to peer network emerges, allowing to combine agents from different private networks
using tree-like or even chaotic structures.
This offers a high flexibility in retrieving data and sending commands from/to private devices.

## Disadvantages
Compared to a centralized cloud-based storage solution,
the _private agent_ suffers under limited availability and increased configuration effort.
* The service is only available when the devices in the private network are on and connected.
  However, this could be circumvented by letting the commander cache encrypted responses to selected requests.
* The user needs to install, configure and maintain instances on at least two different devices.

## Challenges
* While the agent needs to open the connection to the controller,
  the controller must be able to push requests to the agent, so that the user won't expect a high latency.
  Similarly, the agent must also be able to push his responses back to the controller.
* When thinking about authentication, it is clear that the agent should not need to trust the controller.
  Otherwise, in terms of security, the controller would be a single point of failure.
  Instead, it would be better if both could rely on a common third party for establishing authentication,
  as it is possible with e.g. OpenID.

## Use cases and examples
The _private agent_ could be used in following scenarios.
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

### Command Line Interface

The main program is `privage-agent`, providing at least the following arguments.

* `-agent /path/to/agent.yml` makes the pocess act as an agent, connecting to a controller and waiting for commands.
* `-commander /path/to/commander.yml` makes the process act as a controller, opening a port providing a REST interface for both users and agents.
* `-auth /path/to/auth.yml` configures authentication for both agent and controller.

### `agent.yml`

The file `agent.yml` contains the configuration for an agent.
It must configure at least the following aspects.

* A textual identifier of the agent instance. Should be unique amongst all agents on the same controller.
* URL of one (or even multiple) controllers.
* A reference to a private (and possibly public) key file, used to decrypt received messages and to identify the agent's instance against a commander.
* Enabled receivers, with the appropriate configuration for each.
 * A file system receiver must be configured with the file system location, include/exclude patterns and allowed/disallowed operations.
 * An executing receiver must be configured with the command to execute when invoked.
 * Per receiver, it must be possible to enable/disable it for individual authenticated users.

### `commander.yml`

The file `commander.yml` contains the configuration for a commander.
It must configure at least the following aspects.

* Hostname and port to bind on.
* A whitelist of allowed agents, identified by a fingerprint of their key.

### `auth.yml`

The file `auth.yml` configures the authentication method to be used by both agents and commanders.
It must contain at least the following information.

* Authentication method to be used, e.g. htpasswd file or OAuth 2.0 server.
* Parameters for the authentication method, e.g. location of the passwd file or URL and access data to the OAuth 2.0 server.
 * Note, that commander and agent might need different parameters.
   While the commander must be able to establish user authentication, the agent must simply be able to consume a request.
