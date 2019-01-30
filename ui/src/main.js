import Vue from 'vue';
import App from './App.vue';
import router from './router';
import store from './store';
import 'semantic-ui-css/semantic.css';
import './global.css';

Vue.config.productionTip = false;

console.log(`
# AlbinoDrought/creamy-videos

Repo: https://github.com/AlbinoDrought/creamy-videos

Source: https://${window.location.host}/source.tar.gz
`);

new Vue({
  router,
  store,
  render: h => h(App),
}).$mount('#app');
