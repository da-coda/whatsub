<template>
  <div>Ende</div>
  <div
    v-for="(score, name) in endScore"
    :key="name"
  >
    {{ name }} {{ score }}
  </div>
  <div>
    <router-link
      :to="{name: 'StartScreen'}"
    >
      <el-button>back</el-button>
    </router-link>
  </div>
</template>

<script>
export default {
  name: 'GameEndScreen',
  props: {
    code: {
      type: String,
      default: ''
    }
  },
  computed: {
    endScore () {
      return this.$store.state.scoreBoard
    }
  }
}
</script>

<style scoped>

</style>
