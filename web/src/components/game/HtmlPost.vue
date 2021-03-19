<template>
  <div v-html="htmlContent" />
</template>

<script>
export default {
  name: 'HtmlPost',
  props: {
    round: {
      type: Object,
      default: null
    }
  },
  computed: {
    /**
     * Fetch the html content from the payload and decodes it to
     * raw HTML
     * //TODO This might be a security risk
     * @return {string} Raw HTML that has been decoded
     */
    htmlContent () {
      const doc = new DOMParser().parseFromString(this.round.Content, 'text/html')
      return doc.documentElement.textContent
    }
  }
}
</script>

<style scoped>

</style>
