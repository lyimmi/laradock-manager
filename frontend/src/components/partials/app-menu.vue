<template>
    <div>
        <v-navigation-drawer v-model="drawer" app clipped>
            <v-list dense>
                <v-list-item to="/">
                    <v-list-item-action>
                        <v-icon>dashboard</v-icon>
                    </v-list-item-action>
                    <v-list-item-content>
                        <v-list-item-title>Dashboard</v-list-item-title>
                    </v-list-item-content>
                </v-list-item>
                <v-list-item to="/settings">
                    <v-list-item-action>
                        <v-icon>settings</v-icon>
                    </v-list-item-action>
                    <v-list-item-content>
                        <v-list-item-title>Settings</v-list-item-title>
                    </v-list-item-content>
                </v-list-item>
            </v-list>
        </v-navigation-drawer>

        <v-app-bar app clipped-left>
            <v-app-bar-nav-icon @click.stop="drawer = !drawer"></v-app-bar-nav-icon>
            <v-toolbar-title>Laradock manager</v-toolbar-title>

            <v-spacer></v-spacer>

            <v-progress-circular
                    :rotate="90"
                    :size="34"
                    :value="refreshProgress"
                    :color="refreshProgress < 70 ? 'success' : 'warning'"
            >

                <v-btn
                        icon
                        @click="$root.$emit('refreshData')" prevent
                >
                    <v-icon>refresh</v-icon>
                </v-btn>
            </v-progress-circular>
        </v-app-bar>
    </div>
</template>

<script>
    export default {
        name: "app-menu",
        data() {
            return {
                drawer: null,
                refreshCounter: 0,
                refreshProgress: 0
            }
        },
        mounted() {

            setInterval(() => {
                this.refreshCounter++;
                this.refreshProgress = Math.floor(this.refreshCounter / 30 * 100)
                if (this.refreshCounter === 30) {
                    this.$root.$emit('refreshData')
                    this.refreshCounter = 0
                }
            }, 2000)
        }
    }
</script>

<style scoped>
    .router-link-exact-active {
    }

    .router-link-active {
    }
</style>