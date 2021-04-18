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
        @click="chooseSubreddit(subreddit)"
      >
        {{ subreddit }}
      </el-button>
    </div>
  </div>
</template>
<script>

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
    const that = this
    const websocketConnection = this.$store.state.websocketConnection
    websocketConnection.onmessage = (event) => {
      const msg = JSON.parse(event.data)
      console.log(event.data)
      switch (msg.Type) {
        case 'round': {
          const data = JSON.parse(event.data)
          console.log('Round: ' + data.Payload)
          that.round = data.Payload
          break
        }
        case 'score': {
          console.log('I got a score')
          const data = JSON.parse(event.data)
          this.$store.commit('updateScoreBoard', data.Payload.Scores)
          break
        }
        case 'finished': {
          console.log('We are done!')
          this.$router.push('/game/' + this.code + '/finished')
        }
      }
    }
    websocketConnection.send('{}')// Dummy Ack
  },
  methods: {
    chooseSubreddit (subreddit) {
      console.log('You have chosen: ' + subreddit)
      const payload = {
        Type: 'answer',
        Payload: {
          Answer: subreddit
        }
      }
      this.$store.state.websocketConnection.send(JSON.stringify(payload))
    }
  }
}
</script>
<style scoped>

</style>
