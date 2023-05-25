# Prometheus Exporter Release Process

The repo uses [goreleaser](https://goreleaser.com/) and [ko](https://ko.build/) to release the different artifacts.
To make a new release just create a new git tag, this will trigger a new Github action release [workflow](https://github.com/tfadeyi/sloth-simple-comments/blob/main/.github/workflows/release.yml).

```shell
git tag -a v0.1.0 -m "First release"
git push origin v0.1.0
```
