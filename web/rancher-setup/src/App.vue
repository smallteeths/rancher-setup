<script setup lang="ts">
import { reactive, onMounted, inject } from 'vue';
import { ElMessage } from 'element-plus'
import InstallRancherView from './views/InstallRancherView.vue'

const http = inject<MyHttp | null>('http');

if (!http) {
  throw new Error('HTTP client not provided');
}

interface State {
  data: any[];
  loading: boolean;
  rancherInstalled: boolean; 
  rancherIsActive: boolean;
  error: Error | null;
}

const state: State = reactive({
  data: [],
  loading: false,
  error: null,
  rancherInstalled: false,
  rancherIsActive: false,
});

const reloadRancherState = () => {
  fetchData()
};

const fetchData = async () => {
  try {
    const url: string = '/rancherState';
    const response = await http.get(url);
    state.rancherInstalled = response?.data?.data?.rancherPodExists;
    state.rancherIsActive = response?.data?.data?.rancherPodActive;
  } catch (error) {
    state.error = error;
    ElMessage.error(`ACK 不可用: ${state.error?.message}`)
  }
};

const fetchPeriodically = () => {
  fetchData();
  setInterval(fetchData, 2 * 60 * 1000);
};

onMounted(() => {
  fetchPeriodically();
});

</script>

<template>
  <div style="position: absolute; width: 50%; height: 100%; overflow:hidden; z-index: -1">
    <img style="width: 900px; height: 100%;" src="@/assets/background.svg" width="125" height="125" />
  </div>
  <header>
    <img alt="Vue logo" class="logo" src="@/assets/logo.svg" width="125" height="125" />
    <el-card style="width: 500px; height: 600px">
      <p class="text item"><strong>Host</strong>: 你的 Rancher Server 域名，注意阿里云需要已备案的域名</p>
      <p class="text item"><strong>Version</strong>: 请参考 ACK 版本选择 Rancher 版本</p>
      <p class="text item"><strong>Harbor User</strong>: Harbor 账户用户名</p>
      <p class="text item"><strong>Harbor Passwd</strong>: Harbor 账户密码</p>
      
      <div v-if="!state.rancherIsActive">
        <el-result v-if="!state.rancherInstalled" style="margin-top: 80px;" icon="warning" title="Rancher" subTitle="Rancher 尚未部署" />
        <el-result v-else style="margin-top: 80px;" icon="success" title="Rancher" subTitle="Rancher 已经部署请等待 Rancher 部署完成" />
      </div>
      <div v-else>
        <el-result style="margin-top: 80px;" icon="success" title="Rancher" subTitle="Rancher Rancher 部署完成请访问 Host" />
      </div>
    </el-card>
  </header>
  <InstallRancherView
    v-loading="state.rancherInstalled && !state.rancherIsActive"
    :rancherInstalled="state.rancherInstalled"
    :rancherIsActive="state.rancherIsActive"
    @reloadRancherState="reloadRancherState"
  />
</template>

<style scoped>
header {
  line-height: 1.5;
  max-height: 100vh;
}

.logo {
  display: block;
  margin: 0 auto 2rem;
}

nav {
  width: 100%;
  font-size: 12px;
  text-align: center;
  margin-top: 2rem;
}

nav a.router-link-exact-active {
  color: var(--color-text);
}

nav a.router-link-exact-active:hover {
  background-color: transparent;
}

nav a {
  display: inline-block;
  padding: 0 1rem;
  border-left: 1px solid var(--color-border);
}

nav a:first-of-type {
  border: 0;
}


header {
  display: flex;
  flex-direction: column;
  flex-wrap: wrap;
  align-items: center;
  justify-content: center;
  padding-right: calc(var(--section-gap) / 2);
}

.logo {
  margin: 0 2rem 0 0;
}

header .wrapper {
  display: flex;
  place-items: flex-start;
}

header .text {
  line-height: 40px;
}

header strong {
  font-weight: 700;
}

nav {
  text-align: left;
  margin-left: -1rem;
  font-size: 1rem;

  padding: 1rem 0;
  margin-top: 1rem;
}

</style>
