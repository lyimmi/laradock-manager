<template>
  <v-container fluid grid-list-md>
    <v-layout>
      <v-flex xs12>

      </v-flex>
    </v-layout>
  </v-container>
</template>

<script>
  import dockerCompose from '../../shared/dockerCompose'

  export default {
    name: 'index',
    mixins:[dockerCompose],
    mounted() {
      this.getAvailableContainers()
    }
  }
</script>

<style scoped>

</style>