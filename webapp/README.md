# deadbox/webapp

This is the webapp to access a Drop.

## Security Considerations

* The deadbox webapp is meant to process private, sensitive data. When not being secure, there'd be no reason to use this app, so this has highest priority.
* We must be careful who to trust.
  * Third party libraries and build tools [must not be trusted](https://hackernoon.com/im-harvesting-credit-card-numbers-and-passwords-from-your-site-here-s-how-9a8cb347c5b5) blindly.
  * Everyone must be able to verify the whole code of the webapp at any time. Therefore, the codebase must be as small as possible, but _not minified_.
  * _When_ using third party libraries and build tools, they must be very popular, thus having the highest probability of being reviewed by a lot of people.
* As many security relevant functions must be taken from third party libraries (as that's way more secure than implementing them ourselves). Candidates are:
  * [digitalbazaar/forge](https://github.com/digitalbazaar/forge), which is [fast](http://dominictarr.github.io/crypto-bench/), mature
    ([since 2009/2010](http://digitalbazaar.com/2010/07/20/javascript-tls-1/), >1,400 commits, 50 contributors) and gained a lot of attention
    (>2,000 stars on GitHub, >100,000 GitHub dependents).
    It implements a huge lot of crypto standards (amongst them AES, RSA, X509, SHA1+2, HMAC, PRNG, some encodings).
    On the downside, it uses quite a lot build tools / devDependencies.
  * [bitewiseshiftleft/sjcl](https://github.com/bitwiseshiftleft/sjcl) aka Stanford Javascript Crypto Library, which is also
    [fast](http://dominictarr.github.io/crypto-bench/),
    mature ([since 2009](http://bitwiseshiftleft.github.io/sjcl/acsac.pdf), >400 commits, and 45 contributors) and gained a lot of attention
    (>4,000 stars on GitHub, but only 1,000 GitHub dependents).
    It implements a fair set of crypto standards (amongst them AES, SHA1+2, ECC, HMAC, PRNG, some encodings).
    It has only two devDependencies (eslint and jsDoc), but brings its own build tools in the repo.

## Docs

* [create-react-app](https://github.com/facebookincubator/create-react-app/blob/master/packages/react-scripts/template/README.md), used to bootstrap this webapp
* [bulma](https://bulma.io/documentation/overview/start/), CSS framework
