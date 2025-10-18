# GnoVM Cosmos SDK Module

Cosmos SDK module for [GnoVM](https://github.com/gnolang/gno), a virtual machine for the Gno programming language.  
Cosmos SDK module scaffolded with [Ignite](https://ignite.com), a developer-friendly framework for building Cosmos SDK applications.

> [!WARNING]  
> This module is still in its alpha phase. Expect bugs and breaking changes.
> Please report any issues you encounter.
> Additionally, we currently rely on a fork of the GnoVM containing only this PR: https://github.com/gnolang/gno/pull/4852.

## Installation

To install the GnoVM module in your Cosmos SDK application, please follow the instructions below:

```bash
ignite s chain github.com/ignite/gnovm --minimal --no-module
ignite app install github.com/ignite/apps/gnovm@main
ignite gnovm add
ignite chain serve
```

The [Ignite GnoVM App](https://github.com/ignite/apps/tree/main/gnovm) simplifies the wiring of GnoVM into your Cosmos SDK application.

## Usage

### Start a localnet

The following command runs a single node locally, with the home located in
`~/.gnovm-localnet`.

```bash
make localnet-start
```

To target the node, use `./build/gnovmd --home ~/.gnovm-localnet ...`.

### Add Realm / Package

```bash
gnovmd tx gnovm add-package github.com/gno/examples/gno.land/r/demo/counter 5000stake --from alice
```

### Call Realm / Package

```bash
gnovmd tx gnovm call 1stake gno.land/r/demo/counter Increment --from alice
```

### Run Realm / Package

```bash
gnovmd tx gnovm run github.com/gno/examples/gno.land/r/demo/counter 5000stake --from alice
```

### Render Realm / Package

You can either query the `Render` function on the realm via cli:

```bash
gnovmd eval gno.land/r/demo/counter 'Render("")'
```

Or directly access its RPC endpoint on your node: <http://localhost:1317/ignite/gnovm/gnovm/v1/render/gno.land/r/demo/counter>

## Scaffolded with Ignite

This repo has been scaffolded with Ignite.
