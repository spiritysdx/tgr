<template>
  <div id="userLayout">
    <div class="login-panel">
      <div class="login-panel-form">
        <div class="login-panel-form-title">
          <img class="login-panel-form-title-logo" src="~@/assets/logo.png" alt>
          <p class="login-panel-form-title-p">{{ $GIN_VUE_ADMIN.appName }}</p>
        </div>
        <el-form ref="loginForm" :model="loginFormData" :rules="rules" @keyup.enter="submitForm">
          <el-form-item prop="username">
            <el-input v-model="loginFormData.username" placeholder="请输入用户名">
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
            <el-input v-model="loginFormData.password" :type="lock === 'lock' ? 'password' : 'text'"
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
              <el-input v-model="loginFormData.captcha" placeholder="请输入验证码" class="captcha-input" />
              <div class="vPic">
                <img v-if="picPath" :src="picPath" alt="请输入验证码" @click="loginVerify()">
              </div>
            </div>
          </el-form-item>
          <!-- Newly added input for TG code -->
          <el-form-item prop="code">
            <el-input v-model="loginFormData.code" placeholder="请输入TG验证码"></el-input>
          </el-form-item>
          <el-form-item>
            <!-- Button to send TG code -->
            <el-button @click="sendTGCode">发送TG验证码</el-button>
            <el-button type="primary" size="large" @click="submitForm">
              <div v-if="loginType.status">注册</div>
              <div v-else>登录</div>
            </el-button>
            <el-switch v-model="loginType.status" />
          </el-form-item>
        </el-form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { captcha } from '@/api/user'
import { getCode } from '@/plugin/register/api/api'
import { reactive, ref } from 'vue'
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
  } else if (!/^(?=.*[a-zA-Z])(?=.*[^\x00-\xff]).{6,}$/.test(value)) {
    return callback(new Error('密码必须包含中英文混合'))
  } else {
    callback()
  }
}
const loginVerify = () => {
  captcha({}).then((ele) => {
    rules.captcha[1].max = ele.data.captchaLength
    rules.captcha[1].min = ele.data.captchaLength
    picPath.value = ele.data.picPath
    loginFormData.captchaId = ele.data.captchaId
  })
}
loginVerify()

const lock = ref('lock')
const changeLock = () => {
  lock.value = lock.value === 'lock' ? 'unlock' : 'lock'
}
const loginForm = ref(null)
const picPath = ref('')
const loginFormData = reactive({
  username: 'admin',
  password: '123456',
  captcha: '',
  captchaId: '',
  code: '',
})
const loginType = reactive({
  status: false,
})
const rules = reactive({
  username: [{ validator: checkUsername, trigger: 'blur' }],
  password: [{ validator: checkPassword, trigger: 'blur' }],
  captcha: [
    { required: true, message: '请输入验证码', trigger: 'blur' },
    { message: '验证码格式不正确', trigger: 'blur' },
  ],
  code: [{ required: true, message: '请输入TG验证码', trigger: 'blur' }],
})
const userStore = extendedUseUserStore()
const login = async () => {
  return await userStore.LoginIn(loginFormData)
}
const register = async () => {
  return await userStore.Register(loginFormData)
}
const submitForm = () => {
  loginForm.value.validate(async (v) => {
    if (v) {
      let flag
      if (loginType.status) {
        flag = await register()
      } else {
        flag = await login()
      }
      if (!flag) {
        loginVerify()
        loginType.status = false
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
  const tgid = "your_tg_id_here"; // Replace this with actual TG ID
  getCode({ tgid })
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
  width: 100px; /* Adjust width as needed */
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
  width: calc(100% - 90px); /* Adjust width as needed */
}
</style>
