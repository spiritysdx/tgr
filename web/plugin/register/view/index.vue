<template>
  <div id="userLayout">
    <div class="login-panel">
      <div class="login-panel-form">
        <div class="login-panel-form-title">
          <img class="login-panel-form-title-logo" src="~@/assets/logo.png" alt>
          <p class="login-panel-form-title-p">{{ $GIN_VUE_ADMIN.appName }}</p>
        </div>
        <el-form ref="loginForm" :model="currentFormData" :rules="formRules" @keyup.enter="submitForm">
          <el-form-item prop="username">
            <el-input v-model="currentFormData.username" placeholder="请输入用户名">
              <template #suffix>
                <span class="input-icon">
                  <el-icon>
                    <user />
                  </el-icon>
                </span>
              </template>
            </el-input>
          </el-form-item>
          <el-form-item prop="password">
            <el-input v-model="currentFormData.password" :type="lock === 'lock' ? 'password' : 'text'"
                      placeholder="请输入密码">
              <template #suffix>
                <span class="input-icon">
                  <el-icon>
                    <component :is="lock" @click="changeLock" />
                  </el-icon>
                </span>
              </template>
            </el-input>
          </el-form-item>
          <el-form-item prop="captcha">
            <div class="vPicBox">
              <el-input v-model="currentFormData.captcha" placeholder="请输入验证码" class="captcha-input" />
              <div class="vPic">
                <img v-if="picPath" :src="picPath" alt="请输入验证码" @click="loginVerify()">
              </div>
            </div>
          </el-form-item>
          <el-form-item v-if="registerType" prop="tg_id">
            <el-input v-model="currentFormData.tg_id" placeholder="请输入TGID"></el-input>
          </el-form-item>
          <el-form-item v-if="registerType" prop="code">
            <el-input v-model="currentFormData.code" placeholder="请输入TG验证码"></el-input>
          </el-form-item>
          <el-form-item v-if="registerType">
            <el-button @click="sendTGCode">发送TG验证码</el-button>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" size="large" @click="submitForm">
              <div v-if="registerType">注册</div>
              <div v-else>登录</div>
            </el-button>
          </el-form-item>
          <el-form-item>
            <el-switch v-model="registerType" />
          </el-form-item>
        </el-form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { captcha } from '@/api/user'
import { getCode } from '@/plugin/register/api/api'
import { reactive, ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { useRouter } from 'vue-router'
import { extendedUseUserStore } from '@/plugin/register/pinia/modules/user'

const router = useRouter()

// 对用户的输入强制要求符合要求
const checkUsername = (rule, value, callback) => {
  if (value.length < 5) {
    return callback(new Error('用户名必须大于或等于5个字符'))
  } else {
    callback()
  }
}
const checkPassword = (rule, value, callback) => {
  if (value.length < 6) {
    return callback(new Error('密码必须大于或等于6个字符'))
  } else {
    callback()
  }
}
const loginForm = ref(null)
const loginFormData = reactive({
  username: 'admin',
  password: '123456',
  captcha: '',
  captchaId: '',
})
const registerFormData = reactive({
  username: 'admin',
  password: '123456',
  captcha: '',
  captchaId: '',
  tg_id: '',
  code: '',
})
const lock = ref('lock')
const picPath = ref('')
const registerType = ref(false)
const currentFormData = ref(loginFormData)

const rules = reactive({
  username: [{ validator: checkUsername, trigger: 'blur' }],
  password: [{ validator: checkPassword, trigger: 'blur' }],
  captcha: [
    { required: true, message: '请输入验证码', trigger: 'blur' },
    { message: '验证码格式不正确', trigger: 'blur' },
  ],
  tg_id: [{ required: true, message: '请输入TG ID', trigger: 'blur', visible: false }],
  code: [{ required: true, message: '请输入TG验证码', trigger: 'blur', visible: false }],
})

const formRules = ref(rules)

const userStore = extendedUseUserStore()

const loginVerify = () => {
  captcha({}).then((ele) => {
    rules.captcha[1].max = ele.data.captchaLength
    rules.captcha[1].min = ele.data.captchaLength
    picPath.value = ele.data.picPath
    currentFormData.value.captchaId = ele.data.captchaId
  })
}
loginVerify()

const changeLock = () => {
  lock.value = lock.value === 'lock' ? 'unlock' : 'lock'
}

const login = async () => {
  return await userStore.LoginIn(loginFormData)
}
const register = async () => {
  return await userStore.Register(registerFormData)
}
const submitForm = () => {
  const form = loginForm.value
  form.validate(async (v) => {
    if (v) {
      let flag
      if (registerType.value) {
        flag = await register()
      } else {
        flag = await login()
      }
      if (!flag) {
        loginVerify()
        registerType.value = false
      }
    } else {
      ElMessage({
        type: 'error',
        message: '请正确填写登录信息',
        showClose: true,
      })
      loginVerify()
      return false
    }
  })
}
// TG验证码发送
const sendTGCode = () => {
  const tg_id = registerFormData.tg_id
  getCode({ tg_id })
    .then(response => {
      ElMessage({
        type: 'success',
        message: 'TG验证码已发送，请查收',
        showClose: true,
      })
    })
    .catch(error => {
      console.error('Error sending TG code:', error)
      ElMessage({
        type: 'error',
        message: '发送TG验证码时出错，请稍后再试',
        showClose: true,
      })
    })
}

// 监听注册类型的变化，切换表单数据
watch(registerType, (newValue) => {
  currentFormData.value = newValue ? registerFormData : loginFormData
})

</script>

<style scoped>
.login-panel {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
}

.login-panel-form {
  max-width: 400px;
  width: 100%;
}

.login-panel-form-title {
  text-align: center;
}

.login-panel-form-title-logo {
  width: 100px;
}

.login-panel-form-title-p {
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

.captcha-input {
  width: calc(100% - 90px);
}
</style>
