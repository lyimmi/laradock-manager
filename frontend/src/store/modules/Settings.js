const namespaced = true

const state = {
    laradockPath: ""
}

const mutations = {
    SET_LARADOCK_PATH(state, payload) {
        state.laradockPath = payload
    },
}

const getters = {
    laradockPath(state) {
        return state.laradockPath
    },
}

const actions = {
    setLaradockPath(context, payload) {
        context.commit('SET_LARADOCK_PATH', payload)
    }
}

export default {
    namespaced,
    state,
    mutations,
    getters,
    actions
}
