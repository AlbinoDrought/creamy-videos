<template>
  <div class="ui stackable grid">
    <div v-for="value in values" :key="value.id" class="four wide column">
      <video-thumbnail :value="value" />
    </div>
    <infinite-loading
      v-if="loadable"
      class="sixteen wide column"
      @infinite="handleInfinite"
    />
  </div>
</template>

<script>
import InfiniteLoading from 'vue-infinite-loading';
import VideoThumbnail from './VideoThumbnail';

export default {
  components: {
    InfiniteLoading,
    VideoThumbnail,
  },

  methods: {
    handleInfinite(state) {
      // bubble the infinite event
      this.$emit('infinite', state);
    },
  },

  props: {
    loadable: {
      type: Boolean,
      required: false,
      default: false,
    },
    values: {
      type: Array,
      required: true,
      validator: v => v,
    },
  },
};
</script>
