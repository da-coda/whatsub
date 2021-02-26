<template>
  <el-container>
    <h1>WhatSub?</h1>
    <el-main />
    <button
      type="button"
      @click="$router.push('/startGame')"
    >
      Create a new Game
    </button>
  </el-container>
</template>
<script>

export default {
  name: 'NewGame',
  components: {}
}
</script>
