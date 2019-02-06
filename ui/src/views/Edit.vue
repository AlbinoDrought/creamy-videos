<template>
  <div class="edit ui text container">
    <div class="ui inverted form">
      <div class="ui field">
        <label>Title</label>
        <input type="text" placeholder="Title" v-model="video.title">
      </div>
      <div class="ui field">
        <label>Tags (separated by comma)</label>
        <input type="text" placeholder="educational, computer science, wizardry" v-model="stringTags">
      </div>

      <div class="field">
        <label>Description</label>
        <textarea placeholder="Description" v-model="video.description" />
      </div>

      <div class="ui submit button" @click.prevent="edit" :class="{ loading }">
        Save
      </div>
    </div>
  </div>
</template>

<script>
import loadVideoById from './loadVideoById';

export default {
  name: 'edit',
  mixins: [
    loadVideoById,
  ],
  computed: {
    stringTags: {
      get() {
        return Array.isArray(this.video.tags)
          ? this.video.tags.join(',')
          : '';
      },
      set(value) {
        this.video.tags = value.split(',');
      },
    },
  },
  metaInfo() {
    if (this.loading) {
      return {
        title: `Edit ${this.id}`,
      };
    }

    return {
      title: `Edit ${this.video.title}`,
    };
  },
  methods: {
    edit() {
      this.loading = true;
      this.$store.dispatch('edit', this.video).then(({ id }) => this.$router.push({
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
div.edit {
  color: rgb(171, 171, 171);

  & .field>label {
    color: rgb(171, 171, 171);
  }
}
</style>
