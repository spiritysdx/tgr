import service from '@/utils/request'

// @Summary 用户的TG验证码获取
// @Produce  application/json
// @Param data body {tgid: "string"}
// @Router /register/code [post]
export const getCode = (data) => {
  return service({
    url: '/register/code',
    method: 'post',
    data: data,
  })
}

// @Summary 用户注册
// @Produce  application/json
// @Param data body {username:"string", password:"string", tgid:"string", code:"string", captcha:"string", captchaId:"string"}
// @Router /register/register [post]
export const userRegister = (data) => {
  return service({
    url: '/register/register',
    method: 'post',
    data: data,
  })
}