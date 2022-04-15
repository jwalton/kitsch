---
sidebar_position: 1
---

# Configuration

The top level of the configuration file may contain the following keys:

## timeout

This specified the default timeout for modules where the timeout is unspecified. Note that block modules ignore this - the default timeout for a block module is infinite. A value of 0 here will be an infinite timeout.  If not specified, this defaults to 500ms.

## extends

The name of another configuration file to extend (the parent configuration file). We load colors, prompt, and projects from the parent file first, then merge in any custom colors or project configuration from the current file. See [Configuration Merging](../configurationMerging.mdx).

## colors

A map of custom colors. Custom colors must start with a "$". See [Styles](../styles.mdx).

## projectTypes

An array of project types. See [Projects](../projects.mdx).

## prompt

The [module](./modules.mdx) to render as the prompt. Typically this would be a block module with multiple child modules.
