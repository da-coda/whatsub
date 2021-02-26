<template>
  <el-container>
    <el-header style="height: 120px">
      <h1>
        WhatSub?
      </h1>
      <h3>Let the battle begin!</h3>
    </el-header>
    <el-main>
      <p>{{ this.$store.state.gameId }}</p>
      <button
        v-if="isGameHead"
        type="button"
        @click="$router.push('/startGame')"
      >
        Start the game
      </button>
    </el-main>
  </el-container>
</template>
<script>

export default {
  name: 'NewGame',
  props: {
    code: {
      type: String,
      default: null
    }
  },
  computed: {
    isGameHead () {
      return this.$store.getters.isGameHead
    }
  },
  methods: {}
}
</script>
