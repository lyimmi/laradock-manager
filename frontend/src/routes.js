import Dashboard from './components/dashboard';
import Settings from './components/settings/index';
// import Containers from './components/containers/containersList'

const routes = [
    {path: '/', component: Dashboard, name: "home"},
    {path: '/settings', component: Settings, name: "settings"},
    // {path: '/containers', component: Containers, name: "containers"},
    {path: '*', component: Dashboard}
];

export default routes;