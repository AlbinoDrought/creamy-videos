<template>
  <div class="ui progress">
    <div class="bar" :style="{ transitionDuration: '300ms', width: humanizedProgressPercentage }">
      <div class="progress" v-text="humanizedProgressPercentage" />
    </div>
    <div class="label" v-if="labelText" v-text="labelText" />
  </div>
</template>

<script>
export default {
  computed: {
    progressPercentage() {
      if (this.progressMax <= 0) {
        return 0;
      }

      return (this.progressCurrent / this.progressMax) * 100;
    },
    humanizedProgressPercentage() {
      return Math.round(this.progressPercentage) + '%';
    },
  },

  props: {
    progressCurrent: {
      type: Number,
      required: true,
    },
    progressMax: {
      type: Number,
      required: true,
    },
    labelText: {
      type: String,
    },
  },
};
</script>
