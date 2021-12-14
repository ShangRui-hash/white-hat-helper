import axios from 'axios'
import {
  MessageBox,
  Message
} from 'element-ui'
import store from '@/store'
import {
  getToken
} from '@/utils/auth'
import config from '@/config'

// create an axios instance
const service = axios.create({
  baseURL: process.env.VUE_APP_BASE_API, // url = base url + request url
  // withCredentials: true, // send cookies when cross-domain requests
  timeout: 5000 // request timeout
})

// request interceptor
service.interceptors.request.use(
  request => {
    // do something before request is sent
    if (store.getters.token) {
      let token = getToken()
      console.log("token",token)
      if (token) {
        request.headers.Authorization = `Bearer ${token}`;
      }
    }
    return request
  },
  error => {
    return Promise.reject(error)
  }
)

// 拦截接收到的响应
service.interceptors.response.use(
  resp => {
    //如果需要登录
    let code = resp.data.code
    if (code == config.axios.NEED_LOGIN_CODE || code == config.axios.INVALID_TOKEN) {
      router.push("/login")
      return Promise.reject("需要登录")
    }
    //返回的响应码不是成功
    if (code != config.axios.SUCCESS_CODE) {
      let msg
      if (typeof resp.data.msg == "string") {
        msg = resp.data.msg
      } else {
        msg = Object.values(resp.data.msg)[0]
      }
      Message({
        message: msg,
        type: "warning",
      })
      return Promise.reject(resp.data.msg)
    }
    return resp
    //   const res = response.data

    //   // if the custom code is not 20000, it is judged as an error.
    //   if (res.code !== 20000) {
    //     Message({
    //       message: res.message || 'Error',
    //       type: 'error',
    //       duration: 5 * 1000
    //     })

    //     // 50008: Illegal token; 50012: Other clients logged in; 50014: Token expired;
    //     if (res.code === 50008 || res.code === 50012 || res.code === 50014) {
    //       // to re-login
    //       MessageBox.confirm('You have been logged out, you can cancel to stay on this page, or log in again', 'Confirm logout', {
    //         confirmButtonText: 'Re-Login',
    //         cancelButtonText: 'Cancel',
    //         type: 'warning'
    //       }).then(() => {
    //         store.dispatch('user/resetToken').then(() => {
    //           location.reload()
    //         })
    //       })
    //     }
    //     return Promise.reject(new Error(res.message || 'Error'))
    //   } else {
    //     return res
    //   }
  },
  error => {
    console.log('err' + error) // for debug
    Message({
      message: error.message,
      type: 'error',
      duration: 5 * 1000
    })
    return Promise.reject(error)
  }
)

export default service
