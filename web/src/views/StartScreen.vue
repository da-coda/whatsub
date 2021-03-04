<template>
  <el-container>
    <el-header height="260px">
      <BigLogo />
      <p>Battle against your friends and find out who knows Reddit the best!</p>
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
import { createGame, joinGame } from '@/lib/whatsub'
import BigLogo from '@/views/BigLogo'

export default {
  name: 'StartScreen',
  components: { BigLogo },
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
    openNewGame () {
      const username = this.$prompt('Please enter your username', 'Username', {
        confirmButtonText: 'OK',
        cancelButtonText: 'Cancel'
      })
      const game = createGame()
      Promise.all([username, game])
        .then(([username, game]) => {
          return new Promise((resolve) => {
            const connection = joinGame(username.value, game.key)
            resolve([game, connection])
          })
        })
        .then(([game, connection]) => {
          this.$store.commit('setGameId', game.uuid)
          this.$store.commit('setGameHead', true)
          this.$store.commit('setWebsocketConnection', connection)
          this.$router.push('/game/join/' + game.key)
        })
    }
  }
}
</script>
