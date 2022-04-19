---
sidebar_position: 1
---

# Configuration

The top level of the configuration file may contain the following keys:

## timeout

This specified the default timeout for modules where the timeout is unspecified. Note that block modules ignore this - the default timeout for a block module is infinite. A value of 0 here will be an infinite timeout.  If not specified, this defaults to 500ms.

## scanTimeout

In many cases Kitsch needs to determine if a particular file exists, or a file with a particular extension exists.  In order to speed this process up, Kitsch reads the contents of a folder only once and caches this information internally.  If the disk we are reading from is a network disk, or contains an extremely large number of files, however, we could spend a long time reading in file information from that disk.  In order to make sure the prompt is rendered in a timely fashion, we limit the maximum time we spend reading the folder contents.

The `scanTimeout` is the time, in milliseconds, to scan the folder.  The default is 100ms.

For a concrete example, the `projects` module might want to show that your are in a JavaScript project if there is one or more ".js" files in the current folder.  However, if the current folder contains ten thousand files and is mounted over an SMB share, then it could take several seconds to scan the folder contents to see if there are any .js files present.  Instead, Kitsch will read as many files as it can before it hits the `scanTimeout`.  If it doesn't find any ".js" files before the timeout is reached, it will assume there aren't any.

## extends

The name of another configuration file to extend (the parent configuration file). We load colors, prompt, and projects from the parent file first, then merge in any custom colors or project configuration from the current file. See [Configuration Merging](../configurationMerging.mdx).

## colors

A map of custom colors. Custom colors must start with a "$". See [Styles](../styles.mdx).

## projectTypes

An array of project types. See [Projects](../projects.mdx).

## prompt

The [module](./modules.mdx) to render as the prompt. Typically this would be a block module with multiple child modules.
