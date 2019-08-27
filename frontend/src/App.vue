<template>
  <v-app id="inspire">
    <!-- <app-menu></app-menu> -->
    <keep-alive>
      <div>
        <!-- menu start -->
        <v-navigation-drawer v-model="drawer" app clipped>
          <v-list dense>
            <v-list-item to="/">
              <v-list-item-action>
                <v-icon>dashboard</v-icon>
              </v-list-item-action>
              <v-list-item-content>
                <v-list-item-title>Dashboard</v-list-item-title>
              </v-list-item-content>
            </v-list-item>
            <v-list-item to="/containers">
              <v-list-item-action>
                <v-icon>view_module</v-icon>
              </v-list-item-action>
              <v-list-item-content>
                <v-list-item-title>Containers</v-list-item-title>
              </v-list-item-content>
            </v-list-item>
            <v-list-item to="/settings">
              <v-list-item-action>
                <v-icon>settings</v-icon>
              </v-list-item-action>
              <v-list-item-content>
                <v-list-item-title>Settings</v-list-item-title>
              </v-list-item-content>
            </v-list-item>
          </v-list>
        </v-navigation-drawer>

        <v-app-bar app clipped-left dense fixed flat>
          <v-app-bar-nav-icon @click.stop="drawer = !drawer"></v-app-bar-nav-icon>

          <v-spacer></v-spacer>

          <v-btn icon @click="$root.$emit('refreshData')" prevent>
            <v-icon>refresh</v-icon>
            <div style="position: absolute;font-size: 9px;bottom: -7px;">{{60-refreshCounter}}</div>
          </v-btn>
        </v-app-bar>
      </div>
    </keep-alive>
    <!-- menu end -->
    <!-- content start -->
    <v-content>
      <transition
          name="fade"
          mode="out-in"
      >
        <router-view></router-view>
      </transition>
    </v-content>
    <!-- content end -->
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
  //import AppMenu from "./components/partials/app-menu";
  import dockerCompose from "./shared/dockerCompose"
  import {mapActions, mapGetters} from "vuex"

  export default {
    name: "app",
    //components: { AppMenu },
    mixins: [dockerCompose],
    mounted() {
      this.$vuetify.theme.dark = true
      this.waitForSettings(() => {
        this.$vuetify.theme.dark = this.darkTheme
      })
      //Data refresh
      setInterval(() => {
        this.refreshCounter++
        if (this.refreshCounter === 60) {
          this.$root.$emit("refreshData")
          this.refreshCounter = 0
        }
      }, 1000)
      this.$root.$on("resetRefreshConter", () => {
        this.refreshCounter = 0
      })

      //Router settings
      this.$router.push("home")
      this.$root.$on("showError", error => this.addError(error))
    },
    data: () => ({
      errors: [],
      drawer: null,
      refreshCounter: 0
    }),
    computed: {
      ...mapGetters("Settings", [
        "laradockPath",
        'darkTheme'
      ])
    },
    methods: {
      ...mapActions("Settings", ["setLaradockPath"]),
      openBrowser(url) {
        window.wails.Browser = {OpenURL: url}
      },
      addError(error) {
        this.errors.push({
          text: error.substring(0, 255),
          state: true
        })
        this.cleanErrors()
      },
      hideError(index) {
        this.errors[index].state = false
        setTimeout(() => {
          this.cleanErrors()
        }, 1500)
      },
      cleanErrors() {
        this.errors.forEach((v, i) => {
          if (!v.state) {
            this.errors.splice(i, 1)
          }
        })
      },
    }
  }
</script>

<style>
  .logo {
    width: 16em;
  }

  a {
    text-decoration: none;
  }
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
    opacity: 0
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
</style>