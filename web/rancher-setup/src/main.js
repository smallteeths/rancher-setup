import './assets/main.css'

import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import { MyHttp } from './axios';

const app = createApp(App);

app.use(router);
const baseURL = window.location.href;
const http = new MyHttp(baseURL);

app.provide('http', http);
app.use(ElementPlus);

app.mount('#app');
