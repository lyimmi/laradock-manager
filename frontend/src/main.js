import Vue from "vue";
import store from "./store";
import vuetify from "@/plugins/vuetify"; // path to vuetify export
import { preset } from 'vue-cli-plugin-vuetify-preset-rally/preset'
import VueRouter from "vue-router";
import routes from "./routes";
import App from "./App.vue";
import config from "./config/default.json";
import Ticker from "./shared/classes/ticker"

Vue.use(VueRouter);
const router = new VueRouter({
  mode: "history",
  routes
});
router.replace({ path: "*", redirect: "/" });

Vue.config.productionTip = false;
Vue.config.devtools = true;
Vue.prototype.$config = config;
Vue.prototype.$ticker = new Ticker(250);

import Wails from '@wailsapp/runtime';

Wails.Init(() => {
  new Vue({
    preset,
    store,
    vuetify,
    router,
    render: h => h(App)
  }).$mount("#app");
});
