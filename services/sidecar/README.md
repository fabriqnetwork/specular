# Specular Sidecar

## License

This module is licensed under [Apache License 2.0](https://www.apache.org/licenses/LICENSE-2.0); see `LICENSE` for details.

## Development

### Configuration

Sidecar supports configuration injection via [Viper](https://github.com/spf13/viper) (`.json`, `.toml`, `.yaml`, `.env`
and environment variables) autoload from the current working directory and `config` folder.

It additionally uses [CLI](https://github.com/urfave/cli) for the command line arguments parser.

Currently, config is split between `Config` and `SystemConfig`, consolidation of the configuration is WIP.

### Using Wire

Path `internal/service/di` contains providers injected using `inject.go`.
Wire will automatically match demand of the arguments in the constructors injected
with the provided constructors of objects. The initiation order will be determined automatically.
After generation `wire_gen.go` file will be generated containing
`func SetupApplication() (*Application, func(), error)` for `main` and
`SetupApplicationForIntegrationTests(cfg *config.Config) (*TestApplication, func(), error)`
for the testing purposes.

After modifying providers, you need to call `make wire-generate` in order to regenerate.

Object instances are directly assigned to `Application` object fields and ready to access within
the `Run()` entry point function.

For details please refer to [Wire's official documentation](https://github.com/google/wire).
