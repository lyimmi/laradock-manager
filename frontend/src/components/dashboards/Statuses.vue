<template>
  <v-card
      class="mx-auto"
  >
    <v-card-title>
      Status
    </v-card-title>
    <v-card-text>
      <v-list disabled>
        <v-list-item>
          <v-list-item-icon>
            <v-icon color="success">check</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>docker</v-list-item-title>
          </v-list-item-content>
        </v-list-item>

        <v-list-item>
          <v-list-item-icon>
            <v-icon color="success">check</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>docker-compose</v-list-item-title>
          </v-list-item-content>
        </v-list-item>

        <v-list-item>
          <v-list-item-icon>
            <v-icon color="success">check</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>laradock path</v-list-item-title>
          </v-list-item-content>
        </v-list-item>

        <v-list-item>
          <v-list-item-icon>
            <v-icon color="success" v-if="appStatus.dotEnv">check</v-icon>
            <v-icon color="error" v-else>error</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>.env</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </v-list>
    </v-card-text>
  </v-card>
</template>

<script>
  import dockerCompose from '../../shared/dockerCompose'

  export default {
    name: 'statuses',
    mixins: [dockerCompose]
  }
</script>

<style scoped>

</style>