<template>
  <div class="watch">
    <div v-if="loading" class="ui active centered inline loader" />
    <template v-else>
      <div class="ui vertical segment">
        <div class="ui center aligned fluid video container">
          <video :src="video.source" controls autoplay />
        </div>
      </div>
      <div class="ui vertical segment">
        <a class="ui basic inverted right floated icon button" download :href="video.source">
          <i class="download icon" />
          Download
        </a>

        <span class="header" v-text="video.title" />
        <p class="description" v-text="video.description" />
      </div>
    </template>
  </div>
</template>

<script>
export default {
  props: {
    id: {
      required: true,
    },
  },
  data() {
    return {
      video: {},
      loading: true,
    };
  },
  methods: {
    loadVideo() {
      if (!this.id) {
        return;
      }

      this.loading = true;
      this.video = {};
      this.$store.dispatch('video', parseInt(this.id, 10)).then(video => {
        this.video = video;
        this.loading = false;
      });
    },
  },
  mounted() {
    this.loadVideo();
  },
  watch: {
    id() {
      this.loadVideo();
    },
  },
};
</script>

<style lang="scss" scoped>
.ui.video.container {
  background-color: #000;
  max-height: 100%;
  margin: 0px !important;

  &>video {
    min-height: 60vh;
    max-height: 80vh;
    max-width: 100%;
  }
}

div.watch {
  color: rgb(171, 171, 171);

  & .header {
    font-weight: bold;
    font-size: 3em;
  }

  & .description {
    margin-top: 1em;
  }
}
</style>