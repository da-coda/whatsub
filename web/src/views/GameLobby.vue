<template>
  <el-container>
    <el-header style="height: 260px">
      <BigLogo />
      <h4>Let the battle begin!</h4>
    </el-header>
    <el-main>
      <p>Copy this link to invite others: {{ gameLink }} or ask them to join a game with this code {{ code }}</p>
      <el-button
        v-if="isGameHead"
        type="success"
        style="width: 200px; margin-top: 10px"
        @click="$router.push('/startGame')"
      >
        Start the game
      </el-button>
    </el-main>
  </el-container>
</template>
<script>

import BigLogo from '@/views/BigLogo'
export default {
  name: 'NewGame',
  components: { BigLogo },
  props: {
    code: {
      type: String,
      default: null
    }
  },
  computed: {
    isGameHead () {
      return this.$store.state.isGameHead
    },
    gameLink () {
      return window.location.protocol + '//' + window.location.host + this.$route.path
    }
  },
  mounted () {
    this.$store.state.websocketConnection.onmessage = function (event) {
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
        default:
          console.log(event.data)
      }
    }
  },
  methods: {}
}
</script>
