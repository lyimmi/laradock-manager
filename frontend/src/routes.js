import Dashboard from './components/dashboards/index.vue';
import Settings from './components/settings/index.vue';

const routes = [
    {path: '/', component: Dashboard, name: "home"},
    {path: '/settings', component: Settings, name: "settings"},
    {path: '*', component: Dashboard}
];

export default routes;