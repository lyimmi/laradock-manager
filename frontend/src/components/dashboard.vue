<template>
  <v-container></v-container>
</template>

<script>
export default {
  data: () => {
    return {
      ojo: 0
    };
  },
  created() {
    const id = this.$ticker.setTicker(() => {
      this.ojo++;

      if (this.ojo === 5) {
        this.$ticker.clearTicker(id);
      }
      console.log("ojojo");
    }, 5);
  },
  mounted(){
    this.$ticker.nextTick(() => {
      console.log("ájájájá")
    })
  }
};
</script>