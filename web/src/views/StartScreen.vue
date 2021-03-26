<template>
  <div style="width: 100%">
    <div class="main-background">
      <el-row>
        <el-col>
          <BigLogo />
        </el-col>
      </el-row>
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
            @click="joinGame"
          >
            Join Game
          </el-button>
        </el-col>
      </el-row>
    </div>
  </div>
</template>
<script>
import { createGame } from '@/lib/whatsub'
import BigLogo from '@/views/BigLogo'

export default {
  name: 'StartScreen',
  components: { BigLogo },
  mounted () {
    this.$store.commit('clearGameData')
  },
  methods: {
    openNewGame () {
      createGame()
        .then(game => this.$router.push('/game/' + game.key + '/create'))
    },
    joinGame () {
      this.$router.push('/game/join')
    }
  }
}
</script>
<style scoped>
.main-background{
  width: 350px;
  margin: auto;
  border-radius: 25px;
  background-color: rgba(0, 0, 0, 0.2);
  padding-top: 30px;
  padding-bottom: 30px;
}
</style>
