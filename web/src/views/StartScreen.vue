<template>
  <el-container>
    <el-header height="120px">
      <h1>WhatSub?</h1>
      <h3>Battle against your friends and find out who knows Reddit the best!</h3>
    </el-header>
    <el-main>
      <el-row>
        <el-col>
          <el-button
            type="primary"
            style="width: 200px"
            @click="openNewGame"
          >
            Start new Game
          </el-button>
        </el-col>
      </el-row>
      <el-row>
        <el-col>
          <el-button
            type="success"
            style="width: 200px; margin-top: 10px"
            @click="enterGameCode"
          >
            Join Lobby
          </el-button>
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
    },
    async openNewGame () {
      await this.$http.get('/game/create').then((response) => {
        // TODO Error handling
        this.$store.commit('setGameId', response.data.Payload.UUID)
        this.$router.push('/newGame')
      })
    }
  }
}
</script>
