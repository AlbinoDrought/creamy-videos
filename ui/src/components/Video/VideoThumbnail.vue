<template>
  <router-link :to="{ name: 'watch', params: { id: value.id } }" class="ui fluid video card">
    <div class="ui image">
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
    word-wrap: break-word;
  }

  &>.ui.image {
    &, &>img {
      // remove semantic-ui border-radius
      border-radius: 0px;
    }

    // black background makes it appear
    // as if the bounding box is an extension
    // of the actual video thumbnail.
    background-color: #000;

    // force box to be 16:9
    // https://css-tricks.com/aspect-ratio-boxes/
    width: 100%;
    padding-bottom: 56.25%;
    // ...and hide anything that sticks out
    overflow: hidden;

    &>img {
      // make the image take up the full width
      // of the box
      min-width: 100%;
      height: auto;

      // force the image to vertically center itself
      // https://stackoverflow.com/a/28456704/3649573
      position: absolute;
      top: 50%;
      left: 0px;
      transform: translateY(-50%);
    }
  }
}
</style>
