<template>
  <v-card>
    <v-tabs v-model="tab" color="info" class="settings-tabs">
      <v-tabs-slider></v-tabs-slider>

      <v-tab v-for="(tab, i) in tabs" :key="i" :href="`#tab-${i}`">{{ tab }}</v-tab>

      <v-tab-item value="tab-0">
        <app />
      </v-tab-item>

      <v-tab-item value="tab-1">
        <v-card flat tile>
          <v-card-text>excample 2</v-card-text>
        </v-card>
      </v-tab-item>
    </v-tabs>
  </v-card>
</template>

<script>
import app from "./app";
import Docker from "@/shared/dockerMixin";
export default {
  name: "settings",
  mixins: [Docker],
  components: { app },
  data() {
    return { tab: null, tabs: ["App"] };
  },
  mounted() {},
};
</script>

<style scoped>
</style>