<template>
  <v-row>
    <v-col cols="12" class="pt-0">
      <container-list />
    </v-col>
  </v-row>
</template>

<script>
import ContainerList from "./containers/containerList";
export default {
  components: { ContainerList },
  data: () => {
    return {
    };
  },
  created() {
  },
  mounted() {
  }
};
</script>