import Vue from "vue"
const namespaced = true

const state = {
    env: false,
    docker: false,
    dockerCompose: false,
    laradockPath: false
}

const mutations = {
    SET_STATUS(state, payload) {
        Vue.set(state, payload.status, payload.value);
    },
}

const getters = {
    appStatus(state) {
        return state
    },
}

const actions = {
    setAppStatus(context, payload) {
        context.commit('SET_STATUS', payload)
    }
}

export default {
    namespaced,
    state,
    mutations,
    getters,
    actions
}
