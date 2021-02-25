<template>
  <el-container>
    <el-header>
      <h1>WhatSub?</h1>
    </el-header>
    <el-main>
      <el-row>
        <el-col>
          <button
            type="button"
            style="width: 200px"
            @click="$router.push('/newGame')"
          >
            Start new Game
          </button>
        </el-col>
      </el-row>
      <el-row>
        <el-col>
          <button
            type="button"
            style="width: 200px; margin-top: 10px"
            @click="enterGameCode"
          >
            Join Lobby
          </button>
        </el-col>
      </el-row>
    </el-main>
  </el-container>
</template>
<script>

export default {
  name: 'StartScreen',
  components: {},
  methods: {
    enterGameCode () {
      this.$prompt('Please enter the game code:', 'Join Game', {
        confirmButtonText: 'OK',
        cancelButton: 'Cancel'
      }).then(({ value }) => {
        this.$message({
          type: 'success',
          message: 'YAYAYYAY'
        })
      }).catch(() => {
        this.$message({
          type: 'info',
          message: 'Input canceled'
        })
      }
      )
    }
  }
}
</script>
