<template>
  <div>{{ round.Content }}</div>
</template>

<script>
export default {
  name: 'HtmlPost',
  props: {
    round: {
      type: Object,
      default: null
    }
  }
}
</script>

<style scoped>

</style>
