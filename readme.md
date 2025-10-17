# GnoVM Cosmos SDK Module

Cosmos SDK module for [GnoVM](https://github.com/gnolang/gno), a virtual machine for the Gno programming language.  
Cosmos SDK module scaffolded with [Ignite](https://ignite.com), a developer-friendly framework for building Cosmos SDK applications.

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

## Scaffolded with Ignite

This repo has been scaffolded with Ignite.
