export default {
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
      this.$store.dispatch('video', parseInt(this.id, 10)).then((video) => {
        this.video = video;
        this.loading = false;
      });
    },
  },

  props: {
    id: {
      required: true,
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
