<template>
  <img
    src="../../public/img/logo.png"
    width="75"
    style="float: right;"
  >
  <div
    v-if="isGameHead"
    class="box"
  >
    <p>
      Game ID: <strong id="code-span">{{ code }}</strong> <i
        class="el-icon-document-copy"
        @click="copyCode(code)"
      />
    </p>
    <p>
      Game link: <a
        :href="gameLink"
        class="href"
      >{{ gameLink }}</a> <i
        class="el-icon-document-copy"
        @click="copyCode(gameLink)"
      />
    </p>
  </div>
  <h3>Player</h3>
  <div
    id="player"
    class="box flexy"
  >
    <div
      v-for="player in players"
      :key="player"
    >
      {{ player }}
    </div>
  </div>
  <el-button
    v-if="isGameHead"
    type="success"
    style="width: 200px; margin-top: 10px"
  >
    Start Game
  </el-button>
</template>
<script>

import { joinGame, baseUrl } from '@/lib/whatsub'
import { v4 as uuidv4 } from 'uuid'

export default {
  name: 'GameLobby',
  components: { },
  props: {
    code: {
      type: String,
      default: null
    }
  },
  data () {
    return {
      loading: true,
      rejoin: false,
      players: []
    }
  },
  computed: {
    isGameHead () {
      return this.$store.state.isGameHead
    },
    gameLink () {
      return baseUrl + this.$router.resolve({ name: 'JoinScreenByLink', code: this.code }).href
    },
    playerGrid () {
      return this.players.join(' ')
    }
  },
  mounted () {
    if (this.$store.state.websocketConnection === null) {
      let webSocket = null
      let playerName = this.$store.state.playerName
      let playerUUID = this.$store.state.playerUUID
      const isRejoin = [playerName, playerUUID].every(value => value !== null)
      if (!isRejoin) {
        playerName = 'Unnamed Player'
        playerUUID = uuidv4()
      }
      webSocket = joinGame(playerName, this.code, playerUUID)
      this.$store.commit('setWebsocketConnection', webSocket)
      this.$store.commit('setGameData', {
        gameUUID: this.code,
        playerUUID: playerUUID,
        playerName: playerName,
        isGameHead: this.$store.state.isGameHead
      })
    }

    this.loading = false
    const that = this
    this.$store.state.websocketConnection.onmessage = (event) => {
      const msg = JSON.parse(event.data)
      console.log(event.data)
      switch (msg.Type) {
        case 'join': {
          // TODO: display joined user
          break
        }
        case 'left': {
          // TODO: display joined user
          break
        }
        case 'score': {
          console.log('I got a score')
          const data = JSON.parse(event.data)
          console.log(data.Payload.Scores)
          that.players = Object.keys(data.Payload.Scores)
          break
        }
        default:
          console.log(event.data)
      }
    }
  },
  methods: {
    copyCode (clip) {
      navigator.clipboard.writeText(clip).then(() => {
        this.$message({ message: 'Copied to clipboard' })
      })
    }
  }
}
</script>
<style scoped>
div.box{
  width: 700px;
  margin: auto;
  border-radius: 25px;
  background-color: rgba(0, 0, 0, 0.2);
  padding-top: 10px;
  padding-bottom: 10px;
}
div.flexy{
  display: flex;
  flex-direction: row;
  flex-wrap: wrap;
  justify-content: center;
}
.flexy > div {
  margin: 10px;
}
a.href {
  color: whitesmoke;
  font-weight: 700;
}
i.el-icon-document-copy {
  cursor: pointer;
}
</style>
