<template>
  <img
    src="../../public/img/logo.png"
    width="75"
    style="float: right;"
  >
  <div v-if="round">
    <div>Round: {{ round.Number }} of {{ round.From }}</div>
    <div>{{ round.Post.Title }}</div>
    <component
      :is="content"
      v-bind="{ round: round.Post }"
    />
    <div>
      <el-button
        v-for="subreddit in round.Subreddits"
        :key="subreddit"
      >
        {{ subreddit }}
      </el-button>
    </div>
  </div>
</template>
<script>

import { joinGame } from '@/lib/whatsub'
import { v4 as uuidv4 } from 'uuid'
import ImagePost from '@/components/game/ImagePost'
import HtmlPost from '@/components/game/HtmlPost'

export default {
  name: 'Game',
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
      players: [],
      round: {
        Post: {
          Title: null

        }
      }
    }
  },
  computed: {
    content () {
      const type = this.round.Post.Type
      switch (type) {
        case 'Image':
          return ImagePost
        case 'Text':
          return HtmlPost
      }
      return null
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
        case 'round': {
          const data = JSON.parse(event.data)
          console.log('Round: ' + data.Payload)
          that.round = data.Payload
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
