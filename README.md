# Private Agent
Access data and devices in your private network without Dynamic DNS, port opening etc.

**Note:** This is a _concept_, still without implementation!
I'm very happy to hear your feedback and ideas. Simply file an issue!

## Idea

When a user wants to access his private data from the internet, he currently has two possibilities:
* Either he uses a public cloud, like Google Drive, Microsoft OneDrive etc., which forces him to _always_ store his data on foreign servers.
* Or he runs a private cloud like ownCloud, Nextcloud or proprietary software found on NAS devices, and allows to access them from the internet.
  This induces a severe security risk, as using security holes would allow an attacker to enter the users private network and data.
  The risk of such a security hole is high, as Dynamic DNS and private clouds rely on a broad interface and many different technologies.
  
It would be better, if the user wouldn't need to neither store his data on foreign servers, nor open his private network to the public internet.
_Private Agent_ is a message bus that establishes communication between internet and private network.
It consists of a client and a server part.
The client runs on devices in your private network and connects to the server, which is located on a machine somewhere publicy available.

The user accesses the server, sending usecase based requests like "search for a file by name" or "upload a file".
The server stores the requests (if necessary), so that the client can pull the request, process it and send the answer back to the server, which provides it to the user.

The server's user interface provides a user experience similar to one of the well known cloud applications.
Still, the system provides more security, due to the following reasons.
* The private network provides no interface vulnerable to probing techniques like port scanning.
* The server does not store private information over a longer time.
  It will remove all transmitted data after request completion or within a timeout.
  Therefore, attacking the servers data storage will not provide access to the majority of information.
* The private network exposes a very small interface, protected through well-known authentication techniques,
  additionally protected through the server.

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
