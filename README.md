# UwU Note

First of all, I really hate making up names, so `UwU Note` is all you get.

## What's it?

It's a simple note-taking app. Notes stick to your desktop, you can create new ones and delete existing ones.
Every note has it's own window, which saves it's coordinates and size.

## The Future

For the future i have planned some stuff:

* Cloud synchronization (Different Services)
* An android app
* Markdown support

## Installation

As this is a work in progress, there are currently no installation instructions, just build the project yourself instead.
Theretically you should be able to build it on every platform supporting GTk3+.

Official binaries are gonna follow at some point, maybe even a snap or flatpak.

## Building

To build the app you simply run these commands:

```bash
cd $GOPATH
go get github.com/bios-marcel/uwunote
cd src/github.com/bios-marcel/uwunote
go get ./...
go build
```

You might stil need some dependencies:

TODO - List dependencies

## Contribute

Feel free to create issues or create pull requests.