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
          <span style="display: block; margin: auto; width: 200px; text-align: left;">Game Id</span>
          <el-input
            v-model="gameKey"
            maxlength="40"
            style="width: 200px;"
            placeholder="Enter Game Id"
            :disabled="true"
          />
        </el-col>
      </el-row>
      <el-row>
        <el-col>
          <span style="display: block; margin: auto; width: 200px; text-align: left; padding-top: 10px">Username</span>
          <el-input
            v-model="username"
            maxlength="40"
            style="width: 200px;"
            placeholder="Enter username"
          />
        </el-col>
      </el-row>
      <el-row>
        <el-col>
          <el-button
            type="success"
            style="width: 200px; margin-top: 10px"
            @click="openLobby"
          >
            Create Game
          </el-button>
        </el-col>
      </el-row>
    </div>
  </div>
</template>
<script>
import { joinGame } from '@/lib/whatsub'
import BigLogo from '@/views/BigLogo'
import { v4 as uuidv4 } from 'uuid'

export default {
  name: 'StartScreen',
  components: { BigLogo },
  props: {
    code: {
      type: String,
      default: null
    }
  },
  data () {
    return {
      username: null,
      gameKey: this.code
    }
  },
  methods: {
    openLobby () {
      const playerUUID = uuidv4()
      const connection = joinGame(this.username, this.code, playerUUID)
      this.$store.commit('setGameData',
        {
          gameUUID: this.code,
          playerUUID: playerUUID,
          playerName: this.username,
          isGameHead: true
        }
      )
      this.$store.commit('setWebsocketConnection', connection)
      this.$router.push('/game/' + this.code + '/lobby')
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
