<template>
  <v-container fluid grid-list-md>
    <v-layout>
      <v-flex xs12>
        <v-card
            class="mx-auto"
            :loading="loading"
        >
          <v-card-title>
            Containers
          </v-card-title>
          <v-card-text>
            <v-simple-table
                fixed-header v-if="containers.length > 0"
            >
              <thead>
              <tr>
                <th class="text-left">Name</th>
                <th class="text-center">State</th>
                <th>Action</th>
              </tr>
              </thead>
              <tbody>
              <tr
                  v-for="item in availableContainers" :key="item.name"
              >
                <td>
                  {{ item.name }}
                </td>
                <td class="text-center">
                  <v-chip :color="stateColor(item.state)">
                    {{item.state}}
                  </v-chip>
                </td>
                <td>
                  <v-btn
                      small
                      color="warning"
                      v-if="item.state === 'Up'"
                      @click="toggleContainer('stop', item.code)"
                  >
                    <v-icon>stop</v-icon>
                  </v-btn>
                  <v-btn
                      small
                      color="success"
                      v-else-if="item.state !== 'DOWN'"
                      @click="toggleContainer('start', item.code)"
                  >
                    <v-icon>play_arrow</v-icon>
                  </v-btn>
                  <v-btn
                      class="ml-2"
                      small
                      color="info"
                      @click="toggleContainer('build', item.code)"
                  >
                    <v-icon>build</v-icon>
                  </v-btn>
                </td>
              </tr>
              </tbody>
            </v-simple-table>
          </v-card-text>
        </v-card>
      </v-flex>
    </v-layout>
  </v-container>
</template>

<script>
  import dockerCompose from "../../shared/dockerCompose"

  export default {
    name: "index",
    mixins: [dockerCompose],
    data() {
      return {
        loading: true
      }
    },
    mounted() {
      this.getAvailableContainers(() => {
        this.loading = false
      })
      this.$root.$on("refreshData", () => {
        this.loading = true
        this.getAvailableContainers(() => {
          this.loading = false
        })
      })
    },
    methods: {
      stateColor(state) {
        if (state === "UP") {
          return "success"
        } else if (state === "DOWN") {
          return "error"
        } else {
          return "warning"
        }
      }
    }
  }
</script>

<style scoped>

</style>