<template>
  <v-container fluid grid-list-md>
    <exec></exec>
    <v-layout>
      <v-flex xs8>
        <containers></containers>
      </v-flex>
      <v-flex xs4>
        <statuses></statuses>
      </v-flex>
    </v-layout>
  </v-container>
</template>

<script>
  import Statuses from './Statuses'
  import Containers from './Containers'
  import Exec from "../containers/Exec";

  export default {
    name: "index",
    components: {Exec, Containers, Statuses},
    data() {
      return {
      }
    },
  }
</script>

<style scoped>

</style>