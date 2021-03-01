<template>
  <el-container>
    <el-header style="height: 120px">
      <h1>
        WhatSub?
      </h1>
      <h3>Let the battle begin!</h3>
    </el-header>
    <el-main>
      <p>Copy this link to invite others: {{ gameLink }} or ask them to join a game with this code {{ code }}</p>
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
      return this.$store.state.isGameHead
    },
    gameLink () {
      return window.location.protocol + '//' + window.location.host + this.$route.path
    }
  },
  methods: {}
}
</script>
