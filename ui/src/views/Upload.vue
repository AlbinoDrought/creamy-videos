<template>
  <div class="upload ui text container">
    <div class="ui inverted form">
      <div class="ui field">
        <label>Title</label>
        <input
          type="text"
          name="title"
          placeholder="Title"
          v-model="title"
        >
      </div>
      <div class="ui field">
        <label>Tags (separated by comma)</label>
        <input
          type="text"
          name="tags"
          placeholder="educational, computer science, wizardry"
          v-model="tags"
        >
      </div>

      <div class="field">
        <label>Description</label>
        <textarea
          name="description"
          placeholder="Description"
          v-model="description"
        />
      </div>

      <div class="field">
        <label>File</label>
        <input
          type="file"
          name="file"
          @change="handleFileChange"
        >
      </div>

      <div class="ui submit button" @click.prevent="upload" :class="{ loading }">
        Upload
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'upload',
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
  metaInfo: {
    title: 'Upload',
  },
  data() {
    return {
      title: '',
      tags: 'home',
      description: '',
      file: null,
      loading: false,
    };
  },
  methods: {
    handleFileChange(e) {
      this.file = e.target.files[0];
      if (!this.title) {
        this.title = this.file.name;
      }
    },
    upload() {
      this.loading = true;
      this.$store.dispatch('upload', this.formData).then(({ id }) => this.$router.push({
        name: 'watch',
        params: {
          id,
        },
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
