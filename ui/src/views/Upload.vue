<template>
  <div class="upload ui text container">
    <div class="ui inverted form">
      <div class="ui field">
        <label>Title</label>
        <input type="text" placeholder="Title" v-model="title">
      </div>
      <div class="ui field">
        <label>Tags (separated by comma)</label>
        <input type="text" placeholder="educational, computer science, wizardry" v-model="tags">
      </div>

      <div class="field">
        <label>Description</label>
        <textarea placeholder="Description" v-model="description" />
      </div>

      <div class="field">
        <label>File</label>
        <input type="file" @change="handleFileChange">
      </div>

      <div class="ui submit button" @click.prevent="upload" :class="{ loading }">
        Upload
      </div>
    </div>
  </div>
</template>

<script>
import VideoGrid from '@/components/Video/VideoGrid';

export default {
  name: 'upload',
  components: {
    VideoGrid,
  },
  computed: {
    formData() {
      const formData = new FormData();

      formData.append('title', this.title);
      formData.append('description', this.description);
      formData.append('tags', this.tags);
      formData.append('file', this.file);

      return formData;
    },
  },
  data() {
    return {
      title: '',
      tags: '',
      description: '',
      file: null,
      loading: false,
    };
  },
  methods: {
    handleFileChange(e) {
      this.file = e.target.files[0];
    },
    upload() {
      this.loading = true;
      this.$store.dispatch('upload', this.formData).then(() => this.$router.push({
        name: 'home',
      })).catch(() => {
        this.loading = false;
      });
    },
  },
};
</script>

<style lang="scss" scoped>
div.upload {
  color: rgb(171, 171, 171);

  & .field>label {
    color: rgb(171, 171, 171);
  }
}
</style>
