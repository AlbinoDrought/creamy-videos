<template>
  <button @click.prevent="handleClick">
    <slot v-if="!confirming" />
    <template v-else>
      Click {{ remainingClicks }} more times to confirm
    </template>
  </button>
</template>

<script>
export default {
  computed: {
    confirming() {
      return this.clicks > 0;
    },
    remainingClicks() {
      return (this.requiredClicks - this.clicks);
    },
  },

  data() {
    return {
      clicks: 0,
      timeoutHandle: false,
    };
  },

  methods: {
    handleClick() {
      clearInterval(this.timeoutHandle);
      this.clicks += 1;
      if (this.clicks >= this.requiredClicks) {
        this.clicks = 0;
        this.$emit('confirm');
      } else {
        this.timeoutHandle = setTimeout(() => {
          this.clicks = 0;
          this.timeoutHandle = false;
        }, 2000);
      }
    },
  },

  props: {
    requiredClicks: {
      type: Number,
      required: false,
      default: () => 4,
    },
  },
};
</script>
