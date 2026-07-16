# Project Overview

CLI notetaking tool used in conjunction with Vim.

## Dev Environment Tips

* Use `mise` to install any runtimes
* Try to use the latest version of Go and other project dependencies

## Testing Instructions

* You can find the CI workflow in `.github/workflows/tests.yml`
* Only increase test coverage, never let it reduce

## Git Instructions

* Follow conventional commits
* Use the existing Git config, if not available create one with a generic user and email that can be updated later
* Included an `Assisted-by` trailer with the name of your model
* This repo uses trunk based development, make sure tests are passing before making a commit
* Trunk based development, ensure tests are passing before committing

## Other Information

* Ensure that all files are owned by uid 1000 and guid 1000 whilst root can still read, write, and execute where necessary
