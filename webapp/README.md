# deadbox/webapp

This is the webapp to access a Drop.

## Security Considerations

* The deadbox webapp is meant to process private, sensitive data. When not being secure, there'd be no reason to use this app, so this has highest priority.
* We must be careful who to trust.
  * Third party libraries and build tools [must not be trusted](https://hackernoon.com/im-harvesting-credit-card-numbers-and-passwords-from-your-site-here-s-how-9a8cb347c5b5) blindly.
  * Everyone must be able to verify the whole code of the webapp at any time. Therefore, the codebase must be as small as possible, but _not minified_.
  * Unfortunately, this forbids the use of transpilers.
    Yet, we [can use ES6](https://medium.freecodecamp.org/you-might-not-need-to-transpile-your-javascript-4d5e0a438ca) if we decide to drop support for older browsers.
  * _When_ using third party libraries and build tools, they must be very popular, thus having the highest probability of being reviewed by a lot of people.



## Docs

* [webpack](https://webpack.js.org/configuration/)
* [hyperapp](https://github.com/hyperapp/hyperapp/tree/master/docs), UI framework
  * [hyperapp/awesome](https://github.com/hyperapp/awesome#apps-and-boilerplates) is a curated list of projects built with Hyperapp.
* [bulma](https://bulma.io/documentation/overview/start/), CSS framework
