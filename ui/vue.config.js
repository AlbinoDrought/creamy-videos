const fallbackRoutes = [
  /^\/$/,
  /^\/watch\/.+$/,
  /^\/search$/,
];

if (!process.env.VUE_APP_READ_ONLY) {
  fallbackRoutes.push(/^\/edit\/.+$/);
  fallbackRoutes.push(/^\/upload$/);
}

module.exports = {
  devServer: {
    proxy: {
      '^/api': { target: 'http://localhost:3000/' },
      '^/static/videos': { target: 'http://localhost:3000/' },
    },
  },
  pwa: {
    name: 'Creamy Videos',
    themeColor: '#1b1b1b',
    msTileColor: '#000000',
    appleMobileWebAppCapable: 'yes',
    appleMobileWebAppStatusBarStyle: 'black',
    iconPaths: {
      appleTouchIcon: 'icons/icon-152x152.png',
      msTileImage: 'icons/icon-144x144.png'
    },
    workboxOptions: {
      skipWaiting: true,
      navigateFallback: '/index.html',
      navigateFallbackWhitelist: fallbackRoutes,
      runtimeCaching: [{
        urlPattern: /static/,
        handler: 'cacheFirst',
      }, {
        urlPattern: /api/,
        handler: 'networkFirst',
        options: {
          networkTimeoutSeconds: 5,
          cacheableResponse: {
            statuses: [200],
          },
        },
      }],
    },
  }
};
