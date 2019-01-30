<template>
  <div class="ui stackable grid">
    <div v-for="value in values" :key="value.id" class="four wide column">
      <video-thumbnail :value="value" />
    </div>
  </div>
</template>

<script>
import VideoThumbnail from './VideoThumbnail';

export default {
  components: {
    VideoThumbnail,
  },

  props: {
    values: {
      type: Array,
      required: true,
      validator: v => v,
    },
  },
};
</script>
