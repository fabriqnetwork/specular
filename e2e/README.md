# Specular E2E Testsuite

## Setup

Run if genesis changes.

```sh
./sbin/setup.sh
```

## Test

```sh
./sbin/build-container.sh
./sbin/test.sh
```

## Clean up

Run if `./sbin/test.sh` failed to clean up due to testing errors.

```sh
./sbin/clean.sh
```