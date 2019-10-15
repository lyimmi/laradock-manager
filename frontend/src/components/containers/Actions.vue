
<template>
  <div>
    <v-tooltip bottom v-if="item.state === 'Up'">
      <template v-slot:activator="{ on }">
        <v-btn small icon @click="logsContainer(item.name)" v-on="on">
          <v-icon :size="15">fa-bars</v-icon>
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
          <v-icon :size="14">fas fa-terminal</v-icon>
        </v-btn>
      </template>
      <span>Open container in terminal</span>
    </v-tooltip>

    <v-tooltip bottom v-if="item.state === 'Up'">
      <template v-slot:activator="{ on }">
        <v-btn small icon :size="14" @click="toggleContainer('stop', item.name)" v-on="on">
          <v-icon>stop</v-icon>
        </v-btn>
      </template>
      <span>Stop container</span>
    </v-tooltip>

    <v-tooltip bottom v-else-if="item.state !== 'DOWN'">
      <template v-slot:activator="{ on }">
        <v-btn small icon :size="14" @click="toggleContainer('start', item.name)" v-on="on">
          <v-icon>play_arrow</v-icon>
        </v-btn>
      </template>
      <span>Build the containers</span>
    </v-tooltip>

    <v-tooltip bottom v-else-if="item.state === 'DOWN'">
      <template v-slot:activator="{ on }">
        <v-btn small icon :size="14" @click="upContainer(item.name)" v-on="on">
          <v-icon>arrow_upward</v-icon>
        </v-btn>
      </template>
      <span>Up the containers</span>
    </v-tooltip>

    <v-menu offset-y>
      <template v-slot:activator="{ on }">
        <v-btn small icon v-on="on">
          <v-icon :size="14">build</v-icon>
        </v-btn>
      </template>
      <v-list>
        <v-list-item @click="buildContainer(item.name)">
          <v-list-item-title>Build with cache</v-list-item-title>
        </v-list-item>
        <v-list-item @click="buildContainer(item.name, true)">
          <v-list-item-title>Build without cache</v-list-item-title>
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