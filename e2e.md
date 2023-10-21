# E2E test-suite

## Setup

Run if genesis changes.

```sh
./sbin/setup.sh
```

## Test

```sh
./sbin/run_e2e_tests.sh
```

## Cleanup

Run if `./sbin/run_e2e_tests.sh` failed to clean up due to testing errors.

```sh
./sbin/clean.sh
```
