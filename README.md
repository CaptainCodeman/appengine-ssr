# appengine-ssr

Http middleware to call Puppeteer to Server Side Render a
page and cache the results in a cache (such as Memcache)

## Explanation

Say you have the latest and greatest front-end web client
technology for your new site. Unfortunately, Googlebot is
based on Chrome 41 so it probably doesn't understand any
of it. Oh dear. ShadowDOM? WebComponents? What they?

Here's how your sweet UI appears to Google (using the
"Fetch as Google" webmaster tool):
![Fetch as Google with ShadowDom](https://raw.githubusercontent.com/captaincodeman/appengine-ssr/master/examples/fag-shadow.png)

Yeah, your SEO isn't even SE, let alone with added O ...

One solution is to use Server-Side Rendering using the Chrome
headless browser via Puppeteer which is what this middleware
is intended to help with. It means when the Googlebot (or 
any other bot you configure) requests your site it can be
converted into ye-olde HTML like a static site. Now your
site can be rendered and understood. SEO glory awaits ...

![Fetch as Google with ShadyDOM](https://raw.githubusercontent.com/captaincodeman/appengine-ssr/master/examples/fag-ssr.png)

## How it works

The middlware checks if the request is coming from any of
the configured bots by parsing the User-Agent string of the
request. If it isn't, it's process as normal and the code
does nothing.

If it is, the middleware first checks if the URL has already
been server-side rendered and cached (because it can be a bit
expensive to do) and will serve that. If it isn't cached then
it will call the Puppeteer service to render the page but
adding an extra `headless` QueryString parameter to the URL.

The page should use this `headless` parameter to configure
the page to render with `ShadyDOM`. This will create a regular
DOM tree which can then be returned to the middleware and
returned to the caller, added to the cache first.

## Demo

A typical client-side rendered page is available at:
[http://ssr-dot-captain-codeman.appspot.com/](http://ssr-dot-captain-codeman.appspot.com/)

Here's how the DOM looks when rendered in a modern browser
such as Chrome. Note the shadow-root:

![ShadowDOM rendered](https://raw.githubusercontent.com/captaincodeman/appengine-ssr/master/examples/dom-shadow.png)

To view the page as it would look when Server-Side Rendered
you can override the User-Agent parsing by adding `?ssr` to
the URL:

[http://ssr-dot-captain-codeman.appspot.com/?ssr](http://ssr-dot-captain-codeman.appspot.com/?ssr)

Now the content is rendered as a regular DOM tree which
can be understood by any bot. If you generate metadata in
your front-end framework, this will enable services such as
Facebook and Twitter to render it:

![Server Side Rendered](https://raw.githubusercontent.com/captaincodeman/appengine-ssr/master/examples/dom-ssr.png)

## Installation

Install using `go get`

    go get -i github.com/captaincodeman/appengine-ssr

## Usage

TODO ...