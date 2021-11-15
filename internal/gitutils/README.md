# gitutils

This provides some functions for interacting with got repositories.

Why a custom implementation of git, instead of just relying on go-git?  Because go-git is sometimes [atrociously slow](https://github.com/go-git/go-git/issues/181), and isn't very fast in most other cases.  For example, our implementation of finding the current tag name on a particular local repo of mine takes 5ms, whereas the git-go implementation takes around 26ms.  Also, adding git-go approximately doubles the size of the executable.
