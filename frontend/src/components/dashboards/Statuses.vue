<template>
  <v-card class="mx-auto">
    <v-card-title>Status</v-card-title>
    <v-card-text>
      <v-list disabled dense>
        <v-list-item dense>
          <v-list-item-icon small>
            <v-icon color="error" v-if="dockerVersion === ''">mdi-alert-circle-outline</v-icon>
            <v-icon color="success" v-else>mdi-check</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>
              Docker
              <br />
              <small>{{dockerVersion}}</small>
            </v-list-item-title>
          </v-list-item-content>
        </v-list-item>

        <v-list-item>
          <v-list-item-icon>
            <v-icon color="error" v-if="dockerComposeVersion === ''">mdi-alert-circle-outline</v-icon>
            <v-icon color="success" v-else>mdi-check</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>
              Docker Compose
              <br />
              <small>{{dockerComposeVersion}}</small>
            </v-list-item-title>
          </v-list-item-content>
        </v-list-item>

        <v-list-item>
          <v-list-item-icon>
            <v-icon color="error" v-if="laradockPath === ''">mdi-alert-circle-outline</v-icon>
            <v-icon color="success" v-else>mdi-check</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>
              laradock path
              <br />
              <small>{{laradockPath}}</small>
            </v-list-item-title>
          </v-list-item-content>
        </v-list-item>

        <v-list-item>
          <v-list-item-icon>
            <v-icon color="error" v-if="appStatus.dotEnv">mdi-alert-circle-outline</v-icon>
            <v-icon color="success" v-else>mdi-check</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>.env</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
        
        <v-list-item>
          <v-list-item-icon>
            <v-icon color="success">mdi-check</v-icon>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title>{{containerPrefix}}</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </v-list>
    </v-card-text>
  </v-card>
</template>

<script>
import dockerCompose from "../../shared/dockerCompose";
import { mapGetters } from "vuex";

export default {
  name: "statuses",
  mixins: [dockerCompose],
  computed:{
    ...mapGetters("Settings", ["containerPrefix"]),
  },
  mounted() {
    this.checkDotEnv();
    this.checkDockerVersion();
    this.checkDockerComposeVersion();
  }
};
</script>

<style scoped>
</style>