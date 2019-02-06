<template>
  <div class="watch">
    <div v-if="loading" class="ui active centered inline loader" />
    <template v-else>
      <div class="ui vertical segment">
        <div class="ui center aligned fluid video container">
          <video ref="video" :src="video.source" controls autoplay />
        </div>
      </div>
      <div class="ui vertical segment">
        <span class="header" v-text="video.title" />
        <p class="description" v-text="video.description" />
        <div class="ui right floated buttons">
          <a class="ui basic inverted icon button" :download="video.original_file_name" :href="video.source">
            <i class="download icon" />
            Download
          </a>
          <confirm-button class="ui basic red icon button" @confirm="deleteVideo">
            <i class="trash icon" />
            Delete
          </confirm-button>
          <router-link class="ui basic yellow icon button" :to="{ name: 'edit', params: { id: video.id } }">
            <i class="edit icon" />
            Edit
          </router-link>
        </div>
        <div class="tags">
          <router-link
            v-for="(tag, i) in video.tags"
            :key="i"
            class="ui label"
            v-text="tag"
            :to="{ name: 'search', query: { mode: 'tags', tags: tag } }"
          />
        </div>
      </div>
    </template>
  </div>
</template>

<script>
import ConfirmButton from '@/components/ConfirmButton';
import loadVideoById from './loadVideoById';

export default {
  components: {
    ConfirmButton,
  },
  mixins: [
    loadVideoById,
  ],
  metaInfo() {
    if (this.loading) {
      return {
        title: `Video ${this.id}`,
      };
    }

    return {
      title: this.video.title,
    };
  },
  methods: {
    deleteVideo() {
      this.loading = true;
      this.$store.dispatch('delete', parseInt(this.id, 10)).then(() => {
        this.$router.push({ name: 'home' });
      });
    },
  },
  beforeDestroy() {
    // attempt to unload video to free browser socket
    if (this.$refs.video) {
      this.$refs.video.pause();
      this.$refs.video.src = '';
    }
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

  &>.ui.segment:first-child {
    // force remove semantic-ui segment padding.
    // without this, top of video does not
    // match up with top of content on other pages.
    padding-top: 0px;
  }

  & .header {
    font-weight: bold;
    font-size: 3em;
    word-wrap: break-word;
  }

  & .description {
    margin-top: 1em;
  }
}
</style>