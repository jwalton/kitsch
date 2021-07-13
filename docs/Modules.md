# Modules

## Global Module Configuration

* `style` - Style to apply to this module.  May be overridden by the module in certain situations (e.g. the "username" module will use `rootStyle` instead of `style` if the current user is root).
* `template` - Specify a golang template to override the rendering of this module.

## hostname Module

### Template Variables

* `.hostname` - The name of the host.
* `.isRemoteSession` - True if the user is currently logged in via SSH (SSH_CONNECTION, SSH_CLIENT, or SSH_TTY environment variables are set).

### Configuration

* `always` - If true, hostname will always be shown.  If false, hostname will only be shown
if `.isRemoteSession` is true.