
<template>
  <div>
    <v-btn
      icon
      small
      v-if="item.state === 'Up'"
      class="ma-2"
      dark
      v-on:click="$root.$emit('execContiner', item.name)"
    >
      <v-icon :size="14">fas fa-terminal</v-icon>
    </v-btn>

    <v-btn
      small
      icon
      :size="14"
      v-if="item.state === 'Up'"
      @click="toggleContainer('stop', item.name, true)"
    >
      <v-icon>stop</v-icon>
    </v-btn>

    <v-btn
      small
      icon
      :size="14"
      v-else-if="item.state !== 'DOWN'"
      @click="toggleContainer('start', item.name, true)"
    >
      <v-icon>play_arrow</v-icon>
    </v-btn>

    <v-btn small icon :size="14" v-else-if="item.state === 'DOWN'" @click="upContainer(item.name)">
      <v-icon>arrow_upward</v-icon>
    </v-btn>

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