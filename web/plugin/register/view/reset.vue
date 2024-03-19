<template>
  <div id="resetPasswordLayout">
    <div class="reset-password-panel">
      <div class="reset-password-panel-form">
        <div class="reset-password-panel-form-title">
          <img class="reset-password-panel-form-title-logo" src="~@/assets/logo.png" alt>
          <p class="reset-password-panel-form-title-p">{{ $GIN_VUE_ADMIN.appName }}密码重置</p>
        </div>
        <el-form ref="resetPasswordForm" :model="resetFormData" :rules="resetFormRules" @keyup.enter="submitForm">
          <el-form-item prop="tg_id">
            <el-input v-model="resetFormData.tg_id" placeholder="请输入TGID"></el-input>
          </el-form-item>
          <el-form-item prop="captcha">
            <div class="captcha-input-container">
              <el-input v-model="resetFormData.captcha" placeholder="请输入TG验证码" class="captcha-input" style="width: calc(100% - 150px);" />
              <div class="send-tg-code-container">
                <el-button @click="sendTGCode">发送TG验证码</el-button>
              </div>
            </div>
          </el-form-item>
          <el-form-item prop="oldPassword">
            <el-input v-model="resetFormData.oldPassword" :type="lock === 'lock' ? 'password' : 'text'" placeholder="请输入旧密码">
              <template #suffix>
                <span class="input-icon">
                  <el-icon>
                    <component :is="lock" @click="changeLock" />
                  </el-icon>
                </span>
              </template>
            </el-input>
          </el-form-item>
          <el-form-item prop="newPassword">
            <el-input v-model="resetFormData.newPassword" :type="lock === 'lock' ? 'password' : 'text'" placeholder="请输入新密码">
              <template #suffix>
                <span class="input-icon">
                  <el-icon>
                    <component :is="lock" @click="changeLock" />
                  </el-icon>
                </span>
              </template>
            </el-input>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" size="large" @click="submitForm">确认重置</el-button>
          </el-form-item>
        </el-form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { getCode, UChangePassword } from '@/plugin/register/api/api'
import { reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'

const resetFormData = reactive({
  tg_id: '',
  oldPassword: '',
  newPassword: '',
  captcha: ''
})

const lock = ref('lock')

const resetFormRules = reactive({
  tg_id: [{ required: true, message: '请输入TG ID', trigger: 'blur' }],
  oldPassword: [{ required: true, message: '请输入旧密码', trigger: 'blur' }],
  newPassword: [{ required: true, message: '请输入新密码', trigger: 'blur' }],
  captcha: [{ required: true, message: '请输入TG验证码', trigger: 'blur' }]
})

const userStore = useUserStore()

const changeLock = () => {
  lock.value = lock.value === 'lock' ? 'unlock' : 'lock'
}

const submitForm = () => {
  const validationResult = $refs.resetPasswordForm.validate()
  if (validationResult) {
    UChangePassword(resetFormData)
      .then(() => {
        ElMessage({
          type: 'success',
          message: '密码重置成功',
          showClose: true
        })
      })
      .catch(error => {
        console.error('Error resetting password:', error)
        ElMessage({
          type: 'error',
          message: '密码重置失败，请稍后再试',
          showClose: true
        })
      })
  }
}

const sendTGCode = () => {
  const tg_id = resetFormData.tg_id
  getCode({ tg_id })
    .then(() => {
      ElMessage({
        type: 'success',
        message: 'TG验证码已发送，请查收',
        showClose: true
      })
    })
    .catch(error => {
      console.error('Error sending TG code:', error)
      ElMessage({
        type: 'error',
        message: '发送TG验证码时出错，请稍后再试',
        showClose: true
      })
    })
}
</script>

<style scoped>
.reset-password-panel {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
}
.reset-password-panel-form {
  max-width: 400px;
  width: 100%;
}
.reset-password-panel-form-title {
  text-align: center;
}
.reset-password-panel-form-title-logo {
  width: 100px;
}
.reset-password-panel-form-title-p {
  margin-top: 10px;
  font-size: 20px;
  font-weight: bold;
}
.input-icon {
  position: absolute;
  right: 10px;
  top: 50%;
  transform: translateY(-50%);
}
.captcha-input-container {
  position: relative;
  display: flex;
  align-items: center;
}
.captcha-input {
  width: calc(100% - 150px);
}
.send-tg-code-container {
  display: inline-block;
  width: 120px;
  margin-top: 0;
  text-align: right;
}
</style>
