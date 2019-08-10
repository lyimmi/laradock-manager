import Vue from 'vue'
import Vuex from 'vuex'
import VuexPersist from '../plugins/vuex-persist'
import modules from './modules'

const persist = new VuexPersist()

Vue.use(Vuex)

export default new Vuex.Store({
  modules,
  plugins: [persist.subscribe()],
  strict: process.env.NODE_ENV !== 'production'
})
