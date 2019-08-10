<template>
    <v-app id="inspire">
        <app-menu></app-menu>

        <v-content>
            <transition>
                <router-view></router-view>
            </transition>
        </v-content>

        <v-footer app>
            <span><a href="https://github.com/lyimmi"><v-icon>fab fa-github</v-icon> lyimmi</a></span>
        </v-footer>
    </v-app>
</template>

<script>
    import AppMenu from "./components/partials/app-menu";

    export default {
        components: {AppMenu},
        mounted() {
            this.$router.push("home");
        },
        props: {
            source: String,
        },
        data: () => ({}),
        created() {
            this.$vuetify.theme.dark = true
        },
    }
</script>

<style>
    .logo {
        width: 16em;
    }

    a {
        text-decoration: none;
    }
</style>