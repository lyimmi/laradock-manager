import "babel-polyfill";
import * as Wails from "@wailsapp/runtime";
import Vue from "vue";
import store from "./store";
import vuetify from "@/plugins/vuetify"; // path to vuetify export
import VueRouter from "vue-router";
import routes from "./routes";
import App from "./App.vue";
import config from "./config/default.json";

Vue.use(VueRouter);
const router = new VueRouter({
  mode: "history",
  routes
});
router.replace({ path: "*", redirect: "/" });

Vue.config.productionTip = false;
Vue.config.devtools = true;
Vue.prototype.$config = config;

Wails.Init(() => {
  new Vue({
    store,
    vuetify,
    router,
    render: h => h(App)
  }).$mount("#app");
});
