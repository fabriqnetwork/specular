const ReactRefreshWebpackPlugin = require('@pmmmwh/react-refresh-webpack-plugin');

module.exports = {
  webpack: {
    configure: webpackConfig => {
      const scopePluginIndex =
        webpackConfig.resolve.plugins.findIndex(
          ({constructor}) =>
            constructor &&
            constructor.name === 'ModuleScopePlugin',
        );
      webpackConfig.plugins.push(new ReactRefreshWebpackPlugin());
      webpackConfig.resolve.plugins.splice(scopePluginIndex, 1);
      return webpackConfig;
    },
  },
  babel: {
    presets: ['@babel/preset-react'],
    loaderOptions: (babelLoaderOptions, {env, paths}) => {
      return babelLoaderOptions;
    },
  },
};
