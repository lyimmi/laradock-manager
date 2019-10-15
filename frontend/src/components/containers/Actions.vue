
<template>
  <div>
    <v-tooltip bottom v-if="item.state === 'Up'">
      <template v-slot:activator="{ on }">
        <v-btn small icon @click="logsContainer(item.name)" v-on="on">
          <v-icon>mdi-format-list-checks</v-icon>
        </v-btn>
      </template>
      <span>Show container logs</span>
    </v-tooltip>

    <v-tooltip bottom>
      <template v-slot:activator="{ on }">
        <v-btn
          icon
          small
          v-if="item.state === 'Up'"
          class="ma-2"
          v-on:click="$root.$emit('execContiner', item.name)"
          v-on="on"
        >
          <v-icon>mdi-console-line</v-icon>
        </v-btn>
      </template>
      <span>Open container in terminal</span>
    </v-tooltip>

    <v-tooltip bottom v-if="item.state === 'Up'">
      <template v-slot:activator="{ on }">
        <v-btn small icon @click="toggleContainer('stop', item.name)" v-on="on">
          <v-icon>mdi-stop</v-icon>
        </v-btn>
      </template>
      <span>Stop container</span>
    </v-tooltip>

    <v-tooltip bottom v-else-if="item.state !== 'DOWN'">
      <template v-slot:activator="{ on }">
        <v-btn small icon @click="toggleContainer('start', item.name)" v-on="on">
          <v-icon>mdi-play</v-icon>
        </v-btn>
      </template>
      <span>Build the containers</span>
    </v-tooltip>

    <v-tooltip bottom v-else-if="item.state === 'DOWN'">
      <template v-slot:activator="{ on }">
        <v-btn small icon @click="upContainer(item.name)" v-on="on">
          <v-icon>mdi-arrow-up-bold</v-icon>
        </v-btn>
      </template>
      <span>Up the containers</span>
    </v-tooltip>

    <v-menu offset-y>
      <template v-slot:activator="{ on }">
        <v-btn small icon v-on="on">
          <v-icon>mdi-progress-wrench</v-icon>
        </v-btn>
      </template>
      <v-list dense>
        <v-list-item @click="buildContainer(item.name)">
          <v-list-item-content>Build with cache</v-list-item-content>
        </v-list-item>
        <v-list-item @click="buildContainer(item.name, true)">
          <v-list-item-content>Build without cache</v-list-item-content>
        </v-list-item>
      </v-list>
    </v-menu>
  </div>
</template>
<script>
import dockerCompose from "../../shared/dockerCompose";

export default {
  name: "actions",
  mixins: [dockerCompose],
  props: ["item"],
  mounted() {}
};
</script>

<style scoped>
</style>