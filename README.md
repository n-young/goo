# `Goo` - Minimal Sitebuilding using YAML and Markdown

Goo is ([yet another](https://jamstack.org/generators/)) static site generator built atop YAML and Markdown. I built it to make my life easier when building more data-centric websites; I hate setting up an entire framework to do something, and I'd rather have all of my config and specification in one file for very small sites.

## Installation

Run `go get github.com/n-young/goo` to install it. Then, you can run Goo from the command line.

## Usage

### CLI

TODO: Explain this LOL

### `site.yaml`

Most of Goo's functionality is built atop your main `site.yaml` file, which is where you specify the entire shape of the site. A sample is as follows:

```yaml
# Specify site name, output folder, static folder.
name: My Goo Site
output: build
static: static

# Global variable files.
global:
    secret_message: Hello, my darling

# Partials.
partials:
    header: partials/header.tmpl
    footer: partials/footer.tmpl

# Pages
pages:
    - title: Home
      path: /
      template: templates/home.tmpl
      data:
          home: data/home.yaml

    - title: Store
      path: /store
      template: templates/store.tmpl

# Collections.
collections:
    - title: Blog Posts
      path: /blog
      template: templates/blog.tmpl
      posts: posts
```

Let's walk through this example.

`name` is, right now, a purely decorative line. It helps the reader be sure which website this file pertains to.

`output` is the folder to which the generated site should be exported to. In the above example, the site will be exported into `build/`.

`static` is the folder containing static assets, such as CSS, Javascript, or images. It will be copied into the build folder.

`global` can contain a map in any valid YAML format containing data you'd like to access in a template or a collection (more on these later). Data specified here can be accessed through the `global` object.

`partials` should be a map from partial tags to the files that contain those partials. Partials can be directly injected into a template, but (as of now) not into other partials.

`pages` should be a sequence of pages, each having a `title`, `path`, `template`, and optional `data` fields. The `title` is accessible through the `title` object. The `path` is where the page will be accessible, treating `output` as the base directory. The `template` is the base template that this path should build from. And `data` is a list of data points accessible from the template. More on data later.

`collections` should be a sequence of collections, each having a `title`, `path`, `template`, `posts`, and optional `data`. All of these are the same as for `pages`, except for `posts`. `posts` should be a directory which is filled the Markdown files (posts) you wish to convert to HTML. Each post has front matter than can be accessed from the template. Each post's content will be injected into the `content` field of the template.

While the above gives an idea of what each part of the `site.yaml` file means, it is through Templates and Posts that these fields come alive.


### Templates

A template (e.g. `templates/index.tmpl`) is an HTML-like file that can have data and content injected into it. The structure of a template is just like HTML, except that you can inject data using something of the pattern:

```
{{ <action> <payload> }}
```

There are a number of potential actions, details below.

#### Title
`title` takes no payload, and injects the current page or post title.

#### Partial
`partial <partial_name>` takes a name of a partial, specified in the `site.yaml`, and injects its contents. Partials cannot be injected into other partials (no nested evaluation).


#### TODO: Data model

- title
- global
- post
- defined paths
- examples

#### Data
`data` takes the path of a data point and injects it with no special formatting. For example, if we had a `data.yaml` file like:

```yaml
dog:
    color:
        hue: "red"
```

Where our `site.yaml` file specified that:

```yaml
...
data:
    dogs: data.yaml
...
```

Inside of the template that has access to this data, we could inject the "hue" field by using:

```
{{ data dogs.dog.color.hue }}
```

#### Template
`template` is the same as data, just with more complex string replacement. An example illustrates best:

```
{{ template
<tr>
    Hue: ${dogs.dog.color.hue}
    Bark: ${dogs.dog.bark}
</tr>
}}
```

#### Loop
`loop` allows you to iterate over sequential data in YAML, applying a template to each data point. To loop, you specify `loop`, then the data to be looped over, then the template. For example:

```yaml
dogs:
    - name: Fido
      age: "15"
    - name: George
      age: old
    - name: Robin
      age: 10?
```

Could be looped over using:
```
{{ loop
<tr>
    Name: ${dogs.dogs.name}
    Age: ${dogs.dogs.age}
</tr>
}}
```

This would generate two table rows, making it ideal for highly repetitive data.

TODO: Nested loops, also need to change parsing, then.

### Content
`content` is used exclusively in a collection, and is the site where the main content specified in a Markdown file will be injected. Information on parsing options is in the Posts section.

### Posts

Markdown files used in a Collection will each be converted to HTML using Goldmark and written out as a separate page. This is ideal for generating structured content like blog posts. The structure of a post file is as follow:

```md
---
title: My Post
name: Arthur
---

# Hello, world!
```

TODO: Accessing post data (in the post object)

TODO: Clarify the below as other features and give examples.

Our parsing supports unsafe HTML, attributes, :joy:-style emojis, inline LaTex, syntax highlighting, and the entire GitHub Flavored Markdown spec.

#### Note on using MathJax

If you decide to use inline math, link the following in the footer:
```html
<script src="https://polyfill.io/v3/polyfill.min.js?features=es6"></script>
<script id="MathJax-script" async src="https://cdn.jsdelivr.net/npm/mathjax@3/es5/tex-mml-chtml.js"></script>
```


## Reference

Note that the `sample/` directory of this repository has examples of most of the functionality offered by Goo. If you'd like to see new functionality, submit an Issue.


## Contributing

Feel free to make pull requests and contribute to this project! A list of potential contributions can be found in the Issues tab. 
