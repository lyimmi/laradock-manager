// You still need to register Vuetify itself
// src/plugins/vuetify.js
import "material-design-icons-iconfont";
import "@fortawesome/fontawesome-free/css/all.css"; // Ensure you are using css-loader

import Vue from "vue";
import Vuetify from "vuetify/lib";

Vue.use(Vuetify);

export default new Vuetify({
  icons: {
    iconfont: "fa" // default - only for display purposes
  },
  theme: {
    themes: {
      dark: {
        primary: "#00bcd4",
        secondary: "#cdeaff",
        accent: "#009688",
        error: "#f44336",
        warning: "#ffc107",
        info: "#03a9f4",
        success: "#4caf50"
      }
    }
  }
});
