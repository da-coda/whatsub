<template>
  <div>
    <img
      :src="imageSrc"
      style="max-width: 320px; max-height: 320px"
    >
  </div>
</template>

<script>
export default {
  name: 'ImagePost',
  props: {
    round: {
      type: Object,
      default: null
    }
  },
  computed: {
    imageSrc () {
      console.log('Data: ' + JSON.stringify(this.round))
      return this.round.Url
    }
  }
}
</script>

<style scoped>

</style>
