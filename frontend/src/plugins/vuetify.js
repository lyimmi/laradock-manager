// You still need to register Vuetify itself
// src/plugins/vuetify.js
import "material-design-icons-iconfont";
import "@fortawesome/fontawesome-free/css/all.css"; // Ensure you are using css-loader

import Vue from "vue";
import Vuetify from "vuetify/lib";
import LRU from "lru-cache"

Vue.use(Vuetify);

const themeCache = new LRU({
  max: 10,
  maxAge: 1000 * 60 * 60, // 1 hour
})

export default new Vuetify({
  icons: {
    iconfont: "fa" // default - only for display purposes
  },
  theme: {
    options: {
      themeCache,
      minifyTheme: function (css) {
        return process.env.NODE_ENV === 'production'
          ? css.replace(/[\r\n|\r|\n]/g, '')
          : css
      },
    },
    themes: {
      dark: {
        primary: "#00bcd4",
        secondary: "#607d8b",
        accent: "#f44336",
        error: "#ffc107",
        warning: "#009688",
        info: "#2196f3",
        success: "#4caf50"
      }
    }
  }
});
