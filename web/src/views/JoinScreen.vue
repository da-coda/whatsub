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
          <span style="display: block; margin: auto; width: 200px; text-align: left;">Username</span>
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
          <span style="display: block; margin: auto; width: 200px; text-align: left; padding-top: 10px">Game Id</span>
          <el-input
            v-model="gameKey"
            maxlength="40"
            style="width: 200px;"
            placeholder="Enter Game Id"
            :disabled="disableGameId"
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
            {{ joinButtonLabel }}
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
  name: 'JoinScreen',
  components: { BigLogo },
  props: {
    code: {
      type: String,
      default: null
    },
    isGameHead: {
      type: Boolean,
      default: false
    }
  },
  data () {
    return {
      username: null,
      gameKey: this.code
    }
  },
  computed: {
    disableGameId () {
      return this.code !== null
    },
    joinButtonLabel () {
      if (this.isGameHead) {
        return 'Create Game'
      }
      return 'Join Game'
    }
  },
  methods: {
    openLobby () {
      const playerUUID = uuidv4()
      const connection = joinGame(this.username, this.gameKey, playerUUID)
      connection.onopen = () => {
        console.log('Connection to websocket successful')
        this.$store.commit('setGameData',
          {
            gameShortId: this.gameKey,
            playerUUID: playerUUID,
            playerName: this.username,
            isGameHead: this.isGameHead
          }
        )
        this.$store.commit('setWebsocketConnection', connection)
        this.$router.push('/game/' + this.gameKey + '/lobby')
      }
      connection.onerror = () => {
        console.log('Could not establish a web socket connection!')
        this.$router.push('/')// TODO Maybe try again?
      }
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
