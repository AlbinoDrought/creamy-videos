<template>
  <router-link :to="{ name: 'watch', params: { id: value.id } }" class="ui fluid video card">
    <div class="ui medium image">
      <img :src="value.thumbnail">
    </div>
    <div class="content">
      <a class="header" v-text="value.title" />
    </div>
  </router-link>
</template>

<script>
export default {
  props: {
    value: {
      type: Object,
      required: true,
      validator: v => v &&
        v.id &&
        v.title &&
        v.thumbnail &&
        v.tags,
    }
  },
};
</script>

<style lang="scss" scoped>
.ui.video.card {
  border: 0px;
  border-radius: 0px;
  transition: none;
  box-shadow: none;
  background-color: inherit;

  & .content {
    padding: 0.5em 0.5em;
  }

  & .header {
    color: rgb(171, 171, 171);
  }

  &>.ui.image {
    &, &>img {
      border-radius: 0px;
    }

    height: 180px;
    background-color: #000;
    display: flex;
    // so what I actually wanted to do here was force image height to fill container,
    // cut off sides.
    // I had this going well, and then I added "display: flex"
    // it just worked, so I left it.
  }
}
</style>
