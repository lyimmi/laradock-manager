<template>
  <v-container fluid grid-list-md>
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

  export default {
    name: "index",
    components: {Containers, Statuses},
    data() {
      return {
      }
    },
  }
</script>

<style scoped>

</style>