# cp_parser
Parsing test cases from competitive programming contest websites.

## Installation

You can either:
- Download the binary from [Release](https://github.com/thallium/cp_parse/releases)
- Build by yourself:

`go install github.com/thallium/cp_parse@latest`

This will install the `cp_parse` executable with the library and its dependencies.

## Server and Chrome Extension (Experimental)

If viewing the problem or contest requires login (cannot get the content with a http ger request), use `cp_parse server` to turn this application into a server that can receive the html of the problem or contest(in development) page.  The chrome extension that comes with it can get the html and send it to the server when click on the icon.

## Usage
```
Usage:
  cp_parse [command]

Available Commands:
  atc         Parse problems/contests from atcoder.jp
  cf          Parse problems/contests from codeforces.com
  help        Help about any command
  kattis      Parse problems/contests from open.kattis.com
  server      Open a server that can receive html and parse it
 ```
