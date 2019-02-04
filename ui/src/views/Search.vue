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
      this.$store.dispatch('filtered', this.text).then(videos => {
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
  },
  props: {
    text: {
      type: String,
      required: false,
      default() {
        return 'home';
      },
    },
  },
};
</script>
