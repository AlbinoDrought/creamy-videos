<template>
  <div class="search">
    <div v-if="loading" class="ui active centered inline loader" />
    <video-grid :values="videos" />
  </div>
</template>

<script>
import VideoGrid from '@/components/Video/VideoGrid';

export default {
  name: 'search',
  components: {
    VideoGrid,
  },
  data() {
    return {
      videos: [],
      loading: true,
    };
  },
  methods: {
    fetchVideos() {
      this.loading = true;
      this.videos = [];

      const promise = this.mode === 'tags'
        ? this.$store.dispatch('tagged', this.tags)
        : this.$store.dispatch('filtered', this.text);

      promise.then(videos => {
        this.videos = videos;
        this.loading = false;
      });
    },
  },
  mounted() {
    this.fetchVideos();
  },
  watch: {
    text() {
      this.fetchVideos();
    },
    tags() {
      this.fetchVideos();
    },
    mode() {
      this.fetchVideos();
    },
  },
  props: {
    mode: {
      type: String,
      required: false,
      default() {
        return 'text';
      },
    },
    text: {
      type: String,
      required: false,
      default() {
        return 'home';
      },
    },
    tags: {
      type: String,
      required: false,
      default() {
        return 'home';
      },
    },
  },
};
</script>
