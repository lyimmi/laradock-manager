<template>
  <div>
    <v-navigation-drawer v-model="drawer" app clipped>
      <v-list dense>
        <v-list-item to="/">
          <v-list-item-action>
            <v-icon>dashboard</v-icon>
          </v-list-item-action>
          <v-list-item-content>
            <v-list-item-title>Dashboard</v-list-item-title>
          </v-list-item-content>
        </v-list-item><v-list-item to="/containers">
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
      <v-toolbar-title>Laradock manager</v-toolbar-title>

      <v-spacer></v-spacer>

      <v-menu
          left
          bottom
      >
        <template v-slot:activator="{ on }">
          <v-btn icon v-on="on">
            <v-icon>notifications</v-icon>
          </v-btn>
        </template>

        <v-list
            width="250px"
            disabled
        >
          <v-subheader>STATUS</v-subheader>
          <v-list-item
              v-for="n in 5"
              :key="n"
              @click="() => {}"
          >
            <v-list-item-title>Option {{ n }}</v-list-item-title>
          </v-list-item>
        </v-list>
      </v-menu>
      <v-btn
          icon
          @click="$root.$emit('refreshData')" prevent
      >
        <v-icon>refresh</v-icon>
      </v-btn>
    </v-app-bar>
  </div>
</template>

<script>
  import dockerCompose from '../../shared/dockerCompose'

  export default {
    name: 'app-menu',
    data() {
      return {
        drawer: null,
        refreshCounter: 0,
        refreshProgress: 0
      }
    },
    mixins: [dockerCompose],
    mounted() {
      setInterval(() => {
        this.refreshCounter++
        if (this.refreshCounter === 30) {
          this.$root.$emit('refreshData')
          this.refreshCounter = 0
        }
      }, 2000)
    },
    methods: {}
  }
</script>

<style scoped>
  .router-link-exact-active {
  }

  .router-link-active {
  }
</style>