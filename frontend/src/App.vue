<template>
  <v-app id="inspire">
    <!-- <app-menu></app-menu> -->
    <keep-alive>
      <div>
        <!-- menu start -->
        <v-navigation-drawer
          v-model="drawer"
          app
          clipped
          expand-on-hover
          fixed
          permanent
          mini-variant
          mini-variant-width="60"
        >
          <v-list dense>
            <v-list-item to="/home">
              <v-list-item-action>
                <v-icon>mdi-view-dashboard</v-icon>
              </v-list-item-action>
              <v-list-item-content>
                <v-list-item-title>Dashboard</v-list-item-title>
              </v-list-item-content>
            </v-list-item>
            <v-list-item to="/stats">
              <v-list-item-action>
                <v-icon>mdi-chart-bar</v-icon>
              </v-list-item-action>
              <v-list-item-content>
                <v-list-item-title>Stats</v-list-item-title>
              </v-list-item-content>
            </v-list-item>
            <v-list-item to="/settings">
              <v-list-item-action>
                <v-icon>mdi-cog</v-icon>
              </v-list-item-action>
              <v-list-item-content>
                <v-list-item-title>Settings</v-list-item-title>
              </v-list-item-content>
            </v-list-item>
          </v-list>
        </v-navigation-drawer>
      </div>
    </keep-alive>
    <!-- menu end -->
    <!-- content start -->
    <v-main>
      <v-container fluid>
        <transition name="fade" mode="out-in">
          <router-view></router-view>
        </transition>
      </v-container>
    </v-main>
    <!-- content end -->
    <v-snackbar
      v-for="(error, index) in errors"
      :key="index"
      v-model="error.state"
      color="error"
      multi-line
      :timeout="error.timeout"
      :absolute="false"
      top
    >
      {{ error.text }}
      <v-btn text @click="clearError(index)">Close</v-btn>
    </v-snackbar>
    <confirm-dialog ref="confirm" />
  </v-app>
</template>

<script>
import { mapGetters } from "vuex";
import DockerMixin from "./shared/dockerMixin";
import ErrorHandler from "./shared/errorHandlerMixin";
import ConfirmDialog from "./components/confirmDialog";

export default {
  name: "app",
  components: { ConfirmDialog },
  mixins: [DockerMixin, ErrorHandler],
  data: () => ({
    drawer: null,
    refreshCounter: 0,
  }),
  computed: {
    ...mapGetters("Settings", ["darkTheme"]),
  },
  created() {
    if (this.$router.history.current.path !== "/home") {
      this.$router.push("home");
    }
    this.setUpMasterErrorHandler();
  },
  mounted() {
    this.$vuetify.theme.dark = this.darkTheme;
    this.$root.$refs.confirm = this.$refs.confirm;
  },
  methods: {
    log(l) {
      console.log(l);
    },
    // ...mapActions("Settings", ["setLaradockPath"])
  },
  watch: {
    darkTheme(val) {
      this.$vuetify.theme.dark = val;
    },
  },
};
</script>

<style>
@import "./assets/css/main.css";
</style>

<style lang="scss" scoped>
.fade-enter-active,
.fade-leave-active {
  transition-duration: 0.3s;
  transition-property: height, opacity;
  transition-timing-function: ease;
  overflow: hidden;
}

.fade-enter,
.fade-leave-active {
  opacity: 0;
}

.slide-left-enter-active,
.slide-left-leave-active,
.slide-right-enter-active,
.slide-right-leave-active {
  transition-duration: 0.2s;
  transition-property: height, opacity, transform;
  transition-timing-function: cubic-bezier(0.55, 0, 0.1, 1);
  overflow: hidden;
}

@keyframes zoom {
  from {
    opacity: 0;
    transform: scale3d(0.3, 0.3, 0.3);
  }

  100% {
    opacity: 1;
  }
}

.refresh-counter {
  position: absolute;
  font-size: 9px;
  bottom: 0;
  transform: translate(-50%, 0);
  left: 30px;
}
</style>