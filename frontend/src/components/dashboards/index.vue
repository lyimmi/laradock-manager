<template>
  <v-container fluid grid-list-md>
    <v-layout>
      <v-flex xs8>
        <containers-list></containers-list>
      </v-flex>
      <v-flex xs4>
        <statuses></statuses>
      </v-flex>
    </v-layout>
  </v-container>
</template>

<script>
import Statuses from "./Statuses";
import ContainersList from "../containers/ConatinersList";

export default {
  name: "index",
  components: { ContainersList, Statuses },
  data() {
    return {};
  }
};
</script>

<style scoped>
</style>