import Vue from 'vue';
import Router from 'vue-router';
import Meta from 'vue-meta';
import Home from './views/Home.vue';
import Watch from './views/Watch.vue';
import Search from './views/Search.vue';
import Upload from './views/Upload.vue';
import Edit from './views/Edit.vue';

Vue.use(Router);
Vue.use(Meta);

export default new Router({
  mode: 'history',
  base: process.env.BASE_URL,
  scrollBehavior: (to, from, savedPosition) => {
    if (savedPosition) {
      return savedPosition;
    }

    return { x: 0, y: 0 };
  },
  routes: [
    {
      path: '/',
      name: 'home',
      component: Home,
    },
    {
      path: '/watch/:id',
      name: 'watch',
      component: Watch,
      props: true,
    },
    {
      path: '/edit/:id',
      name: 'edit',
      component: Edit,
      props: true,
    },
    {
      path: '/upload',
      name: 'upload',
      component: Upload,
    },
    {
      path: '/search',
      name: 'search',
      component: Search,
      props: route => ({
        page: parseInt(route.query.page, 10) || 1,
        mode: route.query.mode || 'text',
        text: route.query.text,
        tags: route.query.tags,
      }),
    },
    {
      path: '*',
      name: 'not-found',
      component: () => import(/* webpackChunkName: "not-found" */ './views/NotFound.vue'),
    },
  ],
});
