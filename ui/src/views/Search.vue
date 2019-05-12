<template>
  <div class="search">
    <div v-if="loading" class="ui active centered inline loader" />
    <video-grid
      :loadable="true"
      :values="videos"
      @infinite="handleInfinite"
    />
  </div>
</template>

<script>
import VideoGrid from '@/components/Video/VideoGrid.vue';
import sortOptions from '@/sortOptions';

export default {
  name: 'search',
  components: {
    VideoGrid,
  },
  computed: {
    pageToFetch() {
      return this.page + this.infinitePage;
    },
    sortOption() {
      return sortOptions.find(option => option.key === this.sort) || sortOptions[0];
    },
  },
  data() {
    return {
      videos: [],
      infinitePage: 0,
      loading: true,
    };
  },
  metaInfo() {
    if (this.mode === 'tags' && this.tags) {
      return {
        title: `Tag Search: ${this.tags}`,
      };
    }

    if (this.mode === 'text' && this.text) {
      return {
        title: `Search: ${this.text}`,
      };
    }

    return {
      title: 'Search',
    };
  },
  methods: {
    fetchVideos() {
      this.loading = true;
      this.videos = [];
      this.infinitePage = 0;

      const promise = this.actuallyGetVideos();

      promise.then((videos) => {
        this.videos = videos;
        this.loading = false;
      });
    },
    infinitelyLoadVideos() {
      this.infinitePage += 1;
      return this.actuallyGetVideos()
        .then((videos) => {
          this.videos = this.videos.concat(videos);
          return videos;
        });
    },
    actuallyGetVideos() {
      return this.mode === 'tags'
        ? this.$store.dispatch(
            'tagged', 
            {
              tags: this.tags,
              page: this.pageToFetch,
              sortOption: this.sortOption,
            }
          )
        : this.$store.dispatch(
          'filtered', 
          { 
            filter: this.text, 
            page: this.pageToFetch,
            sortOption: this.sortOption,
          }
        );
    },
    handleInfinite(state) {
      this.infinitelyLoadVideos().then((newVideos) => {
        if (newVideos.length > 0) {
          state.loaded();
        } else {
          state.complete();
        }
      }).catch(() => state.error());
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
    sort() {
      this.fetchVideos();
    },
  },
  props: {
    page: {
      type: Number,
      rquired: false,
      default() {
        return 1;
      },
    },
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
    sort: {
      type: String,
      required: false,
      default() {
        return sortOptions[0].key;
      },
    },
  },
};
</script>
