<template>
  <div>Ende</div>
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
  }
}
</script>

<style scoped>

</style>
