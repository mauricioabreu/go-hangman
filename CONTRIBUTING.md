# How to contribute to this project

## Running the tests

In the root project, run the command `go test -v ./...` in order to execute the test suite.

## Did you find a bug?

* Ensure the bug was not reported already. The [issues index](https://github.com/mauricioabreu/go-hangman/issues) is your best friend <3
* If you did not find a bug report, please [open a new issue](https://github.com/mauricioabreu/go-hangman/issues/new) 

## Do you want to send a bug fix?

* Open a new pull request with the patch.
* Use clean and appropriate git commit messages.

## Do you want to send a new feature or change an existing one?

* Open a new issue so we can discuss the new feature.

## Do you want to add a new test?

I love tests but this project grew up a little fast and I had no control over code coverage.
If you want to send a test case, please go ahead. I would really appreciate the efforts.

## Commit messages

Please, follow the diagram below:

```
feature: Add hat wobble
 ^--^  ^------------^
 |     |
 |     +-> Summary in present tense
 |
 +-------> Type: chore, docs, feat, fix, refactor, style, or test
```

Start your git title with one of the available `tags`:

* chore: Add oyster build script
* docs: Explain hat wobble
* feat: Add beta sequence
* fix: Remove broken confirmation message
* refactor: Share logic between 4d3d3d3 and flarhgunnstow
* style: Convert tabs to spaces
* test: Ensure Tayne retains clothing
