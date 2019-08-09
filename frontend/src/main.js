import 'babel-polyfill';
import Vue from "vue";

import VueRouter from 'vue-router'
import routes from './routes';

Vue.use(VueRouter);
const router = new VueRouter({
    mode: 'history',
    routes
});
router.replace({ path: '*', redirect: '/' });

import store from './store'


// Setup Vuetify
import vuetify from '@/plugins/vuetify' // path to vuetify export
import App from "./App.vue";

Vue.config.productionTip = false;
Vue.config.devtools = true;

import Bridge from "./wailsbridge";

Bridge.Start(() => {
    new Vue({
        store,
        vuetify,
        router,
        render: h => h(App)
    }).$mount("#app");
});
