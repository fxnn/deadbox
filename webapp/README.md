# deadbox/webapp

This is the webapp to access a Drop.

## General Security Considerations

The deadbox webapp is meant to process private, sensitive data. 
When not being secure, there'd be no reason to use this app, so this has highest priority.

We must be careful whom to trust.
This means, that

  * third party libraries and build tools 
    [must not be trusted](https://hackernoon.com/im-harvesting-credit-card-numbers-and-passwords-from-your-site-here-s-how-9a8cb347c5b5) blindly.
  * _when_ using third party libraries and build tools, they must be very popular, 
    thus having the highest probability of being reviewed by a lot of people.
  * the server delivering the webapp might be compromised, so everyone must be able
    to verify the whole code of the webapp at any time.

### JavaScript Cryptography

Though articles like [Javascript Cryptography Considered Harmful](https://www.nccgroup.trust/us/about-us/newsroom-and-events/blog/2011/august/javascript-cryptography-considered-harmful/)
state that cryptography implemented in webapps is insecure and should not be used,
we argue that

* it's established to trust the browser itself for security related applications, like
  for example online banking.
* it's possible to verify the webapp code executed on the client's machine.

### Verification

The webapp is compiled and delivered as one single HTML file, containing all required
CSS and JavaScript code.
It's easy for an advanced user

* to verify that no other resources than that single HTML file are loaded, and
* to compute and validate the HTML file's checksum.

Alternatively, in future, [subresource integrity](https://developer.mozilla.org/en-US/docs/Web/Security/Subresource_Integrity)
checking might be used.

## Docs

* [webpack](https://webpack.js.org/configuration/)
* [hyperapp](https://github.com/hyperapp/hyperapp/tree/master/docs), UI framework
  * [hyperapp/awesome](https://github.com/hyperapp/awesome#apps-and-boilerplates) is a curated list of projects built with Hyperapp.
* [bulma](https://bulma.io/documentation/overview/start/), CSS framework
