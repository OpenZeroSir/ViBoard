import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
// import path from 'path';
import { nanoid } from 'nanoid'
import resolve from '@rollup/plugin-node-resolve';
// import videojs from 'video.js';

// import typescript from '@rollup/plugin-typescript';

//import path from 'path';

// // https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    resolve(),
    // copy({
    //   targets: [
    //     { src: 'src/assets/images/charts/*.*', dest: 'dist' },
    //   ],
    //   verbose: true,
    // }),


  ],
  // server: {
  //   proxy: {
  //     "/baidu": {
  //       target: "https://api.map.baidu.com/",
  //       changeOrigin: true,
  //       rewrite: (path) => path.replace(/^\/baidu/, "https://api.map.baidu.com"),
  //     },
  //   },
  // },
  build: {
    assetsInlineLimit: 0, // 将所有文件视为静态资源文件
    // assetsDir: "CeilAI",
    //让文件名不带哈希
    rollupOptions: {
      // external: ['video.js'],
      output: {
        entryFileNames: 'com.zerosir.[hash].js',
        // assetFileNames: '[name].[ext]',
        assetFileNames: (chunkInfo) => {
          if (chunkInfo.name.endsWith('css')) {
            return `css/com.zerosir.${nanoid()}.css`;
          }
          if (chunkInfo.name.endsWith('webp')) {
            return 'webp/[name].[ext]';
          }
          if (chunkInfo.name.endsWith('txt')) {
            return 'txt/[name].[ext]';
          }
          return '[ext]/com.zerosir.[name].[hash].[ext]';
        },

        // plugins: [copy({
        //   targets: [
        //     { src: 'src/assets/images/charts/*', dest: 'dist/assets/' },
        //   ],
        //   verbose: true,
        // })]
      },
      // input: {
      //   app: './src/main.ts'
      // },

    }
  },
  // resolve: {
  //   alias: {
  //     // 将 `@` 指定为 `src` 目录的别名
  //     "@": path.resolve(__dirname, "src/assets/images")
  //   }
  // }


})


