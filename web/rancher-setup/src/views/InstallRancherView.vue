<template>
  <div class="about" v-loading="state.loading">
    <el-form
      ref="ruleFormRef"
      style="min-width: 400px"
      :model="ruleForm"
      status-icon
      :rules="rules"
      label-width="auto"
      label-position="top"
      class="demo-ruleForm"
      :size="size"
    >
      <a v-if="rancherIsActive" style="display: block; position: absolute; display: block; right: 80px; bottom: 506px;" :href="`https://${ruleForm.host}`" target="_blank">
        访问 Rancher
      </a>
      <el-form-item label="Host" prop="pass">
        <el-input v-model="ruleForm.host"/>
      </el-form-item>
      <el-form-item label="Version" prop="checkPass">
        <el-select v-model="ruleForm.version">
          <el-option label="v2.7.9-ent" value="rancher-2.7.9-ent" />
        </el-select>
      </el-form-item>
      <el-form-item label="Harbor User" prop="checkPass">
        <el-input v-model="ruleForm.harborUser"/>
      </el-form-item>
      <el-form-item label="harbor Passwd" prop="checkPass">
        <el-input v-model="ruleForm.harborPasswd" type="password"/>
      </el-form-item>
      <el-form-item>
        <el-button
          type="primary"
          @click="submitForm(ruleFormRef)"
          :disabled="rancherInstalled"
          >Submit</el-button
        >
      </el-form-item>
    </el-form>
  </div>
</template>

<script lang="ts" setup>
import { reactive, ref, onMounted, inject, defineEmits } from 'vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'

const emits = defineEmits(['reloadRancherState']);

const props = defineProps({
  rancherInstalled: { type: Boolean },
  rancherIsActive: { type: Boolean },
})

const http = inject<MyHttp | null>('http');

if (!http) {
  throw new Error('HTTP client not provided');
}

const size = ref('large')
var state = reactive({
  loading: false
})
const ruleFormRef = ref<FormInstance>()
const validateHost = (rule: any, value: any, callback: any) => {
  if (value === '') {
    callback(new Error('Please input the Host'))
  } else {
    callback()
  }
}
const validateVersion = (rule: any, value: any, callback: any) => {
  if (value === '') {
    callback(new Error('Please input the Version'))
  } else {
    callback()
  }
}

const ruleForm = reactive({
  host: '',
  version: 'rancher-2.7.9-ent',
  harborUser: '',
  harborPasswd: '',
})

const rules = reactive<FormRules<typeof ruleForm>>({
  host: [{ validator: validateHost, trigger: 'blur' }],
  version: [{ validator: validateVersion, trigger: 'blur' }],
})

const fetchData = async () => {
  state.loading = true
  try {
    const url: string = '/config';
    const response = await http.get(url);
    ruleForm.host = response?.data?.data?.host
    ruleForm.version = response?.data?.data?.version ? response?.data?.data?.version : 'rancher-2.7.9-ent'
    ruleForm.harborPasswd = response?.data?.data?.harborPasswd
    ruleForm.harborUser = response?.data?.data?.harborUser
  } catch (error) {
    ElMessage.error(error?.response?.data?.message)
  }
  emits('reloadRancherState');
  state.loading = false
};

const applyRancher = async () => {
  state.loading = true
  try {
    const url: string = '/install';
    const response = await http.post(url, {
      "host": ruleForm.host,
      "version": ruleForm.version,
      "harborUser": ruleForm.harborUser,
      "harborPasswd": ruleForm.harborPasswd, 
    });
  } catch (error) {
    ElMessage.error(error?.message)
  }
  fetchData()
};

const submitForm = (formEl: FormInstance | undefined) => {
  if (!formEl) return
  formEl.validate((valid) => {
    if (valid) {
      applyRancher()
    } else {
      console.log('error submit!')
      return false
    }
  })
}

onMounted(() => {
  fetchData();
});
</script>

<style>
.about {
  min-height: 100vh;
  display: flex;
  align-items: center;
  padding-left: 150px;
}
</style>
