//配置文件
export default{
    axios:{
        baseURL:process.env.VUE_APP_BASEURL,
        timeout:5000,
        SUCCESS_CODE:10000,
        NEED_LOGIN_CODE:10004,
        INVALID_TOKEN:10005,
    },
    foot_text:"Copyright 2021@XDU Rick Shang",
    loadmore_text:"加载更多",
    nomore_text:"没有更多了",
    ReCAPTCHA:{
        sitekey:process.env.VUE_APP_RECAPTCHA_SITEKEY
    },
    ws:{
        protocol:process.env.VUE_APP_WS_PROTOCOL,
        domain:process.env.VUE_APP_WS_DOMAIN,
        port:process.env.VUE_APP_WS_PORT,
        path:"api/v1/ws"
    }
}
