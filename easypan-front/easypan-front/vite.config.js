import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

const fs = require('fs')
const path = require('path')

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  define: {
    'process.env': {}
  },
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
  server: {
    disableHostCheck: true,//开启反向代理
    port: 1024,
    hmr: true,
    //secure: false,
    open:true,
    cors:true,//允许跨域
    https: {
      // 主要是下面两行的配置文件，不要忘记引入 fs 和 path 两个对象
      cert: fs.readFileSync(path.join(__dirname, 'src/ssl/cert.crt')),
      key: fs.readFileSync(path.join(__dirname, 'src/ssl/cert.key'))
    },
    //disableHostCheck: true,
    proxy: {
      "/api": {
        target: "https://localhost:443",
        changeOrigin: true,
        ws:true,
        secure:false,
        // onProxyReq:function (proxyReq, req, res, options) {
        //   if (req.body) {
        //     let bodyData = JSON.stringify(req.body);
        //     // incase if content-type is application/x-www-form-urlencoded -> we need to change to application/json
        //     proxyReq.setHeader('Content-Type','application/json');
        //     proxyReq.setHeader('Content-Length', Buffer.byteLength(bodyData));
        //     // stream the content
        //     proxyReq.write(bodyData);
        //   }},
        pathRewrite: {
          "^/api": "/",
        }
      }
    },
    //https: true

  },
  build: {
    chunkSizeWarningLimit: 3000,
    rollupOptions: {
      output: {
        manualChunks(id) {
          if (id.includes('node_modules')) {
            return id.toString().split('node_modules/')[1].split('/')[0].toString();
          }
        }
      }
    },
    chunkFileNames: (chunkInfo) => {
      const facadeModuleId = chunkInfo.facadeModuleId
        ? chunkInfo.facadeModuleId.split('/')
        : [];
      const fileName =
        facadeModuleId[facadeModuleId.length - 2] || '[name]';
      return `js/${fileName}/[name].[hash].js`;
    }
  }
})
