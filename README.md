# Wordsegment

This is a golang migrate version of [python-wordsegment](https://github.com/grantjenks/python-wordsegment).

> See [README.rst](README.rst) for more information of python-wordsegment.

# Usage
See [demo code](cmd/main.go) for more information. 

# Misc
## AI Migration Assistant
Even code migration via AI is a possible solution, but it DOES NOT work well enough.

- Most of the tests migrated by AI.
- Main module code is migrated by AI but many logic bugs exist.
- Though manual check and bugfix cost time, the AI migration help much too.

> Totally, it costs around 2 hours to complete all for the migration.

AI code migration works well for `strong-type -> weak-type`, `strong-type -> strong-type`. 

But it is not good for `weak-type -> strong-type`.

I think the major reason is leak understand of weak-type logic, especially the closure.