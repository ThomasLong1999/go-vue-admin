<!--
  App.vue - 根组件
  【知识点】Vue 单文件组件(SFC): template + script + style 写在一个 .vue 文件中
  - <template>: HTML 模板，定义组件的 UI 结构
  - <script setup>: JavaScript 逻辑，<script setup> 是 Vue3.2+ 的语法糖
  - <style>: CSS 样式，scoped 表示样式只作用于当前组件
-->
<template>
  <!-- 【知识点】<router-view> 是 Vue Router 的核心组件
       它根据当前 URL 显示对应的页面组件
       例如: 访问 /login → 显示 Login.vue, 访问 /dashboard → 显示 Dashboard.vue -->
  <router-view />
</template>

<script setup>
// 【知识点】<script setup> 语法糖
// 在这个标签里写的代码会自动成为组件的逻辑
// 不需要 export default，不需要 return，顶层变量自动在模板中可用
// 这是 Vue3 推荐的写法，比 Options API（data/methods/computed）更灵活
</script>

<style>
/* 全局样式 */
body {
  margin: 0;
  padding: 0;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}
</style>
