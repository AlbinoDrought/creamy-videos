<template>
  <div class="upload ui text container">
    <div class="ui form">
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

      <semantic-progress-bar
        v-if="loading && progressMax > 0"
        class="inverted active"
        label-text="Uploading Video"
        :progress-current="progressCurrent"
        :progress-max="progressMax"
      />

      <div class="ui submit button" @click.prevent="upload" :class="{ loading }">
        Upload
      </div>
    </div>
  </div>
</template>

<script>
import SemanticProgressBar from '@/components/SemanticProgressBar.vue';

export default {
  name: 'upload',
  components: {
    SemanticProgressBar,
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
      progressCurrent: 0,
      progressMax: 0,
    };
  },
  methods: {
    handleFileChange(e) {
      [this.file] = e.target.files;
      if (!this.title) {
        this.title = this.file.name;
      }
    },
    async upload() {
      this.loading = true;
      this.progressCurrent = 0;
      this.progressMax = 0;

      try {
        const { id } = await this.$store.dispatch('upload', {
          formData: this.formData,
          config: {
            onUploadProgress: (progressEvent) => {
              if (!progressEvent.lengthComputable) {
                return;
              }

              this.progressCurrent = progressEvent.loaded;
              this.progressMax = progressEvent.total;
            },
          },
        });

        this.$router.push({
          name: 'watch',
          params: {
            id,
          },
        });
      } catch (ex) {
        this.loading = false;
        console.error(ex);
      }
    },
  },
};
</script>

<style lang="scss" scoped>
$form-text-color: rgb(171, 171, 171);

div.upload {
  color: $form-text-color;

  & .field {
    &>input, &>textarea {
      // intended to override semantic-ui defaults:
      @at-root &, &:focus {
        background-color: rgba(255, 255, 255, 0.1);
        color: white;
      }
    }

    &>label {
      color: $form-text-color;
    }
  }
}
</style>
