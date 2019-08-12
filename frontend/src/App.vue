<template>
  <v-app id="inspire">
    <app-menu></app-menu>

    <v-content>
      <transition>
        <router-view></router-view>
      </transition>
    </v-content>
    <v-snackbar
      v-for="(error, index) in errors"
      :key="index"
      v-model="errors[index].state"
      color="error"
      multi-line
      :timeout="15000"
      :absolute="false"
      top
    >
      {{ error.text }}
      <v-btn text @click="hideError(index)">Close</v-btn>
    </v-snackbar>
  </v-app>
</template>

<script>
import AppMenu from "./components/partials/app-menu";
import dockerCompose from "./shared/dockerCompose";
import { mapActions } from "vuex";
import { mapGetters } from "vuex";

export default {
  name: "app",
  components: { AppMenu },
  mixins: [dockerCompose],
  created() {
    this.$vuetify.theme.dark = true;
  },
  mounted() {
    this.$router.push("home");
    this.$root.$on("showError", error => this.addError(error));
  },
  data: () => ({
    errors: []
  }),
  computed: {
    ...mapGetters("Settings", ["laradockPath"])
  },
  methods: {
    ...mapActions("Settings", ["setLaradockPath"]),
    openBrowser(url) {
      window.wails.Browser = { OpenURL: url };
    },
    addError(error) {
      this.errors.push({
        text: error.substring(0, 255),
        state: true
      });
      this.cleanErrors();
    },
    hideError(index) {
      this.errors[index].state = false;
      setTimeout(() => {
        this.cleanErrors();
      }, 1500);
    },
    cleanErrors() {
      this.errors.forEach((v, i) => {
        if (!v.state) {
          this.errors.splice(i, 1);
        }
      });
    }
  }
};
</script>

<style>
.logo {
  width: 16em;
}

a {
  text-decoration: none;
}
</style>