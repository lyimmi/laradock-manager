import Dashboard from './components/dashboard';
import Settings from './components/settings/index';
import Stats from './components/stats'

const routes = [
    { path: '/', component: Dashboard, name: "home" },
    { path: '/stats', component: Stats, name: "stats" },
    { path: '/settings', component: Settings, name: "settings" },
    { path: '*', component: Dashboard }
];

export default routes;