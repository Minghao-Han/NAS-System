<template>
  <div class="login-body">
<!--    左侧图标-->
    <div class="bg"></div>
<!--    右侧窗口-->
    <div class="login-panel">
      <el-form
        class="login-register"
        :model="formData"
        :rules="rules"
        ref="formDataRef"
      >
        <div class="login-title">Easy云盘</div>
        <!--input输入-->

        <el-form-item prop="account">
          <el-input
              size="large"
              clearable
              placeholder="请输入账号"
              v-model.trim="formData.username"
              maxLength="150"
          >
            <template #prefix>
              <span class="iconfont icon-account"></span>
            </template>
          </el-input>
        </el-form-item>
        <!--登录密码-->
        <el-form-item prop="password" v-if="opType == 1">
          <el-input
            type="password"
            size="large"
            placeholder="请输入密码"
            v-model.trim="formData.password"
            show-password
          >
            <template #prefix>
              <span class="iconfont icon-password"></span>
            </template>
          </el-input>
        </el-form-item>
        <!--注册-->
        <div v-if="opType == 0 || opType == 2">
          <el-form-item prop="nickName" v-if="opType == 0">
            <el-input
              size="large"
              clearable
              placeholder="请输入昵称"
              v-model.trim="formData.nickName"
              maxLength="20"
            >
              <template #prefix>
                <span class="iconfont icon-account"></span>
              </template>
            </el-input>
          </el-form-item>
          <el-form-item prop="registerPassword">
            <el-input
              type="password"
              size="large"
              placeholder="请输入密码"
              v-model.trim="formData.registerPassword"
              show-password
            >
              <template #prefix>
                <span class="iconfont icon-password"></span>
              </template>
            </el-input>
          </el-form-item>
          <el-form-item prop="reRegisterPassword">
            <el-input
              type="password"
              size="large"
              placeholder="请再次输入密码"
              v-model.trim="formData.reRegisterPassword"
              show-password
            >
              <template #prefix>
                <span class="iconfont icon-password"></span>
              </template>
            </el-input>
          </el-form-item>
        </div>

        <el-form-item v-if="opType == 1">
          <div class="rememberme-panel">
            <el-checkbox v-model="formData.rememberMe">记住我</el-checkbox>
          </div>
          <div class="no-account">
            <a href="javascript:void(0)" class="a-link" @click="showPanel(2)"
              >忘记密码？</a
            >
            <a href="javascript:void(0)" class="a-link" @click="showPanel(0)"
              >没有账号？</a
            >
          </div>
        </el-form-item>
        <el-form-item v-if="opType == 0">
          <a href="javascript:void(0)" class="a-link" @click="showPanel(1)"
            >已有账号?</a
          >
        </el-form-item>
        <el-form-item v-if="opType == 2">
          <a href="javascript:void(0)" class="a-link" @click="showPanel(1)"
            >去登录?</a
          >
        </el-form-item>
        <el-form-item>
          <el-button
            type="primary"
            class="op-btn"
            @click="doSubmit"
            size="large"
          >
            <span v-if="opType == 0">注册</span>
            <span v-if="opType == 1">登录</span>
            <span v-if="opType == 2">重置密码</span>
          </el-button>
          <el-button
            type="primary"
            class="op-btn"
            @click="downloadFile"
            size="large"
          > <span>测试下载</span>
          </el-button>
        </el-form-item>

      </el-form>
    </div>
    <!--发送邮箱验证码-->

  </div>
</template>

<script setup>
import { inject } from 'vue'
const $axios = inject('$axios')
import { ref, reactive, getCurrentInstance, nextTick, onMounted } from "vue";
import { useRouter, useRoute } from "vue-router";
import md5 from "js-md5";
import axios from "axios";
import download from "../Func/download";
import exportData  from "../Func/download";
import saveUserInfo from "../Func/saveUserInfo";
//添加文件

const downloadFile = async () => {
  console.log("执行")
  await exportData()
}

const { proxy } = getCurrentInstance();
const router = useRouter();
const route = useRoute();
const api = {

  register: "/register",
  login: "/login",
  resetPwd: "/resetPwd",
};

// 0:注册 1:登录 2:重置密码
const opType = ref(1);
const showPanel = (type) => {
  opType.value = type;
  resetForm();
};

onMounted(() => {
  showPanel(1);
});

const checkRePassword = (rule, value, callback) => {
  if (value !== formData.value.registerPassword) {
    callback(new Error(rule.message));
  } else {
    callback();
  }
};
const formData = ref({});
const formDataRef = ref();
//
const rules = {
  email: [
    { required: true, message: "请输入邮箱" },
    { validator: proxy.Verify.email, message: "请输入正确的邮箱" },
  ],
  password: [{ required: true, message: "请输入密码" }],
  emailCode: [{ required: true, message: "请输入邮箱验证码" }],
  nickName: [{ required: true, message: "请输入昵称" }],
  registerPassword: [
    { required: true, message: "请输入密码" },
    {
      validator: proxy.Verify.password,
      message: "密码只能是数字，字母，特殊字符 8-18位",
    },
  ],
  reRegisterPassword: [
    { required: true, message: "请再次输入密码" },
    {
      validator: checkRePassword,
      message: "两次输入的密码不一致",
    },
  ],
  checkCode: [{ required: true, message: "请输入图片验证码" }],
};
//验证码

//发送邮箱验证码弹窗

//获取邮箱验证码

//发送邮件

