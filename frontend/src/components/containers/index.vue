<template>
  <v-container fluid grid-list-md id="containersContainer">
    <containers-list></containers-list>
  </v-container>
</template>

<script>
import ContainersList from "./ConatinersList.vue";

export default {
  name: "index",
  components: { ContainersList }
};
</script>