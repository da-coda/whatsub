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
    @click="startGame"
  >
    Start Game
  </el-button>
</template>
<script>

import { startGame, baseUrl, askGameState } from '@/lib/whatsub'

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
      rejoin: false,
      players: []
    }
  },
  computed: {
    isGameHead () {
      return this.$store.state.isGameHead
    },
    gameLink () {
      return baseUrl + '/' + this.$router.resolve({ name: 'JoinScreenByLink', code: this.code }).href
    }
  },
  mounted () {
    console.log('Game started')
    const that = this
    askGameState(this.code).then((answer) => {
      this.update(answer.data.Payload)
    })
    this.$store.state.websocketConnection.onmessage = event => {
      const msg = JSON.parse(event.data)
      console.log(event.data)
      switch (msg.Type) {
        case 'round': {
          // TODO: display joined user
          break
        }
        case 'score': {
          console.log('I got a score')
          const data = JSON.parse(event.data)
          console.log(data.Payload.Scores)
          that.players = Object.keys(data.Payload.Scores)
          this.$store.commit('updateScoreBoard', data.Payload.Scores)
          break
        }
        default:
          console.log(event.data)
      }
    }
  },
  methods: {
    update (payload) {
      this.players = payload.Player
    },
    copyCode (clip) {
      navigator.clipboard.writeText(clip).then(() => {
        this.$message({ message: 'Copied to clipboard' })
      })
    },
    startGame () {
      console.log('Starting the game')
      const gameViewURL = this.$router.resolve({ name: 'Game', code: this.code }).path
      startGame(this.code)
        .then(() =>
          this.$router.push(gameViewURL)
        )
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