//重置表单
const resetForm = () => {
  nextTick(() => {
    //changeCheckCode(0);
    formDataRef.value.resetFields();
    formData.value = {};

    //登录
    if (opType.value == 1) {
      const cookieLoginInfo = proxy.VueCookies.get("loginInfo");
      if (cookieLoginInfo) {
        formData.value = cookieLoginInfo;
      }
    }
  });
};

// 登录、注册、重置密码  提交表单
const doSubmit = () => {
  formDataRef.value.validate(async (valid) => {
    if (!valid) {
      return;
    }
    let params = {};
    Object.assign(params, formData.value);
    //注册
    if (opType.value == 0 || opType.value == 2) {
      params.password = params.registerPassword;
      delete params.registerPassword;
      delete params.reRegisterPassword;
    }
    //登录 密码加密了
    // if (opType.value == 1) {
    //   let cookieLoginInfo = proxy.VueCookies.get("loginInfo");
    //   let cookiePassword =
    //     cookieLoginInfo == null ? null : cookieLoginInfo.password;
    //   if (params.password !== cookiePassword) {
    //     params.password = md5(params.password);
    //   }
    // }
    // let url = null;
    // if (opType.value == 0) {
    //   url = api.register;
    // } else if (opType.value == 1) {
    //   url = api.login;
    // } else if (opType.value == 2) {
    //   url = api.resetPwd;
    // }


    let data = {
  "username":params.username,
  "password":params.password
}

    console.log(data)
    //跳转路由
    proxy.Message.success("登录成功");
    router.push('/main/all').catch(err => {
      console.log(err)})


//     axios.post('https://localhost:443/login', data, {
//       headers: {
//         'Content-Type': 'application/json'
//       }
//     })
//         .then(response => {
//           // 处理响应
//           //console.log(response.data);
//           console.log(response);  proxy.Message.success("登录成功");
// //登录成功
//           if(response.status===200)
//           {
//             localStorage.setItem("token",response.data.token)
//
//             //获取用户信息
//             axios.get('https://localhost:443/user/info',  {
//               headers: {
//                 'Content-Type': 'application/json',
//                 'token':localStorage.getItem("token")
//               }
//             })
//                 .then(response => {
//                   if(response.status===200)
//                   {
//                     console.log(response.data)
//                     saveUserInfo(response.data.userInfo)
//                   }
//                   else proxy.Message.error("获取用户数据失败")
//                   //跳转路由

//                    router.push('/main/all').catch(err => {
//                      console.log(err)
//                    })
//                 })
//                 .catch(error => {
//                   // 处理错误
//                   console.error(error);
//                 });
//           }
//
//           //保存token
//
//         })
//         .catch(error => {
//           // 处理错误
//           console.error(error);
//         });

    // let result = await proxy.Request({
    //   url: url,
    //   data:JSON.stringify(data),
    //
    //   headers: {
    //     'Content-Type': 'application/json'
    //   }
    //
    // });






    //await router.replace({path: "/"})
    //注册返回
    // if (opType.value == 0) {
    //   proxy.Message.success("注册成功,请登录");
    //   showPanel(1);
    // } else if (opType.value == 1) {
    //   //登录
    //
    //   if (params.rememberMe) {
    //     const loginInfo = {
    //       email: params.email,
    //       password: params.password,
    //       rememberMe: params.rememberMe,
    //     };
    //     proxy.VueCookies.set("loginInfo", loginInfo, "7d");
    //   } else {
    //     proxy.VueCookies.remove("loginInfo");
    //   }
    //   proxy.Message.success("登录成功");
    //   //存储cookie
    //   proxy.VueCookies.set("userInfo", result.data, 0);
    //
    //   //重定向到原始页面
    //   const redirectUrl = route.query.redirectUrl || "/";
    //   router.push(redirectUrl);
    // } else if (opType.value == 2) {
    //   //重置密码
    //   proxy.Message.success("重置密码成功,请登录");
    //   showPanel(1);
    // }

  });
};

//QQ登录
</script>

<style lang="scss" scoped>
.login-body {
  height: calc(100vh);
  background-size: cover;
  background: url("../assets/login_bg.jpg");
  display: flex;
  .bg {
    flex: 1;
    background-size: cover;
    background-position: center;
    background-size: 800px;
    background-repeat: no-repeat;
    background-image: url("../assets/login_img.png");

  }
  .login-panel {
    width: 430px;
    margin-right: 15%;
    margin-top: calc((100vh - 500px) / 2);
    .login-register {
      padding: 25px;
      background: #fff;
      border-radius: 5px;
      .login-title {
        text-align: center;
        font-size: 18px;
        font-weight: bold;
        margin-bottom: 20px;
      }
      .send-emali-panel {
        display: flex;
        width: 100%;
        justify-content: space-between;
        .send-mail-btn {
          margin-left: 5px;
        }
      }
      .rememberme-panel {
        width: 100%;
      }
      .no-account {
        width: 100%;
        display: flex;
        justify-content: space-between;
      }
      .op-btn {
        width: 100%;
      }
    }
  }

  .check-code-panel {
    width: 100%;
    display: flex;
    .check-code {
      margin-left: 5px;
      cursor: pointer;
    }
  }
  .login-btn-qq {
    margin-top: 20px;
    text-align: center;
    display: flex;
    align-items: center;
    justify-content: center;
    img {
      cursor: pointer;
      margin-left: 10px;
      width: 20px;
    }
  }
}
</style>